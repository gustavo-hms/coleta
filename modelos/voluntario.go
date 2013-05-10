package modelos

import (
	"net/mail"
)

type Voluntario struct {
	Id                  int
	Zona                *Zona
	Lider               *Líder
	Nome                string
	TelefoneResidencial string
	TelefoneCelular     string
	OperadoraCelular    string
	Email               mail.Address
	Turno               string
	ComoSoubeColeta2013 string
}
