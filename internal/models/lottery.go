package models

type Lottery interface {
	Name() string                                         // Название лотереи ("5 из 36")
	Type() string                                         // Тип лотереи (5from36)
	Create() Lottery                                      // Создать экземпляр лотереи
	AddTickets([]*Ticket) error                           // Добавить билеты в лотерею
	CreateTickets(num int, drawId int) ([]*Ticket, error) // Создать новые билеты
}
