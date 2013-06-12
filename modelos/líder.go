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

type LíderComErros struct {
	Líder

	MsgNome                string
	MsgTelefoneResidencial string
	MsgTelefoneCelular     string
	MsgOperadora           string
	MsgEmail               string
	MsgContato             string
	MsgTurnos              string
	MsgZona                string
	MsgCadastradoEm        string
	MsgEsquina             string
}

func NovoLíderComErros() *LíderComErros {
	l := new(LíderComErros)
	l.Líder = *NovoLíder()

	return l
}

func (l *Líder) Preencher(campos map[string][]string) *LíderComErros {
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

	return l.validar()
}

func obterTurnos(ids []string) []Turno {
	turnos := make([]Turno, len(ids))
	for i, id := range ids {
		turnos[i] = TurnoComId(id)
	}

	return turnos
}

func (l *Líder) validar() (comErros *LíderComErros) {
	comErros = l.validarCamposObrigatórios(comErros)
	comErros = l.validarSintaxe(comErros)

	return comErros
}

func (l *Líder) validarCamposObrigatórios(comErros *LíderComErros) *LíderComErros {
	if l.Nome == "" {
		comErros = l.preencherLíderComErros(comErros)
		comErros.MsgNome = "Este campo não pode estar vazio"
	}

	if l.TelefoneResidencial == "" && l.TelefoneCelular == "" && l.Email == "" {
		comErros = l.preencherLíderComErros(comErros)
		comErros.MsgContato = "Ao menos um destes campos não pode estar vazio"
	}

	if len(l.Turnos) == 0 {
		comErros = l.preencherLíderComErros(comErros)
		comErros.MsgTurnos = "Ao menos um turno tem de estar selecionado"
	}

	return comErros
}

func (l *Líder) preencherLíderComErros(validado *LíderComErros) *LíderComErros {
	if validado == nil {
		validado = new(LíderComErros)
		validado.Líder = *l
	}

	return validado
}

func (l *Líder) validarSintaxe(comErros *LíderComErros) *LíderComErros {
	return comErros
}
