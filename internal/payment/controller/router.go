package controller

import (
	"context"
	"homework/internal/auth"
	"homework/pkg/errors"
	"homework/pkg/log"
	"net/http"
)

type handler struct {
	service paymentService
	log     log.Logger
}

type HandlerOption func(*handler)

func NewHandler(opts ...HandlerOption) (*handler, error) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	if h.log == nil {
		return nil, errors.New("logger is missing")
	}

	return h, nil
}

func WithLogger(logger log.Logger) HandlerOption {
	return func(o *handler) {
		o.log = logger
	}
}

// WithService добавляет [paymentService] в обработчик запросов.
func WithService(svc paymentService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type paymentService interface {
	RegisterInvoice(ctx context.Context, ticketId int) (err error)
	RegisterPayment(ctx context.Context, paymentId int, payment float64) (err error)
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	// Invoice
	mux.Handle("POST /api/invoice/{ticket_id}", auth.Authenticated(h.RegisterInvoice))
	// Payment
	mux.Handle("POST /api/payments/{invoice_id}", auth.Authenticated(h.RegisterPayment))
}

// 1. Если идёт выбор билета, то из тикет сервиса получаем список доступных билетов (статус тиража: запланирован, билет без user_id в статусе: готов)
// 2. Если пользователь добавляет билет сам, то тикет сервис просто добавляет билет в лотерею (статус тиража: запланирован)
// 	в обоих случаях возвращается его номер
// 3. Создаём инвойс с номером билета (сообщаем тикет сервису о резервировании)
// 4. Ждём оплаты
// 5. Если оплата инвойса успешна, то тикет сервису сообщается об этом и билет переводится в статус КУПЛЕН
// Если не успешно или истекло время оплаты, то вызывается сценарий отмены в тикет сервисе, и билет:
// 1. Стирается user_id и билет снова готов к покупке
// 2. Билет удаляется из базы

// /ticket/allow_list/{draw_id} - получение списка доступных билетов
// CreateReservedTicket - резервируешь билет, возвращается номер
// /ticket/add/{draw_id} - добавление билета, возвращается номер

// генерируешь инвойс, ждёшь оплату

// Успешно:
// /ticket/success/{ticket_id} - билет куплен, возвращается ОК

// Неуспешно:
// /ticket/failure/{ticket_id} - билет не куплен (ошибка оплаты), билет или удаляется из БД, или возвращается в общий пул доступных билетов (стирается user_id)
