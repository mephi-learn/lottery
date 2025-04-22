package log

import (
	"context"
	"homework/pkg/errors"
	"homework/pkg/log/filter"
	"homework/pkg/log/logutil"
	"io"
	"log/slog"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const zlTimeFormat = time.RFC3339Nano

type zlGroup struct {
	name  string
	attrs []slog.Attr
}

func (zlg zlGroup) String() string {
	return zlg.name
}

// eventAttrs добавляет все групповые события к указанному evt. То же, что и [zlEventAttrs].
func (zlg zlGroup) eventAttrs(evt *zerolog.Event) *zerolog.Event {
	return zlEventAttrs(evt, zlg.attrs)
}

type zlHandler struct {
	zlGroup // Атрибуты верхнего уровня, когда нет групп.

	filter filter.Filter[zlGroup, Level]
	groups []zlGroup

	hook logutil.Hook
	zl   zerolog.Logger
}

// newLogHandler создает новый [slog.Handler], имплементация основана на [zerolog.Logger].
func newLogHandler(opt options) (zlHandler, error) {
	var h zlHandler
	h.hook = opt.hook

	h.filter = filter.NewBackwardFilter[zlGroup](opt.Filters, keySeparator)

	var zlw io.Writer

	switch {
	case len(opt.writers) == 0: // логер по умолчанию.
		zlw = zerolog.NewConsoleWriter()
	case len(opt.writers) == 1:
		var err error
		if zlw, err = zlWriter(opt.writers[0]); err != nil {
			return h, errors.Errorf("create zerolog writer: %w", err)
		}
	default:
		iow := make([]io.Writer, len(opt.writers))

		for i, w := range opt.writers {
			var err error
			if iow[i], err = zlWriter(w); err != nil {
				return h, errors.Errorf("create zerolog writer (%d/%d): %w", i, len(opt.writers), err)
			}
		}

		zlw = zerolog.MultiLevelWriter(iow...)
	}

	level := zlevel(opt.Level)
	h.zl = zerolog.New(zlw).Level(level)

	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Path != "" {
		h.name = bi.Main.Path + "/"
	}

	return h, nil
}

func (h zlHandler) Enabled(_ context.Context, level slog.Level) bool {
	ours := h.zl.GetLevel()

	// Если у нас есть фильтр, то переопределяем уровень.
	group, ok := h.filter.Get(h.groups)
	if ok {
		ours = zlevel(group)
	}

	return ours <= zlevel(level)
}

func (h zlHandler) Handle(ctx context.Context, rec slog.Record) error {
	caller, fn := h.trace(rec.PC)

	//nolint:zerologlint // Msg вызовется в defer.
	evt := h.zl.Log().Ctx(ctx).
		Str("time", rec.Time.Format(zlTimeFormat)). // ПРИМЕЧАНИЕ: zerolog вставляет время только когда вы добавляете hook с помощью метода [zerolog.Context.Timestamp].
		Str("level", rec.Level.String()).           // Если нам нужно изменить имя поля или значение для уровня, это можно сделать здесь.
		Str("caller", caller).                      // Добавляем местоположение вызова функции логирования.
		Str("func", fn)

	defer func() { evt.Msg(rec.Message) }() // Это завершает событие. Убеждаемся что это всегда последняя операция.

	// Добавляем верхний уровень атрибутов.
	evt = h.eventAttrs(evt)

	// Использование атрибутов из контекста.
	// Их не надо добавлять в подгруппу, как атрибуты ниже.
	evt = zlEventAttrs(evt, h.hook.Attrs(ctx))

	// Добавляем атрибуты, предоставленные [slog.Record].
	ln := len(h.groups)
	if ln == 0 {
		rec.Attrs(h.attrIter(ctx, &evt))
		return nil
	}

	// Если есть хотя бы одна группа, вставим их здесь.
	subevt := zerolog.Dict()
	rec.Attrs(h.attrIter(ctx, &subevt))

	// Сворачиваем группу детей в единственный dict.
	// Имена групп добавляются в обратном порядке, при формировании получаем естественный.
	if ln > 1 {
		bottom := h.groups[:ln-1]
		for _, zlg := range bottom {
			// Без нового пустого словаря zerolog добавляет новый dict к концу, а не в основание.
			subevt = zerolog.Dict().Dict(zlg.name, zlg.eventAttrs(subevt))
		}
	}

	top := h.groups[ln-1]
	evt = evt.Dict(top.name, top.eventAttrs(subevt))

	return nil
}

// attrIter добавляет атрибут к [zerolog.Event] с фильтрацией ошибок.
func (h zlHandler) attrIter(ctx context.Context, dict **zerolog.Event) func(attr slog.Attr) bool {
	// Debug логи должны включать трассировку ошибок.
	if h.Enabled(ctx, Debug) {
		return func(attr slog.Attr) bool {
			*dict = zlAppend(*dict, attr)
			return true
		}
	}

	// подавляем излишние сообщения об ошибках.
	return func(attr slog.Attr) bool {
		// если error имплементирует [slog.LogValuer], просто используем строковое представление.
		if val := attr.Value; val.Kind() == slog.KindLogValuer {
			if err, ok := val.Any().(error); ok {
				attr.Value = slog.StringValue(err.Error())
			}
		}

		*dict = zlAppend(*dict, attr)

		return true
	}
}

func (h zlHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Теоретически мы можем кодировать attrs как zerolog.Event и это вероятно быстрее.
	// Но мы не можем добавлять значения словаря без вложенности. Таким образом, здесь нам нужен плоский список атрибутов.

	if len(attrs) == 0 {
		return h
	}

	var hattrs *[]slog.Attr
	if len(h.groups) == 0 {
		hattrs = &h.attrs
	} else {
		hattrs = &h.groups[0].attrs
	}

	*hattrs = append(*hattrs, attrs...)

	return h
}

func (h zlHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	groups := make([]zlGroup, len(h.groups)+1)
	groups[0] = zlGroup{name: name}

	copy(groups[1:], h.groups)

	h.groups = groups

	return h
}

// trace возвращает "func:line" и имя функции как аргументы "caller" и "func".
//
//nolint:nonamedreturns // В противном случае в месте вызова мы не сможем сказать что мы возвращаем.
func (h zlHandler) trace(pc uintptr) (caller, name string) {
	fn := runtime.FuncForPC(pc)
	file, line := fn.FileLine(pc)

	file = strings.TrimPrefix(file, h.name)
	fname := strings.TrimPrefix(fn.Name(), h.name)

	return file + ":" + strconv.Itoa(line), fname
}

// zlevel преобразует уровень slog в уровень zerolog.
func zlevel(slevel slog.Level) zerolog.Level {
	switch {
	case slevel == disabled:
		return zerolog.Disabled
	case slevel >= slog.LevelError:
		return zerolog.ErrorLevel
	case slevel >= slog.LevelWarn:
		return zerolog.WarnLevel
	case slevel >= slog.LevelInfo:
		return zerolog.InfoLevel
	default:
		return zerolog.DebugLevel
	}
}

// zlAppend возвращает объект zerolog с добавленными парами ключ - значение.
func zlAppend[T zlMapper](zlv zlValuer[T], attr slog.Attr) T {
	key, val := attr.Key, attr.Value

	switch val.Kind() {
	// составной тип.
	case slog.KindGroup:
		dict := zlEventAttrs(zerolog.Dict(), val.Group())
		return zlv.Dict(key, dict)
	// тип имплементирует slog.LogValuer.
	case slog.KindLogValuer:
		attr.Value = val.LogValuer().LogValue()
		return zlAppend(zlv, attr)
	// обычные типы.
	case slog.KindBool:
		return zlv.Bool(key, val.Bool())
	case slog.KindFloat64:
		return zlv.Float64(key, val.Float64())
	case slog.KindInt64:
		return zlv.Int64(key, val.Int64())
	case slog.KindUint64:
		return zlv.Uint64(key, val.Uint64())
	case slog.KindString:
		return zlv.Str(key, val.String())
	case slog.KindDuration:
		return zlv.Dur(key, val.Duration())
	case slog.KindTime:
		return zlv.Str(key, val.Time().Format(zlTimeFormat))
	case slog.KindAny:
		// Некоторые ошибки не могут конвертироваться в JSON, рассматриваем их как строки.
		if err, ok := val.Any().(error); ok {
			return zlv.Str(key, err.Error())
		}

		fallthrough
	default:
		// Рассматриваем остальные типы как any.
		return zlv.Any(key, val.Any())
	}
}

// zlEventAttrs добавляет атрибуты slog к [zerolog.Event].
func zlEventAttrs(evt *zerolog.Event, attrs []slog.Attr) *zerolog.Event {
	for _, attr := range attrs {
		evt = zlAppend(evt, attr)
	}

	return evt
}

// zlWriter конфигурирует [io.Writer] от [writer].
func zlWriter(w writer) (io.Writer, error) {
	var iow io.Writer

	cfg := w.Config()
	switch cfg.Format {
	case FormatConsole:
		iow = zerolog.ConsoleWriter{Out: w}
	case FormatJSON:
		iow = w
	case formatNone:
		return nil, errors.New("format not specified")
	default:
		return nil, errors.Errorf("unknown format '%s'", cfg.Format)
	}

	return iow, nil
}

// zlMap определяет ограничение типа zerolog.Context и *zerolog.Event для использования дженериков.
// Их использование одинаково, но поскольку это различные типы, а их методы возвращают значения (не интерфейсы),
// В противном случае невозможно реализовать функции, которые могут использовать оба типа.
type zlMapper interface {
	zerolog.Context | *zerolog.Event
}

// zlValuer оборачивает типы zerolog так чтобы была цепочка из пар.
type zlValuer[T zlMapper] interface {
	Dict(key string, dict *zerolog.Event) T
	Any(key string, val any) T
	Bool(key string, val bool) T
	Str(key string, val string) T
	Float64(key string, val float64) T
	Int64(key string, val int64) T
	Uint64(key string, val uint64) T
	Dur(key string, val time.Duration) T
}
