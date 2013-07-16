package serviços

import (
	"coleta/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	registrar("/", Raiz{})
	registrarSeguro("/", RedirecionamentoRaiz{})
}

type RedirecionamentoRaiz struct{}

func (_ RedirecionamentoRaiz) Get(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://tetocoleta.com.br", http.StatusSeeOther)
}

type Raiz struct{}

func (Raiz) Get(w http.ResponseWriter, r *http.Request) {
	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/index.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo index.html:", err)
		erroInterno(w, r)
		return
	}

	fmt.Fprintf(w, "%s", página)
}
