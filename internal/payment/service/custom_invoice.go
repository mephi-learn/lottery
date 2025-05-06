package service

import (
	"context"
	"fmt"
	"homework/internal/models"
	"homework/pkg/errors"
	"strings"
	"time"
)

// RegisterCustomInvoice регистрация инвойса с указанием билета.
func (s *paymentService) RegisterCustomInvoice(ctx context.Context, drawId int, combination []int) (invoiceId int, err error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return -1, errors.New("unauthenticated user")
	}

	data := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(combination)), ","), "[]")

	ticket, err := s.ticket.CreateReservedTicket(ctx, drawId, data)
	if err != nil {
		return -1, errors.Errorf("failed to reserve custom ticket: %w", err)
	}

	if err = s.ticket.ReserveTicket(ctx, ticket.Id, user.ID); err != nil {
		return -1, errors.Errorf("failed to reserve ticket %d: %w", ticket.Id, err)
	}

	var invoice models.Invoice

	invoice.RegisterTime = time.Now()
	invoice.Status = "pending"
	invoice.TicketID = ticket.Id

	invoiceId, err = s.repo.CreateInvoice(ctx, invoice)

	return invoiceId, err
}
