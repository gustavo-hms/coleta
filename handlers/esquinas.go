package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/esquinas", Esquinas)
}

func Esquinas(w http.ResponseWriter, r *http.Request) {
}
