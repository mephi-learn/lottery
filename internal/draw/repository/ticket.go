package repository

import (
	"context"
	"fmt"
	"homework/internal/models"
	"homework/pkg/errors"
	"strings"
)

const batchSize = 100

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

func (r *repository) LoadParticipatingTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT t.id, t.status_id, t.draw_id, data, t.cost FROM tickets t INNER JOIN draws d ON t.draw_id = d.id WHERE t.draw_id = $1 and t.status_id in ($2, $3, $4)",
		drawId, models.TicketStatusBought, models.TicketStatusWin, models.TicketStatusLose)
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

func (r *repository) MarkDrawTickets(ctx context.Context, drawId int, winTickets []int) error {
	tr, err := r.db.Begin()
	if err != nil {
		return errors.Errorf("failed to begin storage transaction: %w", err)
	}
	defer func() {
		_ = tr.Rollback()
	}()

	// Формируем пачки билетов, поскольку мы не можем запихать миллиард номеров в SQL инструкцию IN, то будем помечать выигрышные билеты пачками
	batches := make([]string, 0)
	for i := range len(winTickets)/batchSize + 1 {
		end := (i + 1) * batchSize
		if end > len(winTickets) {
			end = len(winTickets)
		}
		batches = append(batches, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(winTickets[i*batchSize:end])), ","), "[]"))
	}

	// Всем выигрышным купленным билетам выставляем соответствующий статус
	for _, batch := range batches {
		_, err = tr.ExecContext(ctx, "UPDATE tickets SET status_id = $1 WHERE id in("+batch+") and status_id = $2", models.TicketStatusWin, models.TicketStatusBought)
		if err != nil {
			return errors.Errorf("failed to update win ticket status: %w", err)
		}
	}

	// Всем не выигрышным купленным билетам также выставляем соответствующий статус
	_, err = tr.ExecContext(ctx, "UPDATE tickets SET status_id = $1 WHERE draw_id = $2 and status_id = $3", models.TicketStatusLose, drawId, models.TicketStatusBought)
	if err != nil {
		return errors.Errorf("failed to update lose ticket status: %w", err)
	}

	if err = tr.Commit(); err != nil {
		return errors.Errorf("failed to commit storage transaction: %w", err)
	}

	return nil
}
