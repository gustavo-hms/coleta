package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Conn() (*sql.DB, error) {
	return sql.Open("mysql", "root@tcp(localhost:3306)/coleta")
}
