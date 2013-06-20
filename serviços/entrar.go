package serviços

import (
	"coleta/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	registrar("/entrar", Entrar{})
}

const (
	usuário = "gustavo"
	senha   = "henrique"
)

type Entrar struct{}

func (e Entrar) Get(w http.ResponseWriter, r *http.Request) {
	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/entrar.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo entrar.html:", err)
	}

	fmt.Fprintf(w, "%s", página)
}
