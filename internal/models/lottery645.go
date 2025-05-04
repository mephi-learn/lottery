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
	l645id                = "6from45"
	l645combinationLength = 6
	l645minAllowDigit     = 1
	l645maxAllowDigit     = 45
)

type lottery645Ticket struct {
	Id          int
	Status      TicketStatus
	DrawId      int
	Combination []int
}
type Lottery645 struct {
	notTemplate bool
	Tickets     []*lottery645Ticket
}

func (l *Lottery645) Name() string {
	return "6 из 45"
}

func (l *Lottery645) Type() string {
	return l645id
}

func (l *Lottery645) Create() Lottery {
	return &Lottery645{notTemplate: true}
}

// AddTickets добавляет билет в лотерею.
func (l *Lottery645) AddTickets(tickets []*Ticket) error {
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

// CreateTicket создаёт билет в лотерее.
func (l *Lottery645) CreateTicket(drawId int, data string) (*Ticket, error) {
	// Не даём использовать заготовку
	if !l.notTemplate {
		return nil, errors.New("use template")
	}

	split := strings.Split(data, ",")
	combination := make([]int, len(split))
	for i, num := range split {
		add, err := strconv.Atoi(num)
		if err != nil {
			return nil, errors.Errorf("invalid combination digit: %s", num)
		}
		combination[i] = add
	}

	if err := l.validateCombination(combination); err != nil {
		return nil, errors.New("invalid combination")
	}

	return l.toTicket(&lottery645Ticket{Status: TicketStatusReady, DrawId: drawId, Combination: combination}), nil
}

func (l *Lottery645) CreateTickets(drawId int, num int) ([]*Ticket, error) {
	// Не даём использовать заготовку
	if !l.notTemplate {
		return nil, errors.New("use template")
	}

	newTickets := make([]*lottery645Ticket, 0, num)
	maxDigitIndex := big.NewInt(l645maxAllowDigit - 1)
	for range num {

		// Создаём уникальную комбинацию чисел
		combination := make([]int, l645combinationLength)
		for l.validateCombination(combination) != nil || !l.checkUniqCombination(combination, newTickets...) {
			for i := range l645combinationLength {
				randomNumber, err := rand.Int(rand.Reader, maxDigitIndex)
				if err != nil {
					return nil, errors.Errorf("failed to create combination: %w", err)
				}
				combination[i] = int(randomNumber.Int64()) + 1
			}
		}

		newTickets = append(newTickets, &lottery645Ticket{Status: TicketStatusReady, DrawId: drawId, Combination: combination})
	}

	// Конвертируем список билетов из внутреннего формата во внешний
	result := make([]*Ticket, len(newTickets))
	for i, ticket := range newTickets {
		result[i] = l.toTicket(ticket)
	}

	return result, nil
}

func (l *Lottery645) Drawing(combination []int) (map[string][]*Ticket, error) {
	if len(combination) != l645combinationLength {
		return nil, errors.New("invalid combination")
	}

	result := map[string][]*Ticket{}
	for _, ticket := range l.Tickets {
		matched := l645combinationLength
		for i := range combination {
			if combination[i] != ticket.Combination[i] {
				matched = i
				break
			}
		}
		if matched > 0 {
			value := strconv.Itoa(matched)
			result[value] = append(result[value], l.toTicket(ticket))
		}
	}

	return result, nil
}

// Преобразует билет из общего формата во внутренний.
func (l *Lottery645) fromTicket(rawTicket *Ticket) (*lottery645Ticket, error) {
	data, err := base64.StdEncoding.DecodeString(rawTicket.Data)
	if err != nil {
		return nil, errors.New("unknown decode ticket data")
	}

	split := strings.SplitN(string(data), ";", 2)
	if len(split) != 2 {
		return nil, errors.New("unknown ticket format")
	}

	if len(split[1]) != l645combinationLength*3-1 {
		return nil, errors.New("invalid ticket combination length")
	}

	if split[0] != l645id {
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

	return &lottery645Ticket{
		Id:          rawTicket.Id,
		Status:      rawTicket.Status,
		DrawId:      rawTicket.DrawId,
		Combination: numbers,
	}, nil
}

// Преобразует билет из внутреннего формата в общий.
func (l *Lottery645) toTicket(rawTicket *lottery645Ticket) *Ticket {
	digits := make([]string, len(rawTicket.Combination))
	for i, digit := range rawTicket.Combination {
		digits[i] = fmt.Sprintf("%02d", digit)
	}
	combination := l645id + ";" + strings.Join(digits, ",")

	return &Ticket{
		Id:     rawTicket.Id,
		Status: rawTicket.Status,
		DrawId: rawTicket.DrawId,
		Data:   base64.StdEncoding.EncodeToString([]byte(combination)),
	}
}

// Производит проверку на соответствие комбинации цифр правилам лотереи.
func (l *Lottery645) validateCombination(combination []int) error {
	if len(combination) != l645combinationLength {
		return errors.New("invalid combination length")
	}

	uniq := map[int]struct{}{}

	// Проверяем что все числа из комбинации укладываются в допустимый диапазон
	for _, digit := range combination {
		if digit < l645minAllowDigit || digit > l645maxAllowDigit {
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

// Проверяет, уникальна ли комбинация среди уже существующих билетов.
func (l *Lottery645) checkUniqCombination(combination []int, newTickets ...*lottery645Ticket) bool {
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
