package service

import (
	"context"
	"errors"
	"homework/internal/models"
	"homework/internal/ticket/repository"
	"strconv"
	"homework/pkg/log"
)

var ErrNotFound = errors.New("not found")

type ticketService struct {
	repo *repository.Repository
	log  log.Logger
}

func NewService(repo *repository.Repository, log log.Logger) (Service, error) {
	return &ticketService{
		repo: repo,
		log:  log,
	}, nil
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*models.Ticket, error) {
	return s.repo.Get(ctx, id)
}

func (s *ticketService) GetTicketsByUserID(ctx context.Context, userID int64) ([]*models.Ticket, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *ticketService) CreateTicket(ctx context.Context, userID int64, drawID int64, numbers []int, isAutoNumbers bool) (*models.Ticket, error) {
	s.log.Info("creating ticket",
		"draw_id", drawID,
		"numbers", numbers,
	)

	ticket := &models.Ticket{
		UserID:        strconv.FormatInt(userID, 10),
		DrawID:        strconv.FormatInt(drawID, 10),
		Numbers:       numbers,
		IsAutoNumbers: isAutoNumbers,
		Status:        models.TicketStatusActive,
	}

	return s.repo.Create(ctx, ticket)
}

func (s *ticketService) Get(ctx context.Context, id string) (*models.Ticket, error) {
	ticket, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}
