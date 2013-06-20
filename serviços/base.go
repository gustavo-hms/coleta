package serviços

import (
	"net/http"
	"reflect"
	"strings"
)

func registrar(uri string, provedor interface{}) {
	s := &serviço{
		restrito: len(uri) > 4 && uri[:5] == "/adm/",
		provedor: provedor,
	}
	http.Handle(uri, s)
}

type serviço struct {
	restrito bool
	provedor interface{}
}

func (s *serviço) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.restrito {
		sessão, _ := sessões.Get(r, "coleta")
		if autenticado, ok := sessão.Values["autenticado"]; !ok || !autenticado.(bool) {
			sessão.Values["origem"] = r.URL.String()
			sessão.Save(r, w)
			http.Redirect(w, r, "entrar", http.StatusTemporaryRedirect)
			return
		}
	}

	nome := r.Method[0:1] + strings.ToLower(r.Method[1:])
	s.chamarMétodoSePossível(nome, w, r)
}
func (s serviço) nãoImplementado(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *serviço) chamarMétodoSePossível(nome string, w http.ResponseWriter, r *http.Request) {
	método := reflect.ValueOf(s.provedor).MethodByName(nome)
	if método.IsValid() {
		args := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
		método.Call(args)
	} else {
		s.nãoImplementado(w)
	}
}
