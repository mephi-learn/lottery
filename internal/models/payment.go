package models

import (
	"time"
)

// Invoice представляет инвойс для оплаты билета.
type Invoice struct {
	ID           int       `json:"id"`
	RegisterTime time.Time `json:"register_time"`
	Status       string    `json:"status"`  // Например, "pending", "paid", "failed"
	UserID       string    `json:"user_id"` // ID пользователя, к которому относится инвойс, хз надо или не надо.
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
	InvoiceID  int     `json:"invoice_id"` // Добавлено: ID инвойса для оплаты
	Price      float64 `json:"price"`      // Добавлено: Сумма платежа
	UserID     int     `json:"user_id"`    // Добавлено: ID пользователя, который платит
	TicketID   int     `json:"ticket_id"`
}
