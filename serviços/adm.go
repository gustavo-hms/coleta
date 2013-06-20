package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/adm", Adm)
}

func Adm(w http.ResponseWriter, r *http.Request) {
	página, err := ioutil.ReadFile(gopath + "/src/coleta/páginas/adm.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo adm.html:", err)
	}

	fmt.Fprintf(w, "%s", página)

}
