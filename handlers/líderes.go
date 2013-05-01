package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/líderes", Líderes)
}

func Líderes(w http.ResponseWriter, r *http.Request) {
}
