package serviços

import (
	"coleta/modelos"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)

var sessões = sessions.NewCookieStore(
	securecookie.GenerateRandomKey(32),
	securecookie.GenerateRandomKey(32),
)

func init() {
	sessões.Options.MaxAge = 60 * 30
}

type zonaComSeleção struct {
	Zona        modelos.Zona
	Selecionado bool
}

type esquinaComSeleção struct {
	Esquina     modelos.Esquina
	Selecionado bool
}

type turnoComSeleção struct {
	Turno       modelos.Turno
	Selecionado bool
}

func erroInterno(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}
