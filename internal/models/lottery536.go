package models

type Lottery536 struct {
}

func (l *Lottery536) Name() string {
	return "5 из 36"
}

func (l *Lottery536) Type() string {
	return "5from36"
}

func (l *Lottery536) String() string {
	return l.Type()
}

func (l *Lottery536) New() Lottery {
	return &Lottery536{}
}
