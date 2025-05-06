package models

import (
	"time"
)

type WinTickets map[string][]*Ticket

type DrawResult struct {
	Id             int       `json:"id"`
	DrawId         int       `json:"draw_id"`
	WinCombination []int     `json:"win_combination"`
	ResultTime     time.Time `json:"result_time"`
}

type DrawResultStore struct {
	Id             int    `json:"id"`
	DrawId         int    `json:"draw_id"`
	DrawStatusId   int    `json:"status_id"`
	LotteryType    string `json:"lottery_type"`
	WinCombination []int  `json:"win_combination" db:"win_combination"`
}

type DrawingResult struct {
	WinTickets WinTickets
	Statistic  map[string]int
}
