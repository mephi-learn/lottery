package service

import (
	"context"
	"homework/internal/models"
	"time"

	"github.com/google/uuid"
)

// RegisterInvoice регистрация инвойса.
func (s *paymentService) RegisterInvoice(ctx context.Context, ticketId int) (invoiceId int, err error) {
	var invoice models.Invoice

	invoice.ID = uuid.New()
	invoice.RegisterTime = time.Now()
	invoice.Status = "pending"
	invoice.TicketID = ticketId

	invoiceId, err = s.repo.CreateInvoice(ctx, invoice)

	return invoiceId, err
}
