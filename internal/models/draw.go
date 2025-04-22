package models

import (
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

type Draw struct {
	status atomic.Value

	startTicketBuy time.Time
	endTicketBuy   time.Time
}
