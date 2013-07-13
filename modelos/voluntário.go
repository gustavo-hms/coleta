package modelos

import (
	"log"
	"strconv"
	"time"
)

type Voluntário struct {
	Id                  int       `json:"id"`
	Zona                *Zona     `json:"-"`
	Líder               *Líder    `json:"-"`
	Esquina             *Esquina  `json:"-"`
	Nome                string    `json:"nome"`
	TelefoneResidencial string    `json:"-"`
	TelefoneCelular     string    `json:"-"`
	Operadora           string    `json:"-"`
	RG                  string    `json:"-"`
	CPF                 string    `json:"-"`
	Idade               string    `json:"-"`
	Email               string    `json:"-"`
	Turnos              []Turno   `json:"-"`
	ComoSoube           string    `json:"-"`
	CadastradoEm        time.Time `json:"-"`
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

	if comoSoube, ok := campos["como-soube"]; ok {
		v.ComoSoube = comoSoube[0]
	}

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
