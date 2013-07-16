package main

import (
	"coleta/config"
	"coleta/dao"
	"coleta/serviços"
	"fmt"
	"log"
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

	f, _ := os.Open("coleta.log")
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go http.ListenAndServe(":"+config.Dados.Porta, serviços.MuxSimples)

	err := http.ListenAndServeTLS(
		":"+config.Dados.PortaTLS,
		config.Dados.ArquivoDeCertificado,
		config.Dados.ArquivoDeChave,
		serviços.MuxSeguro,
	)

	if err != nil {
		fmt.Println("Erro ao subir servidor https:", err)
	}
}

func uso() string {
	return "Uso: " + os.Args[0] + " <arquivo de configuração>"
}
