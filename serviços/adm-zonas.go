package serviços

import (
	"coleta/config"
	"coleta/dao"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func init() {
	registrarSeguro("/adm/zonas", AdmZonas{})
}

type AdmZonas struct{}

func (a AdmZonas) Get(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
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

	t, err := template.New("adm-zonas").
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm-zonas.html")
	if err != nil {
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, "adm-zonas.html", zonas)
	if err != nil {
		log.Println(err)
	}
}

func (a AdmZonas) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		erroInterno(w, r)
		return
	}

	if len(r.Form["bloqueadas"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Erro ao ler formulário")
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		zonaDAO.Rollback()
		log.Println(err)
		erroInterno(w, r)
		return
	}

	for _, zona := range zonas {
		if encontrada(zona.Id, r.Form["bloqueadas"]) {
			zona.Bloqueada = true
		} else {
			zona.Bloqueada = false
		}

		err = zonaDAO.Save(zona)
		if err != nil && err != dao.ErrRowsNotAffected {
			zonaDAO.Rollback()
			log.Println(err)
			erroInterno(w, r)
			return
		}
	}

	err = zonaDAO.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		erroInterno(w, r)
		return
	}

	http.Redirect(w, r, ".", http.StatusSeeOther)
}

func encontrada(id int, bloqueadas []string) bool {
	stringDoId := strconv.Itoa(id)
	for _, bloqueada := range bloqueadas {
		if bloqueada == stringDoId {
			return true
		}
	}

	return false
}
