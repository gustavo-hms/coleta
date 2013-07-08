package serviços

import (
	"coleta/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	registrarSeguro("/adm", Adm{})
	registrar("/adm/", Redirecionamento{})
}

type Redirecionamento struct{}

func (_ Redirecionamento) Get(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://tetocoleta.com.br"+r.URL.String(), http.StatusSeeOther)
}

type Adm struct{}

func (a Adm) Get(w http.ResponseWriter, r *http.Request) {
	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/adm.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo adm.html:", err)
	}

	fmt.Fprintf(w, "%s", página)
}
