package server

import (
	"encoding/json"
	"net/http"

	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/web"
)

// Server is an HTTP server adapter wiring API and web handlers.
// It depends only on the use case interfaces.
type Server struct {
	orderUseCase ports.OrderUseCase
	webHandler   *web.WebHandler
}

func NewServer(orderUseCase ports.OrderUseCase) *Server {
	return &Server{
		orderUseCase: orderUseCase,
		webHandler:   web.NewWebHandler(orderUseCase),
	}
}

func (s *Server) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Path[len("/order/"):]
	if orderUID == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
		return
	}

	order, err := s.orderUseCase.GetOrder(orderUID)
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
	stats, err := s.orderUseCase.GetStats()
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
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
