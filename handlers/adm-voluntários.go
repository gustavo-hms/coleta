package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/adm/voluntários", AdmVoluntários)
}

func AdmVoluntários(w http.ResponseWriter, r *http.Request) {
}
