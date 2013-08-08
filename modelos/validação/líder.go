package validação

import (
	"coleta/dao"
	"coleta/modelos"
	"database/sql"
	"errors"
	"log"
	"net/mail"
	"strings"
)

var erroInesperado = errors.New("Erro inesperado")

type LíderComErros struct {
	modelos.Líder

	errosEncontrados bool
	falha            bool

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

func ValidarLíder(líder *modelos.Líder, tx *dao.Tx) (*LíderComErros, error) {
	return NovoLíderComErros(líder).
		validarCamposObrigatórios().
		validarSintaxe().
		validarPolíticas(tx).
		apurarErros()
}

func (l *LíderComErros) apurarErros() (*LíderComErros, error) {
	if l.falha {
		return nil, erroInesperado
	}

	if l.errosEncontrados {
		return l, nil
	}

	return nil, nil
}

func (l *LíderComErros) validarCamposObrigatórios() *LíderComErros {
	if l.Nome == "" {
		l.errosEncontrados = true
		l.MsgNome = "Este campo não pode estar vazio"
	}

	if l.TelefoneResidencial == "" && l.TelefoneCelular == "" {
		l.errosEncontrados = true
		l.MsgContato = "Ao menos um destes campos não pode estar vazio"
	}

	if l.Email == "" {
		l.errosEncontrados = true
		l.MsgEmail = "Este campo não pode estar vazio"
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

func (l *LíderComErros) validarPolíticas(tx *dao.Tx) *LíderComErros {
	if l.Nome != "" && len(strings.Fields(l.Nome)) < 2 {
		l.errosEncontrados = true
		l.MsgNome = "Por favor, informe seu nome completo"
	}

	if l.Email == "" {
		return l
	}

	if l.Id == 0 {
		líderDAO := dao.NewLiderDAO(tx)

		mesmoEmail, err := líderDAO.FindByEmail(l.Email)
		if err != nil && err != sql.ErrNoRows {
			l.falha = true
			log.Println(err)
		}

		if mesmoEmail != nil {
			l.errosEncontrados = true
			l.MsgEmail = "Já existe alguém cadastrado com este mesmo e-mail. Por favor, informe outro endereço. Em caso de dúvidas, contacte seu líder de zona"
		}
	}

	return l
}
