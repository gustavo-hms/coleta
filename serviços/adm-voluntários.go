package serviços

import (
	"coleta/config"
	"coleta/dao"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguro("/adm/voluntarios", new(AdmVoluntários))
}

type AdmVoluntários struct{}

func (e *AdmVoluntários) Get(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntários, err := voluntárioDAO.Todos()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t, err := template.New("voluntários").ParseFiles(
		config.Dados.DiretórioDasPáginas+"/adm-voluntários.html",
		config.Dados.DiretórioDasPáginas+"/voluntário.html",
	)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, "adm-voluntários.html", voluntários)
	if err != nil {
		log.Println(err)
	}
}
