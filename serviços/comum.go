package serviços

import (
	"coleta/modelos"
	"net/http"
)

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
