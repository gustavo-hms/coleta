package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/adm", Adm)
}

func Adm(w http.ResponseWriter, r *http.Request) {
}
