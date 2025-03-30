package monkey

import "fmt"

type Printer interface {
	Println(a ...any)
	Printf(format string, a ...any)
}

type printer struct{}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Println(a ...any) {
	fmt.Println(a...)
}

func (p *printer) Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}
