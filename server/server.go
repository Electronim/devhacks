package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer(port string) {
	router := mux.NewRouter()

	router.HandleFunc("/api/test", TestHandler).Methods("POST")
	router.HandleFunc("/api/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/product/add", ProductAdd).Methods("POST")
	router.HandleFunc("/api/product/get", ProductGet).Methods("POST")
	router.HandleFunc("/api/receipt/create", ReceiptCreate).Methods("POST")
	router.HandleFunc("/api/receipt/confirm", ReceiptConfirm).Methods("POST")
	router.HandleFunc("/api/receipt/get", ReceiptGet).Methods("POST")

	fmt.Println("Server starting on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
