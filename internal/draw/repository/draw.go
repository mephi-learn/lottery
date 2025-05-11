package repository

import (
	"context"
	"database/sql"
	"homework/internal/models"
	"homework/pkg/errors"
	"time"
)

func (r *repository) CreateDraw(ctx context.Context, draw *models.DrawStore) (int, error) {
	var drawId int
	if err := r.db.QueryRowContext(ctx, "insert into draws(status_id, lottery_type, cost, sale_date, start_date) values($1, $2, $3, $4, $5) returning id",
		draw.StatusId, draw.LotteryType, draw.Cost, draw.SaleDate, draw.StartDate).Scan(&drawId); err != nil {
		return -1, errors.Errorf("failed to create draw: %w", err)
	}

	return drawId, nil
}

func (r *repository) GetDraw(ctx context.Context, drawId int) (*models.DrawStore, error) {
	draw := models.DrawStore{}
	if err := r.db.QueryRowContext(ctx, "SELECT id, status_id, lottery_type, cost, sale_date, start_date FROM draws WHERE id = $1", drawId).Scan(&draw.Id, &draw.StatusId, &draw.LotteryType, &draw.Cost, &draw.SaleDate, &draw.StartDate); err != nil {
		return nil, errors.Errorf("failed to get draw: %w", err)
	}

	return &draw, nil
}

func (r *repository) PlannedDraw(ctx context.Context, drawId int) error {
	result, err := r.db.ExecContext(ctx, "update draws set status_id = $1 where id = $2", models.DrawStatusPlanned, drawId)
	if err != nil {
		return errors.Errorf("failed to update draw status: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return errors.Errorf("draw not found: %w", err)
	}

	return nil
}

func (r *repository) ActiveDraw(ctx context.Context, drawId int) error {
	result, err := r.db.ExecContext(ctx, "update draws set status_id = $1 where id = $2", models.DrawStatusActive, drawId)
	if err != nil {
		return errors.Errorf("failed to update draw status: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return errors.Errorf("draw not found: %w", err)
	}

	return nil
}

func (r *repository) CompletedDraw(ctx context.Context, drawId int) error {
	result, err := r.db.ExecContext(ctx, "update draws set status_id = $1 where id = $2", models.DrawStatusCompleted, drawId)
	if err != nil {
		return errors.Errorf("failed to update draw status: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return errors.Errorf("draw not found: %w", err)
	}

	return nil
}

func (r *repository) CancelDraw(ctx context.Context, drawId int) error {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.Errorf("authentificate need: %w", err)
	}

	tr, err := r.db.Begin()
	defer func() {
		_ = tr.Rollback()
	}()
	if err != nil {
		return errors.Errorf("failed to initialize storage transaction: %w", err)
	}

	// Изменение статуса отмененного тиража
	result, err := r.db.ExecContext(ctx, "update draws set status_id = $1 where id = $2", models.DrawStatusCanceled, drawId)
	if err != nil {
		return errors.Errorf("failed to update draw status: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return errors.Errorf("draw not found: %w", err)
	}

	// Получение суммы оплаченных инвойсов для возврату пользователю
	var refundAmount float64
	if err := r.db.QueryRowContext(ctx,
		"select coalesce(sum(price), 0) as price from invoices where user_id = $1 and status_id = 2 and ticket_id in (select id from tickets where draw_id = $2)", user.ID, drawId).
		Scan(&refundAmount); err != nil {
		return errors.Errorf("failed to get refund amount: %w", err)
	}

	// Изменение статусов оплаченных инвойсов на отменён
	_, err = r.db.ExecContext(ctx, "update invoices set status_id = $1 where user_id = $2 and status_id = 2 and ticket_id in (select id from tickets where draw_id = $3)", models.InvoiceStatusCanceled, user.ID, drawId)
	if err != nil {
		return errors.Errorf("failed to update invoice status: %w", err)
	}

	// Если сумма к возврату больше нуля, переводим средства на кошелек пользователя
	if refundAmount > 0 {
		_, err = r.db.ExecContext(ctx, "update users set wallet = wallet + $1 where id = $2", refundAmount, user.ID)
		if err != nil {
			return errors.Errorf("failed funds transfer: %w", err)
		}
	}

	if err = tr.Commit(); err != nil {
		return errors.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) FailedDraw(ctx context.Context, drawId int) error {
	result, err := r.db.ExecContext(ctx, "update draws set status_id = $1 where id = $2", models.DrawStatusFailed, drawId)
	if err != nil {
		return errors.Errorf("failed to update draw status: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return errors.Errorf("draw not found: %w", err)
	}

	return nil
}

func (r *repository) SetDrawSaleDate(ctx context.Context, drawId int, begin time.Time) error {
	result, err := r.db.ExecContext(ctx, "update draws set dale_date = $1 where id = $2", begin.Unix(), drawId)
	if err != nil {
		return errors.Errorf("failed to update draw sale_date: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return errors.Errorf("draw not found: %w", err)
	}

	return nil
}

func (r *repository) SetDrawStartDate(ctx context.Context, drawId int, start time.Time) error {
	result, err := r.db.ExecContext(ctx, "update draws set dale_date = $1 where id = $2", start.Unix(), drawId)
	if err != nil {
		return errors.Errorf("failed to update draw start_date: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		return errors.Errorf("draw not found: %w", err)
	}

	return nil
}

func (r *repository) ListActiveDraw(ctx context.Context) ([]models.DrawStore, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, status_id, lottery_type, cost, sale_date, start_date FROM draws WHERE status_id in ($1, $2)", models.DrawStatusPlanned, models.DrawStatusActive)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var draws []models.DrawStore
	for rows.Next() {
		var draw models.DrawStore
		if err := rows.Scan(&draw.Id, &draw.StatusId, &draw.LotteryType, &draw.Cost, &draw.SaleDate, &draw.StartDate); err != nil {
			return draws, err
		}
		draws = append(draws, draw)
	}
	if err = rows.Err(); err != nil {
		return draws, err
	}

	return draws, nil
}

func (r *repository) ListReadyToBeginDraws(ctx context.Context) ([]*models.DrawStore, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, status_id, lottery_type, cost, sale_date, start_date FROM draws WHERE status_id = $1 and start_date is not null and start_date < $2", models.DrawStatusPlanned, time.Now())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var draws []*models.DrawStore
	for rows.Next() {
		var draw models.DrawStore
		if err := rows.Scan(&draw.Id, &draw.StatusId, &draw.LotteryType, &draw.Cost, &draw.SaleDate, &draw.StartDate); err != nil {
			return nil, err
		}
		if !draw.StartDate.IsZero() {
			draws = append(draws, &draw)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return draws, nil
}

func (r *repository) ListCompletedDraw(ctx context.Context) ([]models.DrawStore, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, status_id, lottery_type, cost, sale_date, start_date FROM draws WHERE status_id = $1", models.DrawStatusCompleted)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var draws []models.DrawStore
	for rows.Next() {
		var draw models.DrawStore
		if err := rows.Scan(&draw.Id, &draw.StatusId, &draw.LotteryType, &draw.Cost, &draw.SaleDate, &draw.StartDate); err != nil {
			return draws, err
		}
		draws = append(draws, draw)
	}
	if err = rows.Err(); err != nil {
		return draws, err
	}

	return draws, nil
}

// GetDrawByTicketId получение тиража по идентификатору билета
func (r *repository) GetDrawByTicketId(ctx context.Context, ticketId int) (*models.DrawStore, error) {
	draw := models.DrawStore{}
	if err := r.db.QueryRowContext(ctx, "SELECT id, cost, status_id, lottery_type, cost, sale_date, start_date FROM draws WHERE id = (select draw_id from tickets where id = $1)", ticketId).
		Scan(&draw.Id, &draw.Cost, &draw.StatusId, &draw.LotteryType, &draw.Cost, &draw.SaleDate, &draw.StartDate); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Errorf("failed to get draw by ticket id: %w", err)
		}

		return nil, nil
	}

	return &draw, nil
}
