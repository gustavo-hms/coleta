package serviços

import (
	"coleta/dao"
	"coleta/modelos"
	"fmt"
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
	if fim < 3 {
		log.Println("URL inesperada:", r.URL.Path)
		erroInterno(w, r)
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

func (z AdmZona) getEsquinas(w http.ResponseWriter, r *http.Request) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
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

func (z AdmZona) jsonDasEsquinas(esquinas []*modelos.Esquina) string {
	objetos := make([]string, len(esquinas))
	for k, esquina := range esquinas {
		objetos[k] = `{"cruzamento": "` + esquina.Cruzamento +
			`", "id": ` + fmt.Sprint(esquina.Id) + `}`
	}

	return `[` + strings.Join(objetos, ",") + `]`
}
