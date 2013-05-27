package modelos

import (
	"time"
)

var (
	Turno1Sexta  Turno
	Turno2Sexta  Turno
	Turno1Sábado Turno
	Turno2Sábado Turno
)

func init() {
	fusoHorárioDeSãoPaulo, _ := time.LoadLocation("America/Sao_Paulo")

	Turno1Sexta = Turno{
		Id:     1,
		Início: time.Date(2013, 8, 23, 8, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 23, 12, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}
}

type Turno struct {
	Id     int
	Início time.Time
	Fim    time.Time
}
