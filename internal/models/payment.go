package models

import (
	"time"

	"github.com/google/uuid"
)

// Invoice представляет инвойс для оплаты билета.
type Invoice struct {
	ID           uuid.UUID `json:"id"`
	TicketData   any       `json:"ticketData"` // Заменить any на данные билета
	RegisterTime time.Time `json:"registerTime"`
	Status       string    `json:"status"` // Например, "pending", "paid", "failed"
	UserID       string    `json:"userID"` // ID пользователя, к которому относится инвойс, хз надо или не надо.
}

// Payment представляет информацию о платеже.
type Payment struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // "SUCCESS", "FAILED"
	PaymentTime time.Time `json:"paymentTime"`
	InvoiceID   string    `json:"invoiceID"` // Ссылка на ID инвойса
}

// Структура для запроса регистрации платежа
type PaymentRequest struct {
	CardNumber string  `json:"cardNumber"`
	CVC        string  `json:"cvc"`
	InvoiceID  string  `json:"invoiceID"` // Добавлено:  ID инвойса для оплаты
	Amount     float64 `json:"amount"`    // Добавлено:  Сумма платежа
	UserID     string  `json:"userID"`    // Добавлено: ID пользователя, который платит
}
