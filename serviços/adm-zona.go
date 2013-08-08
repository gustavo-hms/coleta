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
	registrarSeguroComTransação("/adm/zona/", new(AdmZona))
}

type AdmZona struct {
	id string
}

func (z *AdmZona) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	nós := strings.Split(r.URL.Path, "/")
	fim := len(nós)
	if fim < 4 {
		log.Println("URL inesperada:", r.URL.Path)
		return erroInesperado
	}

	if fim == 4 {
		z.id = nós[3]
		return z.getZona(w, r, tx)
	}

	z.id = nós[fim-2]

	switch nós[fim-1] {
	case "esquinas":
		return z.getEsquinas(w, r, tx)

	default:
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
}

func (z AdmZona) getZona(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	zonaDAO := dao.NewZonaDAO(tx)
	zona, err := zonaDAO.BuscaCompleta(z.id)
	if err != nil {
		log.Println(err)
		return err
	}

	t, err := template.New("zona").
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm-zona.html")
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm-zona.html", zona)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (z AdmZona) getEsquinas(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquinas, err := esquinaDAO.BuscarPorZona(z.id)
	if err != nil {
		log.Println(err)
		return err
	}

	json := z.jsonDasEsquinas(esquinas)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, json)

	return nil
}

func (z AdmZona) jsonDasEsquinas(esquinas []modelos.Esquina) string {
	objetos := make([]string, len(esquinas))
	for k, esquina := range esquinas {
		objetos[k] = `{"cruzamento": "` + esquina.Cruzamento +
			`", "id": ` + fmt.Sprint(esquina.Id) + `}`
	}

	return `[` + strings.Join(objetos, ",") + `]`
}
