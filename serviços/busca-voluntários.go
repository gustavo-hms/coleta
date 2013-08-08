package serviços

import (
	"coleta/dao"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func init() {
	registrarComTransação("/busca/voluntarios", BuscaVoluntários{})
	registrarSeguroComTransação("/adm/busca/voluntarios", BuscaVoluntários{})
}

type BuscaVoluntários struct{}

func (b BuscaVoluntários) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	trecho := r.FormValue("contem")
	if len(trecho) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntários, err := voluntárioDAO.FindAllThatMatches(trecho)
	if err != nil {
		log.Println(err)
		return err
	}

	bytes, err := json.Marshal(voluntários)
	if err != nil {
		log.Println(err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))

	return nil
}
