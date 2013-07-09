package modelos

import (
	"log"
	"strconv"
	"time"
)

const (
	OperadoraClaro = "Claro"
	OperadoraOi    = "Oi"
	OperadoraOutra = "Outra"
	OperadoraTim   = "TIM"
	OperadoraVivo  = "Vivo"
)

type Líder struct {
	Id                  int
	Nome                string
	TelefoneResidencial string
	TelefoneCelular     string
	Operadora           string
	Email               string
	Turnos              []Turno
	Zona                *Zona
	CadastradoEm        time.Time
	Esquina             *Esquina
}

func NovoLíder() *Líder {
	l := new(Líder)
	l.Zona = new(Zona)
	l.Esquina = new(Esquina)

	return l
}

func (l *Líder) Preencher(campos map[string][]string) {
	l.Nome = campos["nome"][0]
	l.TelefoneResidencial = campos["telefone-residencial"][0]
	l.TelefoneCelular = campos["telefone-celular"][0]

	if operadora, ok := campos["operadora"]; ok {
		l.Operadora = operadora[0]
	} else {
		l.Operadora = OperadoraOutra
	}

	l.Email = campos["e-mail"][0]
	l.Turnos = obterTurnos(campos["turnos"])

	id, err := strconv.Atoi(campos["zona"][0])
	if err != nil {
		log.Printf("Erro ao converter %s para um inteiro: %s", campos["zona"][0], err)
	}

	l.Zona = &Zona{Id: id}

	if esquina, ok := campos["esquina"]; ok {
		id, err := strconv.Atoi(esquina[0])
		if err != nil {
			log.Printf("Erro ao converter %s para um inteiro: %s", esquina[0], err)
		}

		l.Esquina = &Esquina{Id: id}
	}

	l.CadastradoEm = time.Now().UTC()
}

func obterTurnos(ids []string) []Turno {
	turnos := make([]Turno, len(ids))
	for i, id := range ids {
		turnos[i] = TurnoComId(id)
	}

	return turnos
}
