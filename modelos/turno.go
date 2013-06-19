package modelos

import (
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
		Início: time.Date(2013, 8, 23, 7, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 23, 13, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turno2Sexta = Turno{
		Id:     "2",
		Início: time.Date(2013, 8, 23, 13, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 23, 20, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turno1Sábado = Turno{
		Id:     "3",
		Início: time.Date(2013, 8, 24, 9, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 24, 14, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turno2Sábado = Turno{
		Id:     "4",
		Início: time.Date(2013, 8, 24, 14, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 24, 20, 0, 0, 0, fusoHorárioDeSãoPaulo),
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
	return t.Início.Format("2/1/2006, das 15h04 às ") + t.Fim.Format("15h04")
}

func Turnos() []Turno {
	return []Turno{turno1Sexta, turno2Sexta, turno1Sábado, turno2Sábado}
}

func TurnoComId(id string) Turno {
	return turnos[id]
}
