package modelos

import (
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

func (l *Líder) Preencher(campos map[string][]string) *LíderValidado {
	l.Nome = campos["nome"][0]
	l.TelefoneResidencial = campos["telefone-residencial"][0]
	l.TelefoneCelular = campos["telefone-celular"][0]
	l.Operadora = campos["operadora"][0]
	l.Email = campos["e-mail"][0]
	l.Turnos = obterTurnos(campos["turnos"])
	//	l.Zona = obterZona("zona")
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
