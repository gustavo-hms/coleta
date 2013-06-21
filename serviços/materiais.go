package serviços

import (
	"coleta/config"
	"net/http"
)

func init() {
	registrar("/materiais/", Materiais{})
}

type Materiais struct{}

func (m Materiais) Get(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.Dados.DiretórioDasPáginas+r.URL.String())
}
