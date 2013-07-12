package serviços

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
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
	turno   string
	esquina *modelos.Esquina
}

func idDaEsquina(endereço *url.URL) string {
	idx := strings.LastIndex(endereço.Path, "/")
	return endereço.Path[idx+1:]
}

func (e *AdmEsquina) Get(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := idDaEsquina(r.URL)
	esquinaDAO := dao.NewEsquinaDAO(tx)
	e.esquina, err = esquinaDAO.BuscaCompletaPorId(id)
	if err != nil {
		esquinaDAO.Rollback()
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esquinaDAO.Commit()

	turno := r.FormValue("turno")
	if turno != "" {
		e.turno = turno
		e.exibirTurno(w, r)
	} else {
		e.exibiçãoGeral(w, r)
	}
}

func (e AdmEsquina) exibirTurno(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"turno": func() string {
			return e.turno
		},
	}

	t, err := template.New("esquina").Funcs(funcMap).
		ParseFiles(
		config.Dados.DiretórioDasPáginas+"/adm-esquina-turno.html",
		config.Dados.DiretórioDasPáginas+"/adm-esquina-líder.html",
		config.Dados.DiretórioDasPáginas+"/adm-esquina-voluntário.html",
	)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, "adm-esquina-turno.html", e.esquina)
	if err != nil {
		log.Println(err)
	}
}

func (e AdmEsquina) exibiçãoGeral(w http.ResponseWriter, r *http.Request) {
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

	err = t.ExecuteTemplate(w, "adm-esquina.html", e.esquina)
	if err != nil {
		log.Println(err)
	}
}
