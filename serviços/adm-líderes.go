package serviços

import (
	"coleta/config"
	"coleta/dao"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguro("/adm/lideres", new(AdmLíderes))
}

type AdmLíderes struct{}

func (e *AdmLíderes) Get(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	líderes, err := líderDAO.Todos()
	if err != nil {
		líderDAO.Rollback()
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := líderDAO.Commit(); err != nil {
		líderDAO.Rollback()
		log.Println(err)
		erroInterno(w, r)
		return
	}

	t, err := template.New("líderes").ParseFiles(
		config.Dados.DiretórioDasPáginas+"/adm-líderes.html",
		config.Dados.DiretórioDasPáginas+"/líder.html",
	)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, "adm-líderes.html", líderes)
	if err != nil {
		log.Println(err)
	}
}
