package models

type Lottery interface {
	Name() string // Название лотереи ("5 из 36")
	Type() string // Тип лотереи (5from36)
	String() string
	New() Lottery
}
