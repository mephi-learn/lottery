package models

import (
	"time"
)

const (
	InvoiceStatusUnknown InvoiceStatus = iota
	InvoiceStatusActive
	InvoiceStatusPaid
	InvoiceStatusCanceled
	InvoiceStatusFailed
)

type InvoiceStatus int

func (d *InvoiceStatus) String() string {
	switch *d {
	case InvoiceStatusActive:
		return "active"
	case InvoiceStatusPaid:
		return "completed"
	case InvoiceStatusCanceled:
		return "canceled"
	case InvoiceStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func InvoiceStatusFromString(status string) InvoiceStatus {
	switch status {
	case "active":
		return InvoiceStatusActive
	case "completed":
		return InvoiceStatusPaid
	case "canceled":
		return InvoiceStatusCanceled
	case "failed":
		return InvoiceStatusFailed
	default:
		return InvoiceStatusUnknown
	}
}

// Invoice представляет инвойс для оплаты билета.
type Invoice struct {
	ID           int       `json:"id"`
	RegisterTime time.Time `json:"register_time"`
	Status       string    `json:"status"`  // Например, "pending", "paid", "failed"
	UserID       string    `json:"user_id"` // ID пользователя, к которому относится инвойс, хз надо или не надо.
	TicketID     int       `json:"ticket_id"`
	Amount       float64   `json:"amount"`
}

// InvoiceStore представляет структуру инвойта для оплаты билета, взятую из базы данных.
type InvoiceStore struct {
	ID           int       `json:"id"`
	RegisterTime time.Time `json:"register_time"`
	StatusId     int       `json:"status_id"` // Например, "pending", "paid", "failed"
	UserID       string    `json:"user_id"`   // ID пользователя, к которому относится инвойс, хз надо или не надо.
	TicketID     int       `json:"ticket_id"`
	Amount       float64   `json:"amount"`
}

// Payment представляет информацию о платеже.
type Payment struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // "SUCCESS", "FAILED"
	PaymentTime time.Time `json:"payment_time"`
	InvoiceID   string    `json:"invoice_id"` // Ссылка на ID инвойса
}

// Структура для запроса регистрации платежа.
type PaymentRequest struct {
	CardNumber string  `json:"card_number"`
	CVC        int     `json:"cvc"`
	Price      float64 `json:"price"`   // Добавлено: Сумма платежа
	UserID     int     `json:"user_id"` // Добавлено: ID пользователя, который платит
	TicketID   int     `json:"ticket_id"`
}
