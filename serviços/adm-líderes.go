package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/adm/líderes", AdmLíderes)
}

func AdmLíderes(w http.ResponseWriter, r *http.Request) {
}
