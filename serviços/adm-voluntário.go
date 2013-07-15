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
	registrarSeguro("/adm/voluntario/", AdmVoluntário{})
}

type AdmVoluntário struct{}

func idDoVoluntário(endereço *url.URL) string {
	idx := strings.LastIndex(endereço.Path, "/")
	return endereço.Path[idx+1:]
}

func (l AdmVoluntário) Get(w http.ResponseWriter, r *http.Request) {
	stringDoId := idDoVoluntário(r.URL)
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

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntário, err := voluntárioDAO.FindById(id)
	if err != nil {
		voluntárioDAO.Rollback()
		log.Printf("Erro ao carregar voluntário com id %d: %s", id, err)
		erroInterno(w, r)
		return
	}

	if voluntário.Líder.Id != 0 {
		líderDAO := dao.NewLiderDAO(tx)
		líder, err := líderDAO.FindById(voluntário.Líder.Id)
		if err != nil {
			tx.Rollback()
			log.Printf("Erro ao carregar líder com id %d: %s", voluntário.Líder.Id, err)
			erroInterno(w, r)
			return
		}

		voluntário.Líder = líder
	}

	tx.Commit()

	l.get(w, r, &validação.VoluntárioComErros{Voluntário: *voluntário})
}

func (l AdmVoluntário) get(
	w http.ResponseWriter,
	r *http.Request,
	voluntário *validação.VoluntárioComErros,
) {
	t := exibiçãoDoVoluntário(voluntário, "adm-voluntário.html")
	if t != nil {
		err := t.ExecuteTemplate(w, "adm-voluntário.html", voluntário)
		if err != nil {
			log.Println(err)
			erroInterno(w, r)
			return
		}
	}
}

func (l AdmVoluntário) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		erroInterno(w, r)
		return
	}

	stringDoId := idDoVoluntário(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro:", stringDoId, err)
		erroInterno(w, r)
		return
	}

	voluntário := modelos.NovoVoluntário()
	voluntário.Id = id
	voluntário.Preencher(r.Form)
	erros := validação.ValidarVoluntário(voluntário)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.get(w, r, erros)
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	if err := voluntárioDAO.Save(voluntário); err != nil {
		voluntárioDAO.Rollback()
		log.Println("Erro ao gravar voluntário:", err)
		erroInterno(w, r)
		return
	}
	if err := voluntárioDAO.Commit(); err != nil {
		voluntárioDAO.Rollback()
		log.Println("Erro no commit:", err)
		erroInterno(w, r)
		return
	}

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/atualização-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo voluntários-sucesso.html:", err)
		erroInterno(w, r)
		return
	}

	fmt.Fprintf(w, "%s", página)
}
