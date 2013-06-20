package main

import (
	"coleta/config"
	"coleta/dao"
	_ "coleta/serviços"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) != 2 {
		fmt.Println(uso())
		os.Exit(1)
	}

	if err := config.Ler(os.Args[1]); err != nil {
		fmt.Printf("Não foi possível ler o arquivo de configuração %s. Erro: %s", os.Args[1], err)
		os.Exit(1)
	}

	if err := dao.Conn(); err != nil {
		fmt.Println("Não foi possível conectar-se ao banco. Erro:", err)
		os.Exit(1)
	}

	http.ListenAndServe(":"+config.Dados.Porta, nil)
}

func uso() string {
	return "Uso: " + os.Args[0] + " <arquivo de configuração>"
}
