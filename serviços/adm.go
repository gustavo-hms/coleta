package serviços

import (
	"coleta/config"
	"coleta/dao"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguro("/adm", Adm{})
	registrar("/adm/", Redirecionamento{})
}

type Redirecionamento struct{}

func (_ Redirecionamento) Get(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://tetocoleta.com.br"+r.URL.String(), http.StatusSeeOther)
}

type Adm struct{}

func (a Adm) Get(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println("Início da transação:", err)
		return
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		zonaDAO.Rollback()
		log.Println(err)
		return
	}

	zonaDAO.Commit()

	t, err := template.New("adm").
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm.html")
	if err != nil {
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, "adm.html", zonas)
	if err != nil {
		log.Println(err)
	}
}
