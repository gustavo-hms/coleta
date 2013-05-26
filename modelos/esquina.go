package modelos

import (
	"log"
	"net/url"
	"strconv"
)

type Esquina struct {
	Id          int
	Cruzamento  string
	Localização string
	Zona        Zona
}

type EsquinaValidada struct {
	Esquina

	MsgCruzamento  string
	MsgLocalização string
	MsgZona        string
}

func (e *Esquina) Preencher(campos func(string) string) *EsquinaValidada {
	e.Cruzamento = campos("cruzamento")

	id, err := strconv.Atoi(campos("zona"))
	if err != nil {
		log.Printf("Erro ao converter %s para um inteiro: %s", campos("zona"), err)
	}

	e.Zona = Zona{Id: id}
	e.Localização = campos("url")

	return e.Validar()
}

func (e *Esquina) Validar() (comErros *EsquinaValidada) {
	comErros = e.validarCamposObrigatórios(comErros)
	comErros = e.validarSintaxe(comErros)

	return comErros
}

func (e *Esquina) validarCamposObrigatórios(comErros *EsquinaValidada) *EsquinaValidada {
	if e.Cruzamento == "" {
		comErros = e.preencherEsquinaValidada(comErros)
		comErros.MsgCruzamento = "Este campo não pode estar vazio"
	}

	if e.Localização == "" {
		comErros = e.preencherEsquinaValidada(comErros)
		comErros.MsgLocalização = "Este campo não pode estar vazio"
	}

	if e.Zona.Id < 0 {
		comErros = e.preencherEsquinaValidada(comErros)
		comErros.MsgZona = "Este campo está com um valor inválido"
	}

	return comErros
}

func (e *Esquina) validarSintaxe(comErros *EsquinaValidada) *EsquinaValidada {
	if comErros == nil || comErros.MsgLocalização == "" {
		u, err := url.ParseRequestURI(e.Localização)
		if err != nil {
			comErros = e.preencherEsquinaValidada(comErros)
			comErros.MsgLocalização = "Endereço inválido"
		} else {
			e.Localização = u.String()
		}
	}

	return comErros
}

func (e *Esquina) preencherEsquinaValidada(validada *EsquinaValidada) *EsquinaValidada {
	if validada == nil {
		validada = new(EsquinaValidada)
		validada.Esquina = *e
	}

	return validada
}
