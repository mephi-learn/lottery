package models

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"homework/pkg/errors"
	"math/big"
	"strconv"
	"strings"
)

const (
	l536id                = "5from36"
	l535combinationLength = 5
	l536minAllowDigit     = 1
	k536maxAllowDigit     = 36
)

type lottery536Ticket struct {
	Id          int
	Status      TicketStatus
	DrawId      int
	Combination []int
}

type Lottery536 struct {
	tickets []*Ticket
}

func NewLottery536() *Lottery536 {
	return &Lottery536{
		tickets: make([]*Ticket, 0),
	}
}

func (l *Lottery536) Type() string {
	return "536"
}

func (l *Lottery536) Name() string {
	return "5 из 36"
}

func (l *Lottery536) Create() Lottery {
	return NewLottery536()
}

func (l *Lottery536) AddTickets(tickets []*Ticket) error {
	for _, ticket := range tickets {
		if err := l.validateTicket(ticket); err != nil {
			return err
		}
		l.tickets = append(l.tickets, ticket)
	}
	return nil
}

func (l *Lottery536) CreateTickets(drawId int, num int) ([]*Ticket, error) {
	tickets := make([]*Ticket, 0, num)
	for i := 0; i < num; i++ {
		ticket, err := l.generateTicket(drawId)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (l *Lottery536) Drawing(combination []int) (map[string][]*Ticket, error) {
	if len(combination) != 5 {
		return nil, errors.Errorf("invalid combination length: %d", len(combination))
	}

	result := make(map[string][]*Ticket)
	for _, ticket := range l.tickets {
		matches := l.countMatches(ticket, combination)
		if matches >= 3 {
			result[strconv.Itoa(matches)] = append(result[strconv.Itoa(matches)], ticket)
		}
	}
	return result, nil
}

func (l *Lottery536) validateTicket(ticket *Ticket) error {
	var numbers []int
	if err := json.Unmarshal([]byte(ticket.Data), &numbers); err != nil {
		return errors.Errorf("invalid ticket data: %w", err)
	}

	if len(numbers) != 5 {
		return errors.Errorf("invalid numbers count: %d", len(numbers))
	}

	unique := make(map[int]bool)
	for _, num := range numbers {
		if num < 1 || num > 36 {
			return errors.Errorf("invalid number: %d", num)
		}
		if unique[num] {
			return errors.Errorf("duplicate number: %d", num)
		}
		unique[num] = true
	}

	return nil
}

func (l *Lottery536) generateTicket(drawId int) (*Ticket, error) {
	numbers := make([]int, 0, 5)
	unique := make(map[int]bool)

	for len(numbers) < 5 {
		num, err := rand.Int(rand.Reader, big.NewInt(36))
		if err != nil {
			return nil, errors.Errorf("failed to generate random number: %w", err)
		}
		n := int(num.Int64()) + 1
		if !unique[n] {
			numbers = append(numbers, n)
			unique[n] = true
		}
	}

	data, err := json.Marshal(numbers)
	if err != nil {
		return nil, errors.Errorf("failed to marshal numbers: %w", err)
	}

	return &Ticket{
		DrawId: drawId,
		Status: TicketStatusReady,
		Data:   string(data),
	}, nil
}

func (l *Lottery536) countMatches(ticket *Ticket, combination []int) int {
	var numbers []int
	if err := json.Unmarshal([]byte(ticket.Data), &numbers); err != nil {
		return 0
	}

	matches := 0
	for _, num := range numbers {
		for _, comb := range combination {
			if num == comb {
				matches++
				break
			}
		}
	}
	return matches
}

// Преобразует билет из общего формата во внутренний
func (l *Lottery536) fromTicket(rawTicket *Ticket) (*lottery536Ticket, error) {
	data, err := json.Marshal(rawTicket.Data)
	if err != nil {
		return nil, errors.Errorf("failed to marshal ticket data: %w", err)
	}

	var combination []int
	if err := json.Unmarshal(data, &combination); err != nil {
		return nil, errors.Errorf("failed to unmarshal ticket data: %w", err)
	}

	if len(combination) != l535combinationLength {
		return nil, errors.Errorf("invalid ticket combination length: %d", len(combination))
	}

	if err = l.validateCombination(combination); err != nil {
		return nil, errors.New("invalid ticket combination")
	}

	return &lottery536Ticket{
		Id:          rawTicket.Id,
		Status:      rawTicket.Status,
		DrawId:      rawTicket.DrawId,
		Combination: combination,
	}, nil
}

// Преобразует билет из внутреннего формата в общий
func (l *Lottery536) toTicket(rawTicket *lottery536Ticket) *Ticket {
	digits := make([]string, len(rawTicket.Combination))
	for i, digit := range rawTicket.Combination {
		digits[i] = fmt.Sprintf("%02d", digit)
	}
	combination := strings.Join(digits, ",")

	return &Ticket{
		Id:     rawTicket.Id,
		Status: rawTicket.Status,
		DrawId: rawTicket.DrawId,
		Data:   combination,
	}
}

// Производит проверку на соответствие комбинации цифр правилам лотереи
func (l *Lottery536) validateCombination(combination []int) error {
	if len(combination) != l535combinationLength {
		return errors.New("invalid combination length")
	}

	uniq := map[int]struct{}{}

	// Проверяем что все числа из комбинации укладываются в допустимый диапазон
	for _, digit := range combination {
		if digit < l536minAllowDigit || digit > k536maxAllowDigit {
			return errors.New("invalid digit in combination: " + strconv.Itoa(digit))
		}
		uniq[digit] = struct{}{}
	}

	// Проверяем, что все цифры уникальные
	if len(uniq) != len(combination) {
		return errors.New("all digit bust be unique")
	}

	return nil
}

// Проверяет, уникальна ли комбинация среди уже существующих билетов
func (l *Lottery536) checkUniqCombination(combination []int, newTickets ...*lottery536Ticket) bool {
	// Проверяем комбинацию на уникальность по существующим билетам
	for _, ticket := range l.tickets {
		var numbers []int
		if err := json.Unmarshal([]byte(ticket.Data), &numbers); err != nil {
			continue
		}
		found := true
		for i := range numbers {
			if combination[i] != numbers[i] {
				found = false
				break
			}
		}
		if found {
			return false
		}
	}

	// Проверяем комбинацию на уникальность по новым билетам
	for _, ticket := range newTickets {
		found := true
		for i := range ticket.Combination {
			if combination[i] != ticket.Combination[i] {
				found = false
				break
			}
		}
		if found {
			return false
		}
	}

	return true
}
