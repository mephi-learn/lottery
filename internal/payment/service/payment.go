package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"homework/internal/models"
	"homework/pkg/errors"
	"math/big"
)

// RegisterPayment регистрация платежа.
func (s *paymentService) RegisterPayment(ctx context.Context, req *models.PaymentRequest) (err error) {

	invoice, err := s.repo.GetInvoice(ctx, req.InvoiceID)
	if err != nil {
		return errors.Errorf("failed to recieve invoice info: %w", err)
	}

	if invoice.Amount > req.Price {
		return errors.New("not enough money")
	}

	req.Price = invoice.Amount
	req.TicketID = invoice.TicketID

	// Обращаемся к платёжной системе
	if err = paymentSystemMock(req); err != nil {
		// Тут нужно вернуть деньги

		return errors.Errorf("failed to pay invoice: %w", err)
	}

	if err := s.ticket.BoughtTicket(ctx, req.TicketID); err != nil {
		// Надо подумать, что будет, если при ошибке маркировки билета купленным и отмене маркировки билета произойдёт ещё одна ошибка
		//if err = s.ticket.CancelTicket(ctx, req.TicketID); err != nil {
		//	return errors.Errorf("failed to buying a ticket: %w", err)
		//}

		return errors.Errorf("failed to buy ticket: %w", err)
	}

	return nil
}

func paymentSystemMock(req *models.PaymentRequest) error {
	if req.CVC == 123 {
		return nil
	}

	// Успешный платёж с вероятностью 80% (примерно)
	nBig, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return fmt.Errorf("failed to generate random number: %w", err)
	}

	if int(nBig.Int64()) < 80 {
		return nil
	}

	return errors.New("processing failed")
}
