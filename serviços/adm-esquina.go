package serviços

import (
	"coleta/config"
	"coleta/dao"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	registrarSeguro("/adm/esquina/", new(AdmEsquina))
}

type AdmEsquina struct {
}

func idDaEsquina(endereço *url.URL) string {
	idx := strings.LastIndex(endereço.Path, "/")
	return endereço.Path[idx+1:]
}

func (e AdmEsquina) Get(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := idDaEsquina(r.URL)
	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquina, err := esquinaDAO.BuscaCompletaPorId(id)
	if err != nil {
		esquinaDAO.Rollback()
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esquinaDAO.Commit()

	funcMap := template.FuncMap{
		"plural": func(tamanho int) bool {
			return tamanho != 1
		},
	}

	t, err := template.New("esquina").Funcs(funcMap).
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm-esquina.html")
	if err != nil {
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, "adm-esquina.html", esquina)
	if err != nil {
		log.Println(err)
	}
}
