package draw

import (
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
	Create(begin time.Time, start time.Time) (drawId int, err error) // Создание тиража, указывается дата начала и окончания приёма билетов
	Cancel(drawId int)                                               // Отмена тиража, все деньги возвращаются клиентам
	SetBeginTime(drawId int, begin time.Time)                        // Установка времени начала продажи билетов
	SetStartTime(drawId int, start time.Time)                        // Установка времени начала тиража
}

type Draw struct {
	status atomic.Value

	startTicketBuy time.Time
	endTicketBuy   time.Time
	repository     Repository

	log zerolog.Logger
}
