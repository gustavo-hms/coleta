package handlers

import (
	"net/http"
)

func init() {
	http.HandleFunc("/voluntário/", Voluntário)
}

func Voluntário(w http.ResponseWriter, r *http.Request) {
}
