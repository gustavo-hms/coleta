package dao

import (
	"coleta/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func Conn() (err error) {
	usuário := config.Dados.Banco.Usuário
	host := config.Dados.Banco.Host
	base := config.Dados.Banco.Base
	parâmetros := fmt.Sprintf("%s@tcp(%s:3306)/%s", usuário, host, base)

	DB, err = sql.Open("mysql", parâmetros)
	if err != nil {
		log.Println("Erro ao conectar-se ao banco:", err)
		return
	}

	DB.SetMaxIdleConns(16)

	return
}
