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
		restrito: len(uri) > 3 && uri[:4] == "/adm",
		provedor: provedor,
	}
	MuxSimples.Handle(uri, s)
}

func registrarSeguro(uri string, provedor interface{}) {
	s := &serviço{
		restrito: len(uri) > 3 && uri[:4] == "/adm",
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

	print("Método: ", r.Method, "\n")
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
