package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/adm/esquinas", AdmEsquinas)
}

func AdmEsquinas(w http.ResponseWriter, r *http.Request) {
}
