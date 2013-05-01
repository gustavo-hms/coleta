package modelos

import (
	"net/mail"
	"time"
)

type Líder struct {
	Nome                string
	TelefoneResidencial string
	TelefoneCelular     string
	Operadora           int
	Email               mail.Address
	Turnos              []Turno
	Zona                Zona
	CadastradoEm        time.Time
	Esquina             Esquina
}
