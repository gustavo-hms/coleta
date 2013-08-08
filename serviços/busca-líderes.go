package serviços

import (
	"coleta/dao"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func init() {
	registrarComTransação("/busca/lideres", BuscaLíderes{})
	registrarSeguroComTransação("/adm/busca/lideres", BuscaLíderes{})
}

type BuscaLíderes struct{}

func (b BuscaLíderes) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	trecho := r.FormValue("contem")
	if len(trecho) < 3 {
		w.WriteHeader(http.StatusNotFound)
	}

	líderDAO := dao.NewLiderDAO(tx)
	líderes, err := líderDAO.FindAllThatMatches(trecho)
	if err != nil {
		log.Println(err)
		return err
	}

	bytes, err := json.Marshal(líderes)
	if err != nil {
		log.Println(err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))

	return nil
}
