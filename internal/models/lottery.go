package models

import "io"

type Lottery interface {
	Name() string                                            // Название лотереи ("5 из 36")
	Type() string                                            // Тип лотереи (5from36)
	Create() Lottery                                         // Создать экземпляр лотереи
	AddTickets([]*Ticket) error                              // Добавить билеты в лотерею
	CreateTickets(drawId int, num int) ([]*Ticket, error)    // Создать новые билеты
	Drawing(combination []int) (map[string][]*Ticket, error) // Провести розыгрыш
	GenerateWinningCombination(randReader io.Reader) ([]int, error)  // Сгенерировать выигрышную комбинацию
}
