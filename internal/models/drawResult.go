package models

import (
	"time"

	"github.com/lib/pq"
)

type DrawResult struct {
	Id             int       `json:"id"`
	DrawId         int       `json:"draw_id"`
	WinCombination []int     `json:"win_combination"`
	ResultTime     time.Time `json:"result_time"`
}

type DrawResultStore struct {
	Id             int           `json:"id"`
	DrawId         int           `json:"draw_id"`
	DrawStatusId   int           `json:"status_id"`
	LotteryType    string        `json:"lottery_type"`
	WinCombination pq.Int64Array `json:"win_combination" db:"win_combination"`
}
