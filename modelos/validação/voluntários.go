package validação

import (
	"coleta/dao"
	"coleta/modelos"
	"log"
	"net/mail"
	"strings"
)

type VoluntárioComErros struct {
	modelos.Voluntário

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
	MsgRG                  string
	MsgCPF                 string
	MsgLíder               string
}

func NovoVoluntárioComErros(líder *modelos.Voluntário) *VoluntárioComErros {
	v := new(VoluntárioComErros)
	v.Voluntário = *líder

	return v
}

func ValidarVoluntário(líder *modelos.Voluntário) *VoluntárioComErros {
	return NovoVoluntárioComErros(líder).
		validarCamposObrigatórios().
		validarSintaxe().
		validarPolíticas().
		apurarErros()
}

func (v *VoluntárioComErros) apurarErros() *VoluntárioComErros {
	if v.errosEncontrados {
		return v
	}

	return nil
}

func (v *VoluntárioComErros) validarCamposObrigatórios() *VoluntárioComErros {
	if v.Nome == "" {
		v.errosEncontrados = true
		v.MsgNome = "Este campo não pode estar vazio"
	}

	if v.TelefoneResidencial == "" && v.TelefoneCelular == "" {
		v.errosEncontrados = true
		v.MsgContato = "Ao menos um destes campos não pode estar vazio"
	}

	if v.Email == "" {
		v.errosEncontrados = true
		v.MsgEmail = "Este campo não pode estar vazio"
	}

	if v.CPF == "" {
		v.errosEncontrados = true
		v.MsgCPF = "Este campo não pode estar vazio"
	}

	if v.Líder.Id == 0 && v.Zona.Id == 0 {
		v.errosEncontrados = true
		v.MsgLíder = "É necessário escolher entre um líder de esquina e uma zona"
	}

	if len(v.Turnos) == 0 {
		v.errosEncontrados = true
		v.MsgTurnos = "Ao menos um turno tem de estar selecionado"
	}

	return v
}

func (v *VoluntárioComErros) validarSintaxe() *VoluntárioComErros {
	if v.Email == "" {
		return v
	}

	if _, err := mail.ParseAddress(v.Email); err != nil {
		v.errosEncontrados = true
		v.MsgEmail = "Este não é um e-mail válido"
	}

	return v
}

func (v *VoluntárioComErros) validarPolíticas() *VoluntárioComErros {
	if v.Nome != "" && len(strings.Fields(v.Nome)) < 2 {
		v.errosEncontrados = true
		v.MsgNome = "Por favor, informe seu nome completo"
	}

	if v.Email == "" {
		return v
	}

	if v.Id == 0 {
		tx, err := dao.DB.Begin()
		defer tx.Commit()
		if err != nil {
			log.Println(err)
			return v
		}

		líderDAO := dao.NewLiderDAO(tx)

		if mesmoEmail, _ := líderDAO.FindByEmail(v.Email); mesmoEmail != nil {
			v.errosEncontrados = true
			v.MsgEmail = "Já existe alguém cadastrado com este mesmo e-mail. Por favor, informe outro endereço. Em caso de dúvidas, contacte seu líder de zona"
		}
	}

	return v
}
