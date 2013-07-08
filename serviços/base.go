package serviços

import (
	"net/http"
	"reflect"
	"strings"
)

var (
	MuxSimples = http.NewServeMux()
	MuxSeguro  = http.NewServeMux()
)

func registrar(uri string, provedor interface{}) {
	s := &serviço{
		restrito: false,
		provedor: provedor,
	}
	MuxSimples.Handle(uri, s)
}

func registrarSeguro(uri string, provedor interface{}) {
	var restrito bool
	if len(uri) > 3 {
		if uri[:4] == "/adm" {
			if len(uri) >= 14 {
				restrito = uri[:14] != "/adm/materiais"
			} else {
				restrito = true
			}
		}
	}

	s := &serviço{
		restrito: restrito,
		provedor: provedor,
	}
	MuxSeguro.Handle(uri, s)
}

type serviço struct {
	restrito bool
	provedor interface{}
}

func (s *serviço) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.restrito {
		sessão, _ := sessões.Get(r, "coleta")

		if autenticado, ok := sessão.Values["autenticado"]; !ok || !autenticado.(bool) {
			sessão.AddFlash(r.URL.String())
			sessão.Save(r, w)
			http.Redirect(w, r, "/entrar", http.StatusTemporaryRedirect)
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
