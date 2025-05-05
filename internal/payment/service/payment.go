package service

import (
	"context"
	"homework/pkg/errors"
)

// RegisterPayment регистрация платежа.
func (s *paymentService) RegisterPayment(ctx context.Context, invoiceId int, payment float64) (err error) {
	// Проверяем, прошёл ли платёж
	// Если не прошёл, возвращаем ошибку и деньги

	// Получаем инвойс
	// Из него получаем номер билета (его id)
	ticketId := 0

	if err := s.ticket.BoughtTicket(ctx, ticketId); err != nil {
		// Надо подумать, что будет, если при ошибке маркировки билета купленным и отмене маркировки билета произойдёт ещё одна ошибка
		if err = s.ticket.CancelTicket(ctx, ticketId); err != nil {
			return errors.Errorf("failed to buying a ticket: %w", err)
		}

		return errors.Errorf("failed to buy ticket: %w", err)
	}

	return nil
}
