package modelos

import (
	"log"
	"strconv"
	"time"
)

type Voluntário struct {
	Id                  int
	Zona                *Zona
	Líder               *Líder
	Esquina             *Esquina
	Nome                string
	TelefoneResidencial string
	TelefoneCelular     string
	Operadora           string
	RG                  string
	CPF                 string
	Idade               string
	Email               string
	Turnos              []Turno
	ComoSoube           string
	CadastradoEm        time.Time
}

func NovoVoluntário() *Voluntário {
	return &Voluntário{
		Zona:    new(Zona),
		Líder:   new(Líder),
		Esquina: new(Esquina),
	}
}

func (v *Voluntário) Preencher(campos map[string][]string) {
	v.Nome = campos["nome"][0]
	v.TelefoneResidencial = campos["telefone-residencial"][0]
	v.TelefoneCelular = campos["telefone-celular"][0]
	v.Idade = campos["idade"][0]
	v.ComoSoube = campos["como-soube"][0]

	if operadora, ok := campos["operadora"]; ok {
		v.Operadora = operadora[0]
	} else {
		v.Operadora = OperadoraOutra
	}

	if rg, ok := campos["rg"]; ok {
		v.RG = rg[0]
	}

	if cpf, ok := campos["cpf"]; ok {
		v.CPF = cpf[0]
	}

	v.Email = campos["e-mail"][0]
	v.Turnos = obterTurnos(campos["turnos"])

	if líder, ok := campos["lider"]; ok {
		id, err := strconv.Atoi(líder[0])
		if err != nil {
			log.Printf("Erro ao converter %s para um inteiro: %s", líder[0], err)
		}

		v.Líder.Id = id

	} else {
		id, err := strconv.Atoi(campos["zona"][0])
		if err != nil {
			log.Printf("Erro ao converter %s para um inteiro: %s", campos["zona"][0], err)
		}

		v.Zona.Id = id
	}

	if esquina, ok := campos["esquina"]; ok {
		id, err := strconv.Atoi(esquina[0])
		if err != nil {
			log.Printf("Erro ao converter %s para um inteiro: %s", esquina[0], err)
		}

		v.Esquina = &Esquina{Id: id}
	}

	v.CadastradoEm = time.Now().UTC()
}
