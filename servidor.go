package main

import (
	_ "coleta/handlers"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", nil)
}
