package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Event struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Date      string  `json:"date"`
	Venue     string  `json:"venue"`
	Price     float64 `json:"price"`
	Available int     `json:"available"`
	Total     int     `json:"total"`
}

type CartItem struct {
	EventID  int `json:"event_id"`
	Quantity int `json:"quantity"`
}

type PurchaseRequest struct {
	Items []CartItem `json:"items"`
}

type Server struct {
	db    *sql.DB
	mutex sync.Mutex
}

func main() {
	// MySQL Connection String Format:
	// username:password@tcp(host:port)/database
	// UPDATE THESE VALUES FOR YOUR SYSTEM
	connStr := "ticketqueen_user:yourpassword@tcp(localhost:3306)/ticketqueen?parseTime=true"

	log.Println("Attempting to connect to MySQL database...")
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	log.Println("✓ MySQL database connection successful!")

	server := &Server{db: db}

	r := mux.NewRouter()
	r.HandleFunc("/api/events", server.getEvents).Methods("GET")
	r.HandleFunc("/api/purchase", server.purchaseTickets).Methods("POST")

	handler := corsMiddleware(r)

	log.Println("✓ Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/events - Fetching events...")

	rows, err := s.db.Query(`
		SELECT id, name, date, venue, price, available, total 
		FROM events 
		ORDER BY date
	`)
	if err != nil {
		log.Printf("ERROR querying events: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Name, &e.Date, &e.Venue, &e.Price, &e.Available, &e.Total)
		if err != nil {
			log.Printf("ERROR scanning event: %v", err)
			http.Error(w, fmt.Sprintf("Scan error: %v", err), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	log.Printf("✓ Successfully fetched %d events", len(events))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) purchaseTickets(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/purchase - Processing purchase...")

	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Purchase request for %d items", len(req.Items))

	// CRITICAL: Mutex lock to prevent race conditions
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("ERROR starting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, item := range req.Items {
		log.Printf("Processing event ID %d, quantity %d", item.EventID, item.Quantity)

		var available int
		// FOR UPDATE locks the row to prevent concurrent modifications
		err := tx.QueryRow("SELECT available FROM events WHERE id = ? FOR UPDATE", item.EventID).Scan(&available)
		if err != nil {
			log.Printf("ERROR fetching event: %v", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		if available < item.Quantity {
			log.Printf("ERROR: Not enough tickets. Available: %d, Requested: %d", available, item.Quantity)
			http.Error(w, fmt.Sprintf("Not enough tickets available. Only %d left", available), http.StatusConflict)
			return
		}

		// Update available tickets
		_, err = tx.Exec("UPDATE events SET available = available - ? WHERE id = ?", item.Quantity, item.EventID)
		if err != nil {
			log.Printf("ERROR updating availability: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Insert purchase record
		_, err = tx.Exec(`
			INSERT INTO purchases (event_id, quantity, purchase_date, total_price)
			SELECT id, ?, ?, price * ?
			FROM events WHERE id = ?
		`, item.Quantity, time.Now(), item.Quantity, item.EventID)
		if err != nil {
			log.Printf("ERROR inserting purchase: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("ERROR committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("✓ Purchase completed successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Purchase completed successfully"})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
