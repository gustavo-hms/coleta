package serviços

import (
	"coleta/config"
	"coleta/dao"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguroComTransação("/adm", Adm{})
	registrar("/adm", Redirecionamento{})
	registrar("/adm/", Redirecionamento{})
}

type Redirecionamento struct{}

func (_ Redirecionamento) Get(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://tetocoleta.com.br"+r.URL.String(), http.StatusSeeOther)
}

type Adm struct{}

func (a Adm) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		log.Println(err)
		return err
	}

	t, err := template.New("adm").
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm.html")
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm.html", zonas)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
