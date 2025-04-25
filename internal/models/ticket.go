package models

import "time"

type TicketStatus string

const (
	TicketStatusActive  TicketStatus = "active"
	TicketStatusWin     TicketStatus = "win"
	TicketStatusLose    TicketStatus = "lose"
	TicketStatusDeleted TicketStatus = "deleted"
)

func (s TicketStatus) String() string {
	return string(s)
}

type LotteryType string

const (
	TypeUserNumbers LotteryType = "USER_NUMBERS"
	TypeAutoNumbers LotteryType = "AUTO_NUMBERS"
)

func (t LotteryType) String() string {
	return string(t)
}

const (
	MinNumber    = 1
	MaxNumber    = 36
	NumbersCount = 5
)

type Ticket struct {
	ID            string       `json:"id"`
	UserID        string       `json:"user_id"`
	DrawID        string       `json:"draw_id"`
	Numbers       []int        `json:"numbers"`
	Status        TicketStatus `json:"status"`
	LotteryType   string       `json:"lottery_type"`
	IsAutoNumbers bool         `json:"is_auto_numbers"`
	CreatedAt     time.Time    `json:"created_at"`
}

type TicketRequest struct {
	DrawID  string `json:"draw_id"`
	Numbers []int  `json:"numbers"`
}
