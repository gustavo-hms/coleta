package serviços

import (
	"net/http"
	"reflect"
	"strings"
)

func registrar(uri string, provedor interface{}) {
	s := &serviço{provedor}
	http.Handle(uri, s)
}

type serviço struct {
	provedor interface{}
}

func (s *serviço) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
