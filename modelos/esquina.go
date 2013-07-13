package modelos

import (
	"log"
	"strconv"
)

const (
	AltaPrioridade  = "alta"
	BaixaPrioridade = "baixa"
)

type Esquina struct {
	Id               int
	Cruzamento       string
	Localização      string
	Prioridade       string
	QtdDeLíderes     int
	QtdDeVoluntários int
	Zona             Zona
	Participantes    map[string]Participantes
}

type Participantes struct {
	Líderes     []Líder
	Voluntários []Voluntário
}

func (e *Esquina) AcrescentarLíder(líder Líder) {
	for _, turno := range líder.Turnos {
		participantes := e.Participantes[turno.Id]
		participantes.Líderes = append(participantes.Líderes, líder)
		e.Participantes[turno.Id] = participantes
	}
}

func (e *Esquina) AcrescentarVoluntário(voluntário Voluntário) {
	for _, turno := range voluntário.Turnos {
		participantes := e.Participantes[turno.Id]
		participantes.Voluntários = append(participantes.Voluntários, voluntário)
		e.Participantes[turno.Id] = participantes
	}
}

func (e *Esquina) Preencher(campos func(string) string) {
	e.Cruzamento = campos("cruzamento")

	id, err := strconv.Atoi(campos("zona"))
	if err != nil {
		log.Printf("Erro ao converter %s para um inteiro: %s", campos("zona"), err)
	}

	e.Zona = Zona{Id: id}
	e.Localização = campos("url")
	e.Prioridade = campos("prioridade")
}
