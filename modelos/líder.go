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

type LíderValidado struct {
	Líder

	MsgNome                string
	MsgTelefoneResidencial string
	MsgTelefoneCelular     string
	MsgOperadora           string
	MsgEmail               string
	MsgTurnos              string
	MsgZona                string
	MsgCadastradoEm        string
	MsgEsquina             string
}

func NovoLíderValidado() *LíderValidado {
	l := new(LíderValidado)
	l.Líder = *NovoLíder()

	return l
}

func (l *Líder) Preencher(campos map[string][]string) *LíderValidado {
	l.Nome = campos["nome"][0]
	l.TelefoneResidencial = campos["telefone-residencial"][0]
	l.TelefoneCelular = campos["telefone-celular"][0]
	l.Operadora = campos["operadora"][0]
	l.Email = campos["e-mail"][0]
	l.Turnos = obterTurnos(campos["turnos"])

	id, err := strconv.Atoi(campos["zona"][0])
	if err != nil {
		log.Printf("Erro ao converter %s para um inteiro: %s", campos["zona"][0], err)
	}

	l.Zona = &Zona{Id: id}

	l.CadastradoEm = time.Now()

	return nil
}

func obterTurnos(ids []string) []Turno {
	turnos := make([]Turno, len(ids))
	for i, id := range ids {
		turnos[i] = TurnoComId(id)
	}

	return turnos
}
