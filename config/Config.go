package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	dbConnectionString := os.Getenv("root:104725@tcp(db:3306)/MENU")
	if dbConnectionString == "" {
		fmt.Println("DB_CONNECTION_STRING environment variable is not set")
		os.Exit(1)
	}
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		panic(err.Error())
	}
	return db
}
