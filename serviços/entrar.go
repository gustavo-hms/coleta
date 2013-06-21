package serviços

import (
	"coleta/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	registrarSeguro("/entrar", Entrar{})
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

func (e Entrar) Post(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("usuário") != usuário || r.FormValue("senha") != senha {
		e.Get(w, r) // TODO escrever mensagem de erro no template
		return
	}

	sessão, _ := sessões.Get(r, "coleta")
	sessão.Values["autenticado"] = true

	mensagens := sessão.Flashes()
	sessão.Save(r, w)

	if len(mensagens) == 0 {
		http.Redirect(w, r, "/adm", http.StatusTemporaryRedirect)

	} else {
		http.Redirect(w, r, mensagens[0].(string), http.StatusTemporaryRedirect)
	}
}
