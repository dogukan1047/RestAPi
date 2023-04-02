package Config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {

	db, err := sql.Open("mysql", "root:104725@tcp(127.0.0.1:3306)/MENU")

	if err != nil {
		panic(err.Error())
	}

	return db

}
