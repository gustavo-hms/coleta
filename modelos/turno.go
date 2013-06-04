package modelos

import (
	"fmt"
	"time"
)

var (
	turno1Sexta  Turno
	turno2Sexta  Turno
	turno1Sábado Turno
	turno2Sábado Turno

	turnos map[string]Turno
)

func init() {
	fusoHorárioDeSãoPaulo, _ := time.LoadLocation("America/Sao_Paulo")

	turno1Sexta = Turno{
		Id:     "1",
		Início: time.Date(2013, 8, 23, 8, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 23, 12, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turno2Sexta = Turno{
		Id:     "2",
		Início: time.Date(2013, 8, 23, 13, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 23, 18, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turno1Sábado = Turno{
		Id:     "3",
		Início: time.Date(2013, 8, 24, 8, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 24, 12, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turno2Sábado = Turno{
		Id:     "4",
		Início: time.Date(2013, 8, 24, 13, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 24, 18, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turnos = map[string]Turno{
		"1": turno1Sexta,
		"2": turno2Sexta,
		"3": turno1Sábado,
		"4": turno2Sábado,
	}
}

type Turno struct {
	Id     string
	Início time.Time
	Fim    time.Time
}

func (t Turno) String() string {
	h, m, _ := t.Fim.Clock()
	return t.Início.Format("2/1/2006, das 15h04 às ") + fmt.Sprintf("%02dh%02d", h, m)
}

func Turnos() []Turno {
	return []Turno{turno1Sexta, turno2Sexta, turno1Sábado, turno2Sábado}
}

func TurnoComId(id string) Turno {
	return turnos[id]
}
