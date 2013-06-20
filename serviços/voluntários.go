package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/voluntários", Voluntários)
}

func Voluntários(w http.ResponseWriter, r *http.Request) {
}
