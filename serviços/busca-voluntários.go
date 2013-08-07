package serviços

import (
	"coleta/dao"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func init() {
	registrar("/busca/voluntarios", BuscaVoluntários{})
	registrarSeguro("/adm/busca/voluntarios", BuscaVoluntários{})
}

type BuscaVoluntários struct{}

func (b BuscaVoluntários) Get(w http.ResponseWriter, r *http.Request) {
	trecho := r.FormValue("contem")
	if len(trecho) < 3 {
		w.WriteHeader(http.StatusNotFound)
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println("Início da transação:", err)
		erroInterno(w, r)
		return
	}

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntários, err := voluntárioDAO.FindAllThatMatches(trecho)
	if err != nil {
		voluntárioDAO.Rollback()
		log.Println(err)
		erroInterno(w, r)
		return
	}
	if err := voluntárioDAO.Commit(); err != nil {
		voluntárioDAO.Rollback()
		log.Println(err)
		erroInterno(w, r)
		return
	}

	bytes, err := json.Marshal(voluntários)
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}
