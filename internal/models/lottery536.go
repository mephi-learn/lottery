package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"homework/pkg/errors"
	"math/big"
	"strconv"
	"strings"
)

const (
	lotteryId536      = "5from36"
	combinationLength = 5
	minAllowDigit     = 1
	maxAllowDigit     = 36
)

type lottery536Ticket struct {
	Id          int
	Status      TicketStatus
	DrawId      int
	Combination []int
}
type Lottery536 struct {
	notTemplate bool
	Tickets     []*lottery536Ticket
}

func (l *Lottery536) Name() string {
	return "5 из 36"
}

func (l *Lottery536) Type() string {
	return lotteryId536
}

func (l *Lottery536) Create() Lottery {
	return &Lottery536{notTemplate: true}
}

// AddTickets добавляет билет в лотерею
func (l *Lottery536) AddTickets(tickets []*Ticket) error {
	// Не даём использовать заготовку
	if !l.notTemplate {
		return errors.New("use template")
	}

	for _, rawTicket := range tickets {
		ticket, err := l.fromTicket(rawTicket)
		if err != nil {
			return errors.Errorf("failed convert ticket: %w", err)
		}
		l.Tickets = append(l.Tickets, ticket)
	}

	return nil
}

func (l *Lottery536) CreateTickets(num int, drawId int) ([]*Ticket, error) {
	// Не даём использовать заготовку
	if !l.notTemplate {
		return nil, errors.New("use template")
	}

	newTickets := make([]*lottery536Ticket, 0, num)
	maxDigitIndex := big.NewInt(maxAllowDigit - 1)
	for range num {

		// Создаём уникальную комбинацию чисел
		combination := make([]int, combinationLength)
		for l.validateCombination(combination) != nil || !l.checkUniqCombination(combination, newTickets...) {
			for i := range combinationLength {
				randomNumber, err := rand.Int(rand.Reader, maxDigitIndex)
				if err != nil {
					return nil, errors.Errorf("failed to create combination: %w", err)
				}
				combination[i] = int(randomNumber.Int64()) + 1
			}
		}

		newTickets = append(newTickets, &lottery536Ticket{Status: TicketStatusReady, DrawId: drawId, Combination: combination})
	}

	// Конвертируем список билетов из внутреннего формата во внешний
	result := make([]*Ticket, len(newTickets))
	for i, ticket := range newTickets {
		result[i] = l.toTicket(ticket)
	}

	return result, nil
}

// Преобразует билет из общего формата во внутренний
func (l *Lottery536) fromTicket(rawTicket *Ticket) (*lottery536Ticket, error) {
	data, err := base64.StdEncoding.DecodeString(rawTicket.Data)
	if err != nil {
		return nil, errors.New("unknown decode ticket data")
	}

	split := strings.SplitN(string(data), ";", 2)
	if len(split) != 2 {
		return nil, errors.New("unknown ticket format")
	}

	if len(split[1]) != combinationLength*3-1 {
		return nil, errors.New("invalid ticket combination length")
	}

	if split[0] != lotteryId536 {
		return nil, errors.New("unknown ticket type")
	}

	rawNumbers := strings.Split(split[1], ",")
	numbers := make([]int, len(rawNumbers))
	for i, rawNumber := range rawNumbers {
		number, err := strconv.Atoi(rawNumber)
		if err != nil {
			return nil, errors.Errorf("invalid combination number: %w", err)
		}
		numbers[i] = number
	}

	if err = l.validateCombination(numbers); err != nil {
		return nil, errors.New("invalid ticket combination")
	}

	return &lottery536Ticket{
		Id:          rawTicket.Id,
		Status:      rawTicket.Status,
		DrawId:      rawTicket.DrawId,
		Combination: numbers,
	}, nil
}

// Преобразует билет из внутреннего формата в общий
func (l *Lottery536) toTicket(rawTicket *lottery536Ticket) *Ticket {
	digits := make([]string, len(rawTicket.Combination))
	for i, digit := range rawTicket.Combination {
		digits[i] = fmt.Sprintf("%02d", digit)
	}
	combination := lotteryId536 + ";" + strings.Join(digits, ",")

	return &Ticket{
		Id:     rawTicket.Id,
		Status: rawTicket.Status,
		DrawId: rawTicket.DrawId,
		Data:   base64.StdEncoding.EncodeToString([]byte(combination)),
	}
}

// Производит проверку на соответствие комбинации цифр правилам лотереи
func (l *Lottery536) validateCombination(combination []int) error {
	if len(combination) != combinationLength {
		return errors.New("invalid combination length")
	}

	uniq := map[int]struct{}{}

	// Проверяем что все числа из комбинации укладываются в допустимый диапазон
	for _, digit := range combination {
		if digit < minAllowDigit || digit > maxAllowDigit {
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
	for _, ticket := range l.Tickets {
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
