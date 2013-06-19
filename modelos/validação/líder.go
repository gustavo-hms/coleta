package validação

import (
	"coleta/dao"
	"coleta/modelos"
	"log"
	"net/mail"
)

type LíderComErros struct {
	modelos.Líder

	errosEncontrados bool

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

func NovoLíderComErros(líder *modelos.Líder) *LíderComErros {
	l := new(LíderComErros)
	l.Líder = *líder

	return l
}

func ValidarLíder(líder *modelos.Líder) *LíderComErros {
	return NovoLíderComErros(líder).
		validarCamposObrigatórios().
		validarSintaxe().
		validarPolíticas().
		apurarErros()
}

func (l *LíderComErros) apurarErros() *LíderComErros {
	if l.errosEncontrados {
		return l
	}

	return nil
}

func (l *LíderComErros) validarCamposObrigatórios() *LíderComErros {
	if l.Nome == "" {
		l.errosEncontrados = true
		l.MsgNome = "Este campo não pode estar vazio"
	}

	if l.TelefoneResidencial == "" && l.TelefoneCelular == "" && l.Email == "" {
		l.errosEncontrados = true
		l.MsgContato = "Ao menos um destes campos não pode estar vazio"
	}

	if len(l.Turnos) == 0 {
		l.errosEncontrados = true
		l.MsgTurnos = "Ao menos um turno tem de estar selecionado"
	}

	return l
}

func (l *LíderComErros) validarSintaxe() *LíderComErros {
	if l.Email == "" {
		return l
	}

	if _, err := mail.ParseAddress(l.Email); err != nil {
		l.errosEncontrados = true
		l.MsgEmail = "Este não é um e-mail válido"
	}

	return l
}

func (l *LíderComErros) validarPolíticas() *LíderComErros {
	if l.Email == "" {
		return l
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		return l
	}

	líderDAO := dao.NewLiderDAO(tx)

	if mesmoEmail, _ := líderDAO.FindByEmail(l.Email); mesmoEmail != nil {
		l.errosEncontrados = true
		l.MsgEmail = "Já existe alguém cadastrado com esse mesmo e-mail. Por favor, informe outro endereço. Em caso de dúvidas, contacte seu líder de zona"
	}

	return l
}
