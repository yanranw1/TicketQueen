package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Server struct {
	db    *sql.DB
	mutex sync.Mutex
}

func main() {
	connStr := "ticketqueen_user:yourpassword@tcp(localhost:3306)/ticketqueen?parseTime=true"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}

	server := &Server{db: db}
	r := mux.NewRouter()

	r.HandleFunc("/api/events", server.getEvents).Methods("GET")
	r.HandleFunc("/api/purchase", server.purchaseTickets).Methods("POST")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}
