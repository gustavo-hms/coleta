package serviços

import (
	"coleta/config"
	"coleta/dao"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguroComTransação("/adm/lideres", new(AdmLíderes))
}

type AdmLíderes struct{}

func (e *AdmLíderes) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	líderDAO := dao.NewLiderDAO(tx)
	líderes, err := líderDAO.Todos()
	if err != nil {
		log.Println(err)
		return err
	}

	t, err := template.New("líderes").ParseFiles(
		config.Dados.DiretórioDasPáginas+"/adm-líderes.html",
		config.Dados.DiretórioDasPáginas+"/líder.html",
	)
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm-líderes.html", líderes)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
