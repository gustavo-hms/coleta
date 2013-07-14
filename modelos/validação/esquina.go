package validação

import (
	"coleta/modelos"
	"net/url"
)

type EsquinaComErros struct {
	modelos.Esquina

	errosEncontrados bool

	MsgCruzamento  string
	MsgLocalização string
	MsgZona        string
}

func NovaEsquinaComErros(esquina *modelos.Esquina) *EsquinaComErros {
	e := new(EsquinaComErros)
	e.Esquina = *esquina

	return e
}

func ValidarEsquina(esquina *modelos.Esquina) *EsquinaComErros {
	return NovaEsquinaComErros(esquina).
		validarCamposObrigatórios().apurarErros()
}

func (e *EsquinaComErros) apurarErros() *EsquinaComErros {
	if e.errosEncontrados {
		return e
	}

	return nil
}

func (e *EsquinaComErros) validarCamposObrigatórios() *EsquinaComErros {
	if e.Cruzamento == "" {
		e.errosEncontrados = true
		e.MsgCruzamento = "Este campo não pode estar vazio"
	}

	if e.Zona.Id < 0 {
		e.errosEncontrados = true
		e.MsgZona = "Este campo está com um valor inválido"
	}

	return e
}

func (e *EsquinaComErros) validarSintaxe() *EsquinaComErros {
	if e.MsgLocalização == "" {
		u, err := url.ParseRequestURI(e.Localização)
		if err != nil {
			e.errosEncontrados = true
			e.MsgLocalização = "Endereço inválido"
		} else {
			e.Localização = u.String()
		}
	}

	return e
}
