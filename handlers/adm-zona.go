package handlers

import (
	"coleta/dao"
	"coleta/modelos"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	http.HandleFunc("/adm/zona/", Zona)
}

func Zona(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquinas, err := esquinaDAO.BuscarPorZona(idDaZona(r.URL))
	if err != nil {
		esquinaDAO.Rollback()
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esquinaDAO.Commit()

	json := jsonDasEsquinas(esquinas)
	fmt.Fprint(w, json)
}

func idDaZona(endereço *url.URL) string {
	return strings.SplitN(endereço.Path[10:], "/", 2)[0]
}

func jsonDasEsquinas(esquinas []*modelos.Esquina) string {
	objetos := make([]string, len(esquinas))
	for k, esquina := range esquinas {
		objetos[k] = `{"cruzamento": "` + esquina.Cruzamento +
			`", "link": "/adm/esquina/` + fmt.Sprint(esquina.Id) + `"}`
	}

	return `[` + strings.Join(objetos, ",") + `]`
}
