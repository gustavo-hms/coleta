package serviços

import (
	"coleta/config"
	"net/http"
	"strings"
)

func init() {
	registrar("/materiais/", Materiais{})
	registrarSeguro("/adm/materiais/", Materiais{})
}

type Materiais struct{}

func (m Materiais) Get(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.String(), "/adm")
	http.ServeFile(w, r, config.Dados.DiretórioDasPáginas+path)
}
