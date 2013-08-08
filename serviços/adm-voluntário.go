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
	registrarSeguroComTransação("/adm/voluntario/", AdmVoluntário{})
}

type AdmVoluntário struct{}

func idDoVoluntário(endereço *url.URL) string {
	idx := strings.LastIndex(endereço.Path, "/")
	return endereço.Path[idx+1:]
}

func (l AdmVoluntário) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	stringDoId := idDoVoluntário(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro: %s", stringDoId, err)
		return err
	}

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntário, err := voluntárioDAO.FindById(id)
	if err != nil {
		log.Printf("Erro ao carregar voluntário com id %d: %s", id, err)
		return err
	}

	if voluntário.Líder.Id != 0 {
		líderDAO := dao.NewLiderDAO(tx)
		líder, err := líderDAO.FindById(voluntário.Líder.Id)
		if err != nil {
			log.Printf("Erro ao carregar líder com id %d: %s", voluntário.Líder.Id, err)
			return err
		}

		voluntário.Líder = líder
	}

	return l.get(w, r, tx, &validação.VoluntárioComErros{Voluntário: *voluntário})
}

func (l AdmVoluntário) get(
	w http.ResponseWriter,
	r *http.Request,
	tx *dao.Tx,
	voluntário *validação.VoluntárioComErros,
) error {
	t := exibiçãoDoVoluntário(voluntário, tx, "adm-voluntário.html")

	if t == nil {
		return erroInesperado
	}

	err := t.ExecuteTemplate(w, "adm-voluntário.html", voluntário)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (l AdmVoluntário) Post(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		return err
	}

	stringDoId := idDoVoluntário(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro:", stringDoId, err)
		return err
	}

	voluntário := modelos.NovoVoluntário()
	voluntário.Id = id
	voluntário.Preencher(r.Form)
	erros := validação.ValidarVoluntário(voluntário)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		return l.get(w, r, tx, erros)
	}

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	if err := voluntárioDAO.Save(voluntário); err != nil {
		log.Println("Erro ao gravar voluntário:", err)
		return err
	}

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/atualização-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo voluntários-sucesso.html:", err)
		return err
	}

	fmt.Fprintf(w, "%s", página)

	return nil
}
