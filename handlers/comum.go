package handlers

import (
	"coleta/modelos"
	"net/http"
	"os"
)

var gopath = os.Getenv("GOPATH") // NOTA Solução temporária. Apenas para testes

type zonaComSeleção struct {
	Zona        modelos.Zona
	Selecionado bool
}

type turnoComSeleção struct {
	Turno       modelos.Turno
	Selecionado bool
}

func erroInterno(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}
