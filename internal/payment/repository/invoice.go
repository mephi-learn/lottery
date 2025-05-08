package repository

import (
	"context"
	"database/sql"
	"homework/internal/models"
	"homework/pkg/errors"
)

// CreateInvoice Создание инвойса
func (r *repository) CreateInvoice(ctx context.Context, invoice models.InvoiceStore) (int, error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return -1, errors.Errorf("authentificate need: %w", err)
	}

	var invoiceId int

	if err := r.db.QueryRowContext(ctx, "insert into invoices(user_id, ticket_id, status_id, date_invoice, price, status_change) values($1, $2, $3, $4, $5, $6) returning id",
		user.ID, invoice.TicketID, 1, invoice.RegisterTime, invoice.Amount, invoice.RegisterTime).Scan(&invoiceId); err != nil {
		return -1, errors.Errorf("failed to create invoice: %w", err)
	}

	return invoiceId, nil
}

// GetInvoice Получние инвойса по идентификатору
func (r *repository) GetInvoice(ctx context.Context, invoiceId int) (*models.InvoiceStore, error) {
	invoice := models.InvoiceStore{}

	if err := r.db.QueryRowContext(ctx, "select id, user_id, ticket_id, status_id, date_invoice, price from invoices where id = $1", invoiceId).
		Scan(&invoice.ID, &invoice.UserID, &invoice.TicketID, &invoice.StatusId, &invoice.RegisterTime, &invoice.Amount); err != nil {
		return nil, errors.Errorf("failed to get invoice: %w", err)
	}

	return &invoice, nil
}

// GetInvoiceByTicketId Получение инвойса по идентификатору билета
func (r *repository) GetInvoiceByTicketId(ctx context.Context, ticketId int) (*models.InvoiceStore, error) {
	invoice := models.InvoiceStore{}

	if err := r.db.QueryRowContext(ctx, "select id, user_id, ticket_id, status_id, date_invoice, price from invoices where ticket_id = $1", ticketId).
		Scan(&invoice.ID, &invoice.UserID, &invoice.TicketID, &invoice.StatusId, &invoice.RegisterTime, &invoice.Amount); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Errorf("failed to get invoice by ticketId: %w", err)
		}
	}

	return &invoice, nil
}
