package models

import (
	"encoding/base64"
	"fmt"
	"homework/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const (
	TicketStatusUnknown TicketStatus = iota // Статус неизвестен (ошибка)
	TicketStatusReady                       // Билет готов к покупке
	TicketStatusBought                      // Билет куплен
	TicketStatusWin                         // Билет выиграл
	TicketStatusLose                        // Билет проиграл
)

func (d TicketStatus) String() string {
	switch d {
	case TicketStatusReady:
		return "ready"
	case TicketStatusBought:
		return "bought"
	case TicketStatusWin:
		return "win"
	case TicketStatusLose:
		return "lose"
	default:
		return "unknown"
	}
}

func TicketStatusFromString(status string) TicketStatus {
	switch status {
	case "ready":
		return TicketStatusReady
	case "bought":
		return TicketStatusBought
	case "win":
		return TicketStatusWin
	case "lose":
		return TicketStatusLose
	default:
		return TicketStatusUnknown
	}
}

type TicketStatus int

type TicketStore struct {
	Id       int    `json:"id"`
	StatusId int    `json:"status_id"`
	DrawId   int    `json:"draw_id"`
	Data     string `json:"data"`
	UserId   int    `json:"user_id"`
}

func (t *TicketStore) Marshal(data string) {
	t.Data = base64.StdEncoding.EncodeToString([]byte(data))
}

func (t *TicketStore) Unmarshal() (string, error) {
	result, err := base64.StdEncoding.DecodeString(t.Data)
	return string(result), err
}

type TicketResult struct {
	WinCombination []int `json:"win_combination"`
	Combination    []int `json:"combination"`
	WinCount       int   `json:"win_count"`
	Id             int   `json:"id"`
}

type Ticket struct {
	Id       int          `json:"id"`
	Status   TicketStatus `json:"status"`
	DrawId   int          `json:"draw_id"`
	UserId   int          `json:"user_id"`
	Data     string       `json:"data"`
	Cost     float64      `json:"cost"`
	LockTime time.Time    `json:"lock_time"`
}

func ParseTicketCombination(combination string) (ticketNumbers []int, err error) {
	data, err := base64.StdEncoding.DecodeString(combination)
	if err != nil {
		return nil, errors.New("unknown decode ticket data")
	}
	parts := strings.SplitN(string(data), ";", 2)

	numberStr := parts[1]

	digitStrings := strings.Split(numberStr, ",")
	ticketNumbers = make([]int, 0, len(digitStrings))

	for i, digitStr := range digitStrings {
		digitStr = strings.TrimSpace(digitStr) // Handle potential spaces
		if digitStr == "" {
			err = fmt.Errorf("invalid ticket combination format: empty number string at index %d", i)
			return
		}

		digit, parseErr := strconv.Atoi(digitStr)
		if parseErr != nil {
			err = fmt.Errorf("invalid ticket combination format: failed to parse number '%s' at index %d: %w", digitStr, i, parseErr)
			return
		}
		ticketNumbers = append(ticketNumbers, digit)
	}

	return ticketNumbers, nil
}
