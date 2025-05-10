package repository

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

// DebitingFundsFromWallet Списание средств с кошелька пользователя
func (r *repository) DebitingFundsFromWallet(ctx context.Context, amount float64) error {
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

	var balance float64
	if err = r.db.QueryRowContext(ctx, "update users set wallet = wallet - $1 where id = $2 returning wallet", amount, user.ID).Scan(&balance); err != nil {
		return errors.Errorf("failed to debiting funds from wallet: %w", err)
	}

	if balance < 0 {
		return errors.New("no money in wallet")
	}

	if err = tr.Commit(); err != nil {
		return errors.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAmountInUserWallet Получение суммы на кошельке пользователя
func (r *repository) GetAmountInUserWallet(ctx context.Context) (float64, error) {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return -1, errors.Errorf("authentificate need: %w", err)
	}

	var amountInWallet float64

	if err := r.db.QueryRowContext(ctx, "select wallet from users where id = $1", user.ID).Scan(&amountInWallet); err != nil {
		return -1, errors.Errorf("failed to debiting funds from wallet: %w", err)
	}

	return amountInWallet, nil
}
