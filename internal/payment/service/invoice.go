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

	// Проверка существования инвойса
	storedInsvoice, err := s.repo.GetInvoiceByTicketId(ctx, ticketId)
	if err != nil {
		return -1, errors.New(err)
	}
	if storedInsvoice != nil {
		return -1, errors.New("invoice for ticketId already exists")
	}

	// Получение билета и проверка его статуса
	ticket, err := s.ticket.GetTicketById(ctx, ticketId)
	if err != nil {
		return -1, errors.Errorf("failed to get ticket: %w", err)
	}
	if ticket.Status > 1 {
		status := models.TicketStatus(ticket.Status)
		return -1, errors.Errorf("ticket has status %q", status.String())
	}

	// Получение тиража и проверка его статуса
	draw, err := s.draw.GetDrawByTicketId(ctx, ticketId)
	if err != nil {
		return -1, errors.Errorf("failed to get draw: %w", err)
	}
	if draw.StatusId > 2 {
		status := models.DrawStatus(draw.StatusId)
		return -1, errors.Errorf("draw has status %q", status.String())
	}

	// Резервирование билета
	if err = s.ticket.ReserveTicket(ctx, ticketId, user.ID); err != nil {
		return -1, errors.Errorf("failed to reserve ticket %d: %w", ticketId, err)
	}

	var invoice models.InvoiceStore

	invoice.RegisterTime = time.Now()
	invoice.StatusId = 1
	invoice.TicketID = ticketId
	invoice.Amount = draw.Cost

	// Создание инвойса
	invoiceId, err = s.repo.CreateInvoice(ctx, invoice)
	if err != nil {
		return -1, errors.Errorf("failed to create invoice: %w", err)
	}

	return invoiceId, err
}
