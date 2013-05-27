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

func (l *Líder) Preencher(campos func(string) string) *LíderValidado {
	l.Nome = campos("nome")
	l.TelefoneResidencial = campos("telefone-residencial")
	l.TelefoneCelular = campos("telefone-celular")
	l.Operadora = campos("operadora")
	l.Email = campos("e-mail")
	l.Turnos = obterTurnos(campos("turnos"))
	l.Zona = obterZona("zona")
	l.CadastradoEm = time.Now()

	return nil
}
