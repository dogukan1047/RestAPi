package main

import (
	controller "awesomeProject/controller"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/happymoons", controller.Handler).Methods("GET")
	router.HandleFunc("/happymoons/csv", controller.CsvHandler).Methods("GET")
	happymoons := router.PathPrefix("/happymoons").Subrouter()
	happymoons1 := router.PathPrefix("/happymoons").Subrouter()
	happymoons1.HandleFunc("", controller.InHandler).Queries("in", "{in}")
	happymoons.HandleFunc("", controller.ExHandler).Queries("ex", "{ex}")
	log.Fatal(http.ListenAndServe(":9000", router))

}
