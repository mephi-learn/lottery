package models

type Lottery interface {
	Type() string                                            // Тип лотереи (5from36)
	Name() string                                            // Название лотереи ("5 из 36")
	Create() Lottery                                         // Создать экземпляр лотереи
	AddTickets(tickets []*Ticket) error                       // Добавить билеты в лотерею
	CreateTickets(drawId int, num int) ([]*Ticket, error)    // Создать новые билеты
	Drawing(combination []int) (map[string][]*Ticket, error) // Провести розыгрыш
}
