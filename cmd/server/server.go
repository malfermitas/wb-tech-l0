package server

import (
	"encoding/json"
	"net/http"
	"wb-tech-l0/internal/repository/database"
	"wb-tech-l0/internal/web"
)

type Server struct {
	db         *database.DB
	webHandler *web.WebHandler
}

func NewServer(db *database.DB) *Server {
	return &Server{
		db:         db,
		webHandler: web.NewWebHandler(db),
	}
}

func (s *Server) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Path[len("/order/"):]
	if orderUID == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
		return
	}

	order, err := s.db.GetOrder(orderUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) StatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"cache_size": s.db.Cache.Size(),
	}

	count, err := s.db.GetOrderCount()
	if err == nil {
		stats["db_count"] = count
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) Start(addr string) error {
	// API routes
	http.HandleFunc("/order/", s.GetOrderHandler)
	http.HandleFunc("/stats", s.StatsHandler)

	// Web routes
	http.HandleFunc("/", s.webHandler.IndexHandler)
	http.HandleFunc("/order", s.webHandler.OrderPageHandler)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return http.ListenAndServe(addr, nil)
}
