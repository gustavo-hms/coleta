package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/líder/", Líder)
}

func Líder(w http.ResponseWriter, r *http.Request) {
}
