package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, name, date, venue, price, available, total FROM events ORDER BY date")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		rows.Scan(&e.ID, &e.Name, &e.Date, &e.Venue, &e.Price, &e.Available, &e.Total)
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) purchaseTickets(w http.ResponseWriter, r *http.Request) {
	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	tx, _ := s.db.Begin()
	defer tx.Rollback()

	for _, item := range req.Items {
		var available int
		err := tx.QueryRow("SELECT available FROM events WHERE id = ? FOR UPDATE", item.EventID).Scan(&available)
		if err != nil {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		if available < item.Quantity {
			http.Error(w, "Not enough tickets", http.StatusConflict)
			return
		}

		tx.Exec("UPDATE events SET available = available - ? WHERE id = ?", item.Quantity, item.EventID)
		tx.Exec(`INSERT INTO purchases (event_id, quantity, purchase_date, total_price)
                 SELECT id, ?, ?, price * ? FROM events WHERE id = ?`,
			item.Quantity, time.Now(), item.Quantity, item.EventID)
	}

	tx.Commit()
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
