package serviços

import (
	"coleta/config"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguro("/entrar", Entrar{})
}

const (
	usuário = "nome-de-usuário"
	senha   = "senha-do-usuário"
)

type Entrar struct{}

type Autenticação struct {
	Usuário, Senha, Msg string
}

func (e Entrar) Get(w http.ResponseWriter, r *http.Request) {
	e.get(w, r, new(Autenticação))
}

func (e Entrar) get(w http.ResponseWriter, r *http.Request, entrar *Autenticação) {
	t, err := template.New("entrar").
		ParseFiles(config.Dados.DiretórioDasPáginas + "/entrar.html")
	if err != nil {
		log.Println("Ali:", err)
		return
	}

	err = t.ExecuteTemplate(w, "entrar.html", entrar)
	if err != nil {
		log.Println("Aqui:", err)
	}
}

func (e Entrar) Post(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("usuário") != usuário || r.FormValue("senha") != senha {
		autenticação := &Autenticação{
			Usuário: r.FormValue("usuário"),
			Senha:   r.FormValue("senha"),
			Msg:     "Usuário ou senha incorretos",
		}
		e.get(w, r, autenticação)
		return
	}

	sessão, _ := sessões.Get(r, "coleta")
	sessão.Values["autenticado"] = true

	mensagens := sessão.Flashes()
	sessão.Save(r, w)

	if len(mensagens) == 0 {
		http.Redirect(w, r, "/adm", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, mensagens[0].(string), http.StatusSeeOther)
	}
}
