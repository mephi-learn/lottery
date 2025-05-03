package models

import "encoding/base64"

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

type Ticket struct {
	Id     int          `json:"id"`
	Status TicketStatus `json:"status"`
	DrawId int          `json:"draw_id"`
	Data   string       `json:"data"`
}

type Ticket1 interface {
	Draw() Draw                  // Вывод тиража
	Lottery() Lottery            // Пип лотереи
	String() string              // Вывод карточки
	Status() TicketStatus        // Вывод статуса
	Marshal() ([]byte, error)    // Маршалинг
	Unmarshal(data []byte) error // Демаршалинг
}

type TicketResult struct {
	WinCombination []int `json:"win_combination"`
	Combination    []int `json:"combination"`
	WinCount       int   `json:"win_count"`
}
