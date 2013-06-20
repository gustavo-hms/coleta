package serviços

import (
	"coleta/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	registrar("/adm", Adm{})
}

type Adm struct{}

func (a Adm) Get(w http.ResponseWriter, r *http.Request) {
	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/adm.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo adm.html:", err)
	}

	fmt.Fprintf(w, "%s", página)
}
