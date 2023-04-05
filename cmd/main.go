package main

import (
	controller "RestApiGo/controller"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net"
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

	// Rastgele bir port seçiliyor
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	// Dinamik olarak atanan port bilgisi alınıyor
	addr := listener.Addr().String()

	// Servis Başlatılıyor
	server := &http.Server{
		Handler: router,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Seçilen dinamik portun log çıktısı alınır
	log.Printf("Service is running on port: %v", addr)

	// Servis sonlandığında kapatılır
	defer func() {
		if err := server.Shutdown(nil); err != nil {
			log.Fatal(err)
		}
		log.Println("Service stopped")
	}()

	// Program kapatılana kadar servisin çalışmasını sağlamak için bir döngü oluşturulur
	select {}
}

/*
• Servis config dosyası ile yönetilebilir olmalı (+++---)
• Solid prensiplerine uygun olmalı(++++++++++)
• Servisin hangi port üzerinde çalıştığı client tarafından bilinmeyecek(+++++++++)
• Api gateway ile servis güvenli hale getirilecek(------)
• Servis endpointleri aşağıdaki gibi olmalıdır(+++++++++++)
*/
