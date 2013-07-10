package modelos

import (
	"net/mail"
)

type Voluntário struct {
	Id                  int
	Zona                *Zona
	Líder               *Líder
	Nome                string
	TelefoneResidencial string
	TelefoneCelular     string
	Operadora           string
	RG                  string
	CPF                 string
	Idade               string
	Email               mail.Address
	Turnos              []Turno
	ComoSoube           string
}
