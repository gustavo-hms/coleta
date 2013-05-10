package modelos

import (
	"net/http"
	"net/url"
)

type Esquina struct {
	Id          int
	Cruzamento  string
	Localização string
	Zona        Zona
}

func (e *Esquina) Preencher(req *http.Request) *EsquinaValidada {
	e.Cruzamento = req.FormValue("cruzamento")
	e.Zona = Zona{Nome: req.FormValue("zona")}
	e.Localização = req.FormValue("url")

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

	if e.Zona.Nome == "" {
		comErros = e.preencherEsquinaValidada(comErros)
		comErros.MsgZona = "Este campo não pode estar vazio"
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

type EsquinaValidada struct {
	Esquina

	MsgCruzamento  string
	MsgLocalização string
	MsgZona        string
}
