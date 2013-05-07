package modelos

import (
	"net/http"
	"net/url"
)

type Esquina struct {
	Cruzamento  string
	Localização url.URL
	Zona        Zona
}

func (e *Esquina) Preencher(req *http.Request) {
	e.Cruzamento = req.FormValue("cruzamento")

	u, _ := url.Parse(req.FormValue("url"))
	e.Localização = *u

	e.Zona = Zona(req.FormValue("zona"))
}

type EsquinaValidada struct {
	Esquina

	MsgCruzamento  string
	MsgLocalização string
	MsgZona        string
}
