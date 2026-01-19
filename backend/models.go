package main

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
