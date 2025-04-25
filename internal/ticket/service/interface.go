package service

import (
	"context"
	"homework/internal/models"
)

type Service interface {
	GetTicket(ctx context.Context, id string) (*models.Ticket, error)
	GetTicketsByUserID(ctx context.Context, userID int64) ([]*models.Ticket, error)
	CreateTicket(ctx context.Context, userID int64, drawID int64, numbers []int, isAutoNumbers bool) (*models.Ticket, error)
}
