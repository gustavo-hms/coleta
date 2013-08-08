package serviços

import (
	"coleta/dao"
	"log"
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

func registrarComTransação(uri string, provedor interface{}) {
	s := &serviçoComTransação{
		restrito: false,
		provedor: provedor,
	}
	MuxSimples.Handle(uri, s)
}

func registrarSeguroComTransação(uri string, provedor interface{}) {
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

	s := &serviçoComTransação{
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

type serviçoComTransação struct {
	restrito bool
	provedor interface{}
	erro     bool
}

func (s *serviçoComTransação) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.restrito {
		sessão, _ := sessões.Get(r, "coleta")

		if autenticado, ok := sessão.Values["autenticado"]; !ok || !autenticado.(bool) {
			sessão.AddFlash(r.URL.String())
			sessão.Save(r, w)
			http.Redirect(w, r, "/entrar", http.StatusTemporaryRedirect)
			return
		}
	}

	tx := s.préTratamento(w, r)
	if tx == nil {
		return
	}

	nome := r.Method[0:1] + strings.ToLower(r.Method[1:])
	s.chamarMétodoSePossível(nome, w, r, tx)

	s.pósTratamento(w, r, tx)
}

func (s serviçoComTransação) préTratamento(w http.ResponseWriter, r *http.Request) *dao.Tx {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return nil
	}

	return &dao.Tx{tx}
}

func (s serviçoComTransação) pósTratamento(w http.ResponseWriter, r *http.Request, tx *dao.Tx) {
	if s.erro {
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			erroInterno(w, r)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		erroInterno(w, r)

		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
	}
}

func (s serviçoComTransação) nãoImplementado(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *serviçoComTransação) chamarMétodoSePossível(nome string, w http.ResponseWriter, r *http.Request, tx *dao.Tx) {
	método := reflect.ValueOf(s.provedor).MethodByName(nome)
	if método.IsValid() {
		args := []reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			reflect.ValueOf(tx),
		}

		erros := método.Call(args)
		if len(erros) != 1 {
			log.Printf("Slice de erros tem %d elementos quando deveria ter 1", len(erros))
			s.erro = true
			return
		}

		s.erro = !erros[0].IsNil()

	} else {
		s.erro = true
		s.nãoImplementado(w)
	}
}
