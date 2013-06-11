package db

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
	senha := config.Dados.Banco.Senha
	host := config.Dados.Banco.Host
	base := config.Dados.Banco.Base

	parâmetros := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=require host=%s",
		usuário, senha, base, host)

	DB, err = sql.Open("mysql", parâmetros)
	if err != nil {
		log.Println(err)
		return
	}

	DB.SetMaxIdleConns(16)

	return
}
