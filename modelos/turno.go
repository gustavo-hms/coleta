package modelos

import (
	"time"
)

var (
	Turno1Sexta  Turno
	Turno2Sexta  Turno
	Turno1Sábado Turno
	Turno2Sábado Turno

	turnos map[string]Turno
)

func init() {
	fusoHorárioDeSãoPaulo, _ := time.LoadLocation("America/Sao_Paulo")

	Turno1Sexta = Turno{
		Id:     "1",
		Início: time.Date(2013, 8, 23, 7, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 23, 13, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	Turno2Sexta = Turno{
		Id:     "2",
		Início: time.Date(2013, 8, 23, 13, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 23, 20, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	Turno1Sábado = Turno{
		Id:     "3",
		Início: time.Date(2013, 8, 24, 9, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 24, 14, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	Turno2Sábado = Turno{
		Id:     "4",
		Início: time.Date(2013, 8, 24, 14, 0, 0, 0, fusoHorárioDeSãoPaulo),
		Fim:    time.Date(2013, 8, 24, 20, 0, 0, 0, fusoHorárioDeSãoPaulo),
	}

	turnos = map[string]Turno{
		"1": Turno1Sexta,
		"2": Turno2Sexta,
		"3": Turno1Sábado,
		"4": Turno2Sábado,
	}
}

type Turno struct {
	Id     string
	Início time.Time
	Fim    time.Time
}

func (t Turno) String() string {
	var diaDaSemana string
	switch t.Início.Weekday() {
	case time.Friday:
		diaDaSemana = "SEXTA-FEIRA"
	default:
		diaDaSemana = "SÁBADO"
	}
	return t.Início.Format("2/1/2006, das 15h04 às ") + t.Fim.Format("15h04, ") + diaDaSemana
}

func Turnos() []Turno {
	return []Turno{Turno1Sexta, Turno2Sexta, Turno1Sábado, Turno2Sábado}
}

func TurnoComId(id string) Turno {
	return turnos[id]
}
