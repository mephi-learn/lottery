package repository

import (
	"context"
	"homework/internal/models"
)

func (r *repository) LoadTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT t.id, t.status_id, t.draw_id, data, t.cost FROM tickets t INNER JOIN draws d ON t.draw_id = d.id WHERE t.draw_id = $1", drawId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var tickets []*models.Ticket
	for rows.Next() {
		var ticket models.Ticket
		if err = rows.Scan(&ticket.Id, &ticket.Status, &ticket.DrawId, &ticket.Data, &ticket.Cost); err != nil {
			return nil, err
		}
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}

func (r *repository) LoadBoughtTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT t.id, t.status_id, t.draw_id, data, t.cost FROM tickets t INNER JOIN draws d ON t.draw_id = d.id WHERE t.draw_id = $1 and t.status_id = $2", drawId, models.TicketStatusBought)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var tickets []*models.Ticket
	for rows.Next() {
		var ticket models.Ticket
		if err = rows.Scan(&ticket.Id, &ticket.Status, &ticket.DrawId, &ticket.Data, &ticket.Cost); err != nil {
			return nil, err
		}
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}
