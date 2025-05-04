package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

// RegisterInvoice регистрация инвойса.
func (s *paymentService) RegisterInvoice(ctx context.Context, ticketId int) (err error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.New("authenticate need")
	}

	if err = s.ticket.ReserveTicket(ctx, ticketId, user.ID); err != nil {
		return errors.Errorf("failed to reserve ticket: %w", err)
	}

	return nil
}
