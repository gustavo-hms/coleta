package modelos

import (
	"net/mail"
	"time"
)

const (
	OperadoraClaro = "CLARO"
	OperadoraOi    = "OI"
	OperadoraOutra = "OUTRA"
	OperadoraTim   = "TIM"
	OperadoraVivo  = "VIVO"
)

type Líder struct {
	Id                  int
	Nome                string
	TelefoneResidencial string
	TelefoneCelular     string
	Operadora           string
	Email               mail.Address
	Turnos              []Turno
	Zona                Zona
	CadastradoEm        time.Time
	Esquina             Esquina
}
