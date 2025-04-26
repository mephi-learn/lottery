package repository

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (r *repository) StoreTicket(ctx context.Context, ticket *models.Ticket) error {
	var ticketId int
	if err := r.db.QueryRowContext(ctx, "insert into tickets(status_id, draw_id, data) values($1, $2, $3) returning id",
		ticket.Status, ticket.DrawId, ticket.Data).Scan(&ticketId); err != nil {
		return errors.Errorf("failed to store ticket: %w", err)
	}

	return nil
}

func (r *repository) StoreTickets(ctx context.Context, tickets []*models.Ticket) error {
	tr, err := r.db.Begin()
	if err != nil {
		return errors.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tr.Rollback()
	}()

	for _, ticket := range tickets {
		var ticketId int
		if err := tr.QueryRowContext(ctx, "insert into tickets(status_id, draw_id, data) values($1, $2, $3) returning id",
			ticket.Status, ticket.DrawId, ticket.Data).Scan(&ticketId); err != nil {
			return errors.Errorf("failed to store ticket: %w", err)
		}
		ticket.Id = ticketId
	}

	return tr.Commit()
}

func (r *repository) LoadTickets(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT t.id, t.status_id, t.draw_id, data FROM tickets t INNER JOIN draws d ON t.draw_id = d.id WHERE t.draw_id = $1", drawId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var tickets []*models.Ticket
	for rows.Next() {
		var ticket models.Ticket
		if err = rows.Scan(&ticket.Id, &ticket.Status, &ticket.DrawId, &ticket.Data); err != nil {
			return nil, err
		}
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}

func (r *repository) GetTicket(ctx context.Context, ticketId int) (*models.Ticket, error) {
	ticket := models.Ticket{}
	if err := r.db.QueryRowContext(ctx, "select id, status_id, draw_id, data from tickets where id = $1", ticketId).Scan(&ticket.Id, &ticket.Status, &ticket.DrawId, &ticket.Data); err != nil {
		return nil, errors.Errorf("failed to get ticket: %w", err)
	}

	return &ticket, nil
}
