package repository

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

// Создание инвойс
func (r *repository) CreateInvoice(ctx context.Context, invoice models.Invoice) (int, error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return -1, errors.Errorf("authentificate need: %w", err)
	}

	var invoiceId int
	if err := r.db.QueryRowContext(ctx, "insert into invoices(user_id, ticket_id, status_id, date_invoice, price, status_change) values($1, $2, $3, $4, $5, $6) returning id",
		user.ID, invoice.TicketID, 1, invoice.RegisterTime, 0, invoice.RegisterTime).Scan(&invoiceId); err != nil {
		return -1, errors.Errorf("failed to create invoice: %w", err)
	}

	return invoiceId, nil
}
