package repository

import (
	"context"
	"database/sql"
	"homework/internal/models"
	"homework/pkg/errors"
	"time"
)

func (r *repository) StoreTicket(ctx context.Context, ticket *models.Ticket) error {
	var ticketId int
	if err := r.db.QueryRowContext(ctx, "insert into tickets(status_id, draw_id, data, user_id, lock_time) values($1, $2, $3, $4, $5) returning id",
		ticket.Status, ticket.DrawId, ticket.Data, ticket.UserId, ticket.LockTime).Scan(&ticketId); err != nil {
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

func (r *repository) LoadTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error) {
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

func (r *repository) ListAvailableTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	var userId sql.NullInt64
	rows, err := r.db.QueryContext(ctx, `
SELECT t.id, t.status_id, t.draw_id, data, user_id
FROM tickets t INNER JOIN draws d ON t.draw_id = d.id
WHERE 1 = 1
    and d.status_id = $1
    and t.status_id = $2
    and t.draw_id = $3
  	and t.user_id is null
  	and t.lock_time is null`, models.DrawStatusPlanned, models.TicketStatusReady, drawId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var tickets []*models.Ticket
	for rows.Next() {
		var ticket models.Ticket
		if err = rows.Scan(&ticket.Id, &ticket.Status, &ticket.DrawId, &ticket.Data, &userId); err != nil {
			return nil, err
		}
		if userId.Valid {
			ticket.UserId = int(userId.Int64)
		}
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}

func (r *repository) LoadTicketsByUserId(ctx context.Context, userId int) ([]*models.Ticket, error) {
	//rows, err := r.db.QueryContext(ctx, "SELECT t.id, t.status_id, t.draw_id, data FROM tickets t INNER JOIN draws d ON t.draw_id = d.id WHERE t.draw_id = $1", drawId)
	//if err != nil {
	//	return nil, err
	//}
	//defer func() {
	//	_ = rows.Close()
	//}()
	//
	//var tickets []*models.Ticket
	//for rows.Next() {
	//	var ticket models.Ticket
	//	if err = rows.Scan(&ticket.Id, &ticket.Status, &ticket.DrawId, &ticket.Data); err != nil {
	//		return nil, err
	//	}
	//	tickets = append(tickets, &ticket)
	//}
	//
	//return tickets, nil

	return nil, nil
}

func (r *repository) GetTicketById(ctx context.Context, ticketId int) (*models.Ticket, error) {
	ticket := models.Ticket{}
	if err := r.db.QueryRowContext(ctx, "select id, status_id, draw_id, data from tickets where id = $1", ticketId).Scan(&ticket.Id, &ticket.Status, &ticket.DrawId, &ticket.Data); err != nil {
		return nil, errors.Errorf("failed to get ticket: %w", err)
	}

	return &ticket, nil
}

func (r *repository) MarkTicketAsBought(ctx context.Context, ticketId int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tickets SET status_id = $1, lock_time = NULL WHERE id = $2",
		models.TicketStatusBought, ticketId)
	if err != nil {
		return errors.Errorf("failed to update ticket status: %w", err)
	}
	return nil
}

func (r *repository) ReserveTicket(ctx context.Context, ticketId int, userId int, lockTime time.Time) error {
	var resultTicketId int
	err := r.db.QueryRowContext(ctx, "UPDATE tickets SET status_id = $1, user_id = $2, lock_time = $3 WHERE id = $4 and lock_time is null returning id",
		models.TicketStatusReady, userId, lockTime, ticketId).Scan(&resultTicketId)
	if err != nil {
		return errors.Errorf("failed to reserve ticket: %w", err)
	}
	return nil
}

func (r *repository) CancelTicket(ctx context.Context, ticketId int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tickets SET status_id = $1, user_id = NULL, lock_time = NULL WHERE id = $2",
		models.TicketStatusReady, ticketId)
	if err != nil {
		return errors.Errorf("failed to cancel ticket: %w", err)
	}
	return nil
}

func (r *repository) GetExpiredTickets(ctx context.Context) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM tickets WHERE lock_time < NOW() AND lock_time IS NOT NULL")
	if err != nil {
		return nil, errors.Errorf("failed to get expired tickets: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ticketIds []int
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ticketIds = append(ticketIds, id)
	}

	return ticketIds, nil
}
