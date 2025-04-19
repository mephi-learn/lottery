package draw

import (
	"context"
	"github.com/rs/zerolog"
	"sync/atomic"
	"time"
)

const (
	StatusUnknown Status = iota
	StatusPlanned
	StatusActive
	StatusCompleted
	StatusCanceled
	StatusFailed
)

type Status int

func (d *Status) String() string {
	switch *d {
	case StatusPlanned:
		return "planned"
	case StatusActive:
		return "active"
	case StatusCompleted:
		return "completed"
	case StatusCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}

type Option func(*Draw) error

// Repository реализует интерфейс репозитория тиража
type Repository interface {
	Create(ctx context.Context, begin time.Time, start time.Time, lotteryType string) (drawId int, err error) // Создание тиража, указывается дата начала и окончания приёма билетов
	ListActive(ctx context.Context) ([]Draw, error)                                                           // Список активных тиражей
	Cancel(ctx context.Context, drawId int) error                                                             // Отмена тиража, все деньги возвращаются клиентам
	SetBeginTime(ctx context.Context, drawId int, begin time.Time) error                                      // Установка времени начала продажи билетов
	SetStartTime(ctx context.Context, drawId int, start time.Time) error                                      // Установка времени начала тиража
}

type Draw struct {
	status atomic.Value

	startTicketBuy time.Time
	endTicketBuy   time.Time
	repository     Repository

	log zerolog.Logger
}
