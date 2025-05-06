package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"time"
)

// RegisterInvoice регистрация инвойса.
func (s *paymentService) RegisterInvoice(ctx context.Context, ticketId int) (invoiceId int, err error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return -1, errors.New("unauthenticated user")
	}
	if err = s.ticket.ReserveTicket(ctx, ticketId, user.ID); err != nil {
		return -1, errors.Errorf("failed to reserve ticket %d: %w", ticketId, err)
	}

	var invoice models.Invoice

	invoice.RegisterTime = time.Now()
	invoice.Status = "pending"
	invoice.TicketID = ticketId

	invoiceId, err = s.repo.CreateInvoice(ctx, invoice)

	return invoiceId, err
}
