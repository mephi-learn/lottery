package models

import (
	"time"
)

const (
	DrawStatusUnknown DrawStatus = iota
	DrawStatusPlanned
	DrawStatusActive
	DrawStatusCompleted
	DrawStatusCanceled
	DrawStatusFailed
)

type DrawStatus int

func (d *DrawStatus) String() string {
	switch *d {
	case DrawStatusPlanned:
		return "planned"
	case DrawStatusActive:
		return "active"
	case DrawStatusCompleted:
		return "completed"
	case DrawStatusCanceled:
		return "canceled"
	case DrawStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func DrawStatusFromString(status string) DrawStatus {
	switch status {
	case "planned":
		return DrawStatusPlanned
	case "active":
		return DrawStatusActive
	case "completed":
		return DrawStatusCompleted
	case "canceled":
		return DrawStatusCanceled
	case "failed":
		return DrawStatusFailed
	default:
		return DrawStatusUnknown
	}
}

type DrawInput struct {
	Status    string    `json:"status"`
	Lottery   string    `json:"lottery"`
	SaleDate  time.Time `json:"sale_date"`
	StartDate time.Time `json:"start_date"`
	Cost      float64   `json:"cost"`
}

type DrawOutput struct {
	Id        int       `json:"id"`
	Status    string    `json:"status"`
	Lottery   string    `json:"lottery"`
	SaleDate  time.Time `json:"sale_date"`
	StartDate time.Time `json:"start_date"`
}

type DrawStore struct {
	Id          int       `json:"id"`
	StatusId    int       `json:"status_id"`
	LotteryType string    `json:"lottery_type"`
	SaleDate    time.Time `json:"sale_date"`
	StartDate   time.Time `json:"start_date"`
	Cost        float64   `json:"cost"`
}

type Draw struct {
	Id        int        `json:"id"`
	Status    DrawStatus `json:"status"`
	Lottery   Lottery    `json:"lottery"`
	SaleDate  time.Time  `json:"sale_date"`
	StartDate time.Time  `json:"start_date"`
}
