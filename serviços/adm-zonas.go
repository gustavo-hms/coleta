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
	registrarSeguroComTransação("/adm/zonas", AdmZonas{})
}

type AdmZonas struct{}

func (a AdmZonas) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		log.Println(err)
		return err
	}

	t, err := template.New("adm-zonas").
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm-zonas.html")
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm-zonas.html", zonas)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a AdmZonas) Post(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		return err
	}

	if len(r.Form["bloqueadas"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Erro ao ler formulário")
		return nil
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, zona := range zonas {
		if encontrada(zona.Id, r.Form["bloqueadas"]) {
			zona.Bloqueada = true
		} else {
			zona.Bloqueada = false
		}

		err = zonaDAO.Save(zona)
		if err != nil && err != dao.ErrRowsNotAffected {
			log.Println(err)
			return err
		}
	}

	http.Redirect(w, r, ".", http.StatusSeeOther)

	return nil
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
