package main

import (
	"coleta/config"
	"coleta/db"
	_ "coleta/handlers"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(uso())
	}

	if err := config.Ler(os.Args[1]); err != nil {
		fmt.Printf("Não foi possível ler o arquivo de configuração %s. Erro: %s", os.Args[1], err)
		os.Exit(1)
	}

	if err := db.Conn(); err != nil {
		fmt.Println("Não foi possível conectar-se ao banco. Erro:", err)
		os.Exit(1)
	}

	http.ListenAndServe(":8080", nil)
}

func uso() string {
	return "Uso: " + os.Args[0] + " <arquivo de configuração>"
}
