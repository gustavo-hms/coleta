package serviços

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func init() {
	registrarSeguro("/adm/zona/", new(AdmZona))
}

type AdmZona struct {
	id string
}

func (z *AdmZona) Get(w http.ResponseWriter, r *http.Request) {
	nós := strings.Split(r.URL.Path, "/")
	fim := len(nós)
	if fim < 4 {
		log.Println("URL inesperada:", r.URL.Path)
		erroInterno(w, r)
		return
	}

	if fim == 4 {
		z.id = nós[3]
		z.getZona(w, r)
		return
	}

	z.id = nós[fim-2]

	switch nós[fim-1] {
	case "esquinas":
		z.getEsquinas(w, r)
	case "lideres":
		//		z.getLíderes(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (z AdmZona) getZona(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zona, err := zonaDAO.BuscaCompleta(z.id)
	if err != nil {
		zonaDAO.Rollback()
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zonaDAO.Commit()

	t, err := template.New("zona").
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm-zona.html")
	if err != nil {
		log.Println("Ali:", err)
		return
	}

	err = t.ExecuteTemplate(w, "adm-zona.html", zona)
	if err != nil {
		log.Println(err)
	}
}

func (z AdmZona) getEsquinas(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esquinaDAO := dao.NewEsquinaDAO(&dao.Tx{tx})
	esquinas, err := esquinaDAO.BuscarPorZona(z.id)
	if err != nil {
		esquinaDAO.Rollback()
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esquinaDAO.Commit()

	json := z.jsonDasEsquinas(esquinas)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (z AdmZona) jsonDasEsquinas(esquinas []modelos.Esquina) string {
	objetos := make([]string, len(esquinas))
	for k, esquina := range esquinas {
		objetos[k] = `{"cruzamento": "` + esquina.Cruzamento +
			`", "id": ` + fmt.Sprint(esquina.Id) + `}`
	}

	return `[` + strings.Join(objetos, ",") + `]`
}
