package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/esquina/", Esquina)
}

func Esquina(w http.ResponseWriter, r *http.Request) {
}
