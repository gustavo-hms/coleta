package serviços

import (
	"net/http"
	"reflect"
)

func registrarServiço(uri string, provedor interface{}) {
	s := &serviço{provedor}
	http.Handle(uri, s)
}

type serviço struct {
	provedor interface{}
}

func (s *serviço) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		s.nãoImplementado(w)
	case "GET":
		s.chamarMétodoSePossível("Get", w, r)
	case "PUT":
		s.chamarMétodoSePossível("Put", w, r)
	case "POST":
		s.chamarMétodoSePossível("Post", w, r)
	case "DELETE":
		s.chamarMétodoSePossível("Delete", w, r)
	case "PATCH":
		s.chamarMétodoSePossível("Patch", w, r)
	}
}
func (s serviço) nãoImplementado(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *serviço) chamarMétodoSePossível(name string, w http.ResponseWriter, r *http.Request) {
	method := reflect.ValueOf(s.provedor).MethodByName(name)
	if method.IsValid() {
		args := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
		method.Call(args)
	} else {
		s.nãoImplementado(w)
	}
}
