package serviços

import (
	"coleta/dao"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func init() {
	registrar("/busca/lideres", BuscaLíderes{})
	registrarSeguro("/adm/busca/lideres", BuscaLíderes{})
}

type BuscaLíderes struct{}

func (b BuscaLíderes) Get(w http.ResponseWriter, r *http.Request) {
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

	líderDAO := dao.NewLiderDAO(tx)
	líderes, err := líderDAO.FindAllThatMatches(trecho)
	if err != nil {
		líderDAO.Rollback()
		log.Println(err)
		erroInterno(w, r)
		return
	}
	if err := líderDAO.Commit(); err != nil {
		líderDAO.Rollback()
		log.Println("Erro no commit:", err)
		erroInterno(w, r)
		return
	}

	bytes, err := json.Marshal(líderes)
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}
