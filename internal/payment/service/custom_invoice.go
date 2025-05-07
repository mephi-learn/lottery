package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"time"
)

// RegisterCustomInvoice регистрация инвойса с указанием билета.
func (s *paymentService) RegisterCustomInvoice(ctx context.Context, drawId int, combination []int) (invoiceId int, err error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return -1, errors.New("unauthenticated user")
	}

	ticket, err := s.ticket.CreateReservedTicket(ctx, drawId, combination)
	if err != nil {
		return -1, errors.Errorf("failed to reserve custom ticket: %w", err)
	}

	if err = s.ticket.ReserveTicket(ctx, ticket.Id, user.ID); err != nil {
		return -1, errors.Errorf("failed to reserve ticket %d: %w", ticket.Id, err)
	}

	var invoice models.InvoiceStore

	invoice.RegisterTime = time.Now()
	invoice.StatusId = 1
	invoice.TicketID = ticket.Id

	invoiceId, err = s.repo.CreateInvoice(ctx, invoice)

	return invoiceId, err
}
