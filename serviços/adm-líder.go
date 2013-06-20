package serviços

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"coleta/modelos/validação"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func init() {
	registrar("/adm/líder/", AdmLíder{})
}

type AdmLíder struct{}

func idDoLíder(endereço *url.URL) string {
	return strings.SplitN(endereço.Path[12:], "/", 2)[0]
}

func (l AdmLíder) Get(w http.ResponseWriter, r *http.Request) {
	stringDoId := idDoLíder(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro: %s", stringDoId, err)
		erroInterno(w, r)
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	líder, err := líderDAO.FindById(id)
	if err != nil {
		líderDAO.Rollback()
		log.Printf("Erro ao carregar líder com id %d: %s", id, err)
		erroInterno(w, r)
		return
	}

	líderDAO.Commit()

	l.get(w, r, &validação.LíderComErros{Líder: *líder})
}

func (l AdmLíder) get(
	w http.ResponseWriter,
	r *http.Request,
	líder *validação.LíderComErros,
) {
	t := exibiçãoDoLíder(líder, "adm-líder.html")
	if t != nil {
		err := t.ExecuteTemplate(w, "adm-líder.html", líder)
		if err != nil {
			log.Println("Aqui:", err)
			erroInterno(w, r)
			return
		}
	}
}

func (l AdmLíder) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		erroInterno(w, r)
		return
	}

	líder := modelos.NovoLíder()
	líder.Preencher(r.Form)
	erros := validação.ValidarLíder(líder)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.get(w, r, erros)
		return
	}

	stringDoId := idDoLíder(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro:", stringDoId, err)
		erroInterno(w, r)
		return
	}

	líder.Id = id

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	if err := líderDAO.Save(líder); err != nil {
		líderDAO.Rollback()
		log.Println("Erro ao gravar líder:", err)
		erroInterno(w, r)
		return
	}
	if err := líderDAO.Commit(); err != nil {
		líderDAO.Rollback()
		log.Println("Erro no commit:", err)
		erroInterno(w, r)
		return
	}

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/atualização-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo líderes-sucesso.html:", err)
		erroInterno(w, r)
		return
	}

	fmt.Fprintf(w, "%s", página)
}
