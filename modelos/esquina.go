package modelos

import (
	"log"
	"strconv"
)

type Esquina struct {
	Id               int
	Cruzamento       string
	Localização      string
	Prioridade       string
	QtdDeLíderes     int
	QtdDeVoluntários int
	Zona             Zona
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
