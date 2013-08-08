package serviços

import (
	"coleta/config"
	"coleta/dao"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguroComTransação("/adm/voluntarios", new(AdmVoluntários))
}

type AdmVoluntários struct{}

func (e *AdmVoluntários) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntários, err := voluntárioDAO.Todos()
	if err != nil {
		log.Println(err)
		return err
	}

	t, err := template.New("voluntários").ParseFiles(
		config.Dados.DiretórioDasPáginas+"/adm-voluntários.html",
		config.Dados.DiretórioDasPáginas+"/voluntário.html",
	)
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm-voluntários.html", voluntários)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
