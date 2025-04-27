package service

import (
	"testing"

	"errors"
	"homework/internal/models"
	"homework/pkg/log"
)

type mockLotteryService struct {
	lottery models.Lottery
	err     error
}

func (m *mockLotteryService) LotteryByType(name string) (models.Lottery, error) {
	return m.lottery, m.err
}

func TestCreateTickets(t *testing.T) {
	logger, err := log.New(log.LoggerConfig{Level: "debug"})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	lottery := models.NewLottery536()
	mockLotteryService := &mockLotteryService{lottery: lottery}

	svc, err := NewTicketService(
		WithLogger(logger),
		WithLotteryService(mockLotteryService),
	)
	if err != nil {
		t.Fatalf("failed to create ticket service: %v", err)
	}

	tickets, err := svc.CreateTickets(1, "536", 5)
	if err != nil {
		t.Fatalf("failed to create tickets: %v", err)
	}

	if len(tickets) != 5 {
		t.Errorf("expected 5 tickets, got %d", len(tickets))
	}

	for _, ticket := range tickets {
		if ticket.DrawId != 1 {
			t.Errorf("expected draw id 1, got %d", ticket.DrawId)
		}
		if ticket.Status != models.TicketStatusReady {
			t.Errorf("expected status ready, got %v", ticket.Status)
		}
	}
}

func TestCreateTicketsInvalidLottery(t *testing.T) {
	logger, err := log.New(log.LoggerConfig{Level: "debug"})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	mockLotteryService := &mockLotteryService{err: errors.New("lottery not found")}

	svc, err := NewTicketService(
		WithLogger(logger),
		WithLotteryService(mockLotteryService),
	)
	if err != nil {
		t.Fatalf("failed to create ticket service: %v", err)
	}

	_, err = svc.CreateTickets(1, "invalid", 5)
	if err == nil {
		t.Error("expected error for invalid lottery type")
	}
}
