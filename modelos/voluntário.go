package modelos

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
}

func NovoVoluntário() *Voluntário {
	return &Voluntário{
		Zona:    new(Zona),
		Líder:   new(Líder),
		Esquina: new(Esquina),
	}
}
