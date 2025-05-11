package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"homework/pkg/errors"
	"math/big"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Lottery interface {
	Name() string                                                            // Название лотереи ("5 из 36", "6 из 45" и т.д.)
	Type() string                                                            // Тип лотереи (5from36, 6from45, etc...)
	Create() Lottery                                                         // Создать новый экземпляр лотереи
	AddTickets([]*Ticket) error                                              // Добавить билеты в лотерею
	AddTicketWithCombination(drawId int, combination []int) (*Ticket, error) // Создать билет с указанными номерами и добавить его в лотерею
	CreateTickets(drawId int, cost float64, num int) ([]*Ticket, error)      // Создать новые билеты с уникальными номерами и добавить их в лотерею
	Drawing(combination []int) (map[string][]*Ticket, error)                 // Провести розыгрыш (отсортировать билеты по выигрышным комбинациям)
	GenerateWinningCombination() ([]int, error)                              // Сгенерировать выигрышную комбинацию, но не применять её
}

func NewLottery5from36() Lottery {
	l := newLottery("5from36", "5 из 36", 5, 1, 36, 3)
	return &l
}

func NewLottery6from45() Lottery {
	l := newLottery("6from45", "6 из 45", 6, 1, 45, 4)
	return &l
}

type lotteryTicket struct {
	Id          int
	Status      TicketStatus
	DrawId      int
	Combination []int
	Cost        float64
}

type lottery struct {
	id                string
	name              string
	combinationLength int
	minAllowDigit     int
	maxAllowDigit     int
	minWinDigit       int

	combinations map[string]struct{}
	Tickets      []*lotteryTicket
}

func newLottery(id, name string, combinationLength, minAllowDigit, maxAllowDigit int, minWinDigit int) lottery {
	return lottery{
		id:                id,
		name:              name,
		combinationLength: combinationLength,
		minAllowDigit:     minAllowDigit,
		maxAllowDigit:     maxAllowDigit,
		minWinDigit:       minWinDigit,
	}
}

func (l *lottery) Name() string {
	return l.name
}

func (l *lottery) Type() string {
	return l.id
}

func (l *lottery) Create() Lottery {
	return &lottery{
		id:                l.id,
		name:              l.name,
		combinationLength: l.combinationLength,
		minAllowDigit:     l.minAllowDigit,
		maxAllowDigit:     l.maxAllowDigit,
		minWinDigit:       l.minWinDigit,
	}
}

// AddTickets добавляет существующие билеты в лотерею.
func (l *lottery) AddTickets(tickets []*Ticket) error {
	for _, rawTicket := range tickets {
		ticket, err := l.fromTicket(rawTicket)
		if err != nil {
			return errors.Errorf("failed convert ticket: %w", err)
		}
		sort.Ints(ticket.Combination)
		l.addTicket(ticket)
	}

	return nil
}

// AddTicketWithCombination добавляет билет с указанной комбинацией цифр в лотерею.
func (l *lottery) AddTicketWithCombination(drawId int, combination []int) (*Ticket, error) {
	if err := l.validateCombination(combination); err != nil {
		return nil, errors.New("invalid combination")
	}

	ticket := &lotteryTicket{Status: TicketStatusReady, DrawId: drawId, Combination: combination}
	l.addTicket(ticket)

	return l.toTicket(ticket), nil
}

// CreateTickets создаёт указанное количество билетов, гарантируя уникальность комбинаций.
func (l *lottery) CreateTickets(drawId int, cost float64, num int) ([]*Ticket, error) {
	newTickets := make([]*lotteryTicket, 0, num)
	maxDigitIndex := big.NewInt(int64(l.maxAllowDigit) - 1)
	for range num {

		// Создаём уникальную комбинацию чисел
		combination := make([]int, l.combinationLength)
		for l.validateCombination(combination) != nil || !l.checkUniqCombination(combination) {
			for i := range l.combinationLength {
				randomNumber, err := rand.Int(rand.Reader, maxDigitIndex)
				if err != nil {
					return nil, errors.Errorf("failed to create combination: %w", err)
				}
				combination[i] = int(randomNumber.Int64()) + 1
			}
		}
		sort.Ints(combination)
		ticket := &lotteryTicket{Status: TicketStatusReady, DrawId: drawId, Combination: combination, Cost: cost}
		l.addTicket(ticket)
		newTickets = append(newTickets, ticket)
	}

	// Конвертируем список билетов из внутреннего формата во внешний
	result := make([]*Ticket, len(newTickets))
	for i, ticket := range newTickets {
		result[i] = l.toTicket(ticket)
	}

	return result, nil
}

func (l *lottery) Drawing(combination []int) (map[string][]*Ticket, error) {
	if len(combination) != l.combinationLength {
		return nil, errors.New("invalid combination")
	}

	sort.Ints(combination)
	result := map[string][]*Ticket{}
	for _, ticket := range l.Tickets {
		matched := 0
		for i := range combination {
			if slices.Contains(ticket.Combination, combination[i]) {
				matched++
			}
		}
		if matched >= l.minWinDigit {
			value := strconv.Itoa(matched)
			result[value] = append(result[value], l.toTicket(ticket))
		}
	}

	return result, nil
}

// Добавление билета во внутреннее хранилище.
func (l *lottery) addTicket(ticket *lotteryTicket) {
	sort.Ints(ticket.Combination)
	l.Tickets = append(l.Tickets, ticket)
	key := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ticket.Combination)), ","), "[]")
	if l.combinations == nil {
		l.combinations = map[string]struct{}{}
	}
	l.combinations[key] = struct{}{}
}

// Преобразует билет из общего формата во внутренний.
func (l *lottery) fromTicket(rawTicket *Ticket) (*lotteryTicket, error) {
	data, err := base64.StdEncoding.DecodeString(rawTicket.Data)
	if err != nil {
		return nil, errors.New("unknown decode ticket data")
	}

	split := strings.SplitN(string(data), ";", 2)
	if len(split) != 2 {
		return nil, errors.New("unknown ticket format")
	}

	if len(split[1]) != l.combinationLength*3-1 {
		return nil, errors.New("invalid ticket combination length")
	}

	if split[0] != l.id {
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

	return &lotteryTicket{
		Id:          rawTicket.Id,
		Status:      rawTicket.Status,
		DrawId:      rawTicket.DrawId,
		Combination: numbers,
		Cost:        rawTicket.Cost,
	}, nil
}

// Преобразует билет из внутреннего формата в общий.
func (l *lottery) toTicket(rawTicket *lotteryTicket) *Ticket {
	digits := make([]string, len(rawTicket.Combination))
	for i, digit := range rawTicket.Combination {
		digits[i] = fmt.Sprintf("%02d", digit)
	}
	combination := l.id + ";" + strings.Join(digits, ",")

	return &Ticket{
		Id:     rawTicket.Id,
		Status: rawTicket.Status,
		DrawId: rawTicket.DrawId,
		Data:   base64.StdEncoding.EncodeToString([]byte(combination)),
		Cost:   rawTicket.Cost,
	}
}

// Производит проверку на соответствие комбинации цифр правилам лотереи.
func (l *lottery) validateCombination(combination []int) error {
	if len(combination) != l.combinationLength {
		return errors.New("invalid combination length")
	}

	uniq := map[int]struct{}{}

	// Проверяем что все числа из комбинации укладываются в допустимый диапазон
	for _, digit := range combination {
		if digit < l.minAllowDigit || digit > l.maxAllowDigit {
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
func (l *lottery) checkUniqCombination(combination []int) bool {
	// Если список сохранённых комбинаций не инициализирован, то инициализируем список и возвращаем успех
	if l.combinations == nil {
		l.combinations = map[string]struct{}{}
		return true
	}

	key := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(combination)), ","), "[]")
	_, ok := l.combinations[key]

	return !ok
}

func (l *lottery) GenerateWinningCombination() ([]int, error) {
	return GenerateUniqueRandomNumbers(l.combinationLength, l.minAllowDigit, l.maxAllowDigit)
}

// GenerateUniqueRandomNumbers генерирует указанное количество несовпадающих случайных номеров в указанном диапазоне.
func GenerateUniqueRandomNumbers(count, min, max int) ([]int, error) {
	if count > (max - min + 1) {
		return nil, errors.New("cannot generate more unique numbers than available in range")
	}
	if count <= 0 || min > max {
		return nil, errors.New("invalid count or range")
	}

	numbers := make(map[int]struct{})
	result := make([]int, 0, count)
	attempts := 0
	maxAttempts := count * 10 // Safety break

	for len(result) < count && attempts < maxAttempts {
		attempts++
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number: %w", err)
		}
		num := int(nBig.Int64()) + min

		if _, exists := numbers[num]; !exists {
			numbers[num] = struct{}{}
			result = append(result, num)
		}
	}

	if len(result) < count {
		return nil, errors.New("failed to generate sufficient unique numbers within attempts limit")
	}

	sort.Ints(result) // Keep winning numbers sorted

	return result, nil
}
