package server

import (
	"context"
	"encoding/json"
	"net/http"

	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/web"
)

// Server is an HTTP server adapter wiring API and web handlers.
// It depends only on the use case interfaces and wraps http.Server
// to allow graceful shutdown.
type Server struct {
	orderUseCase ports.OrderUseCase
	webHandler   *web.WebHandler
	httpServer   *http.Server
}

func NewServer(orderUseCase ports.OrderUseCase) *Server {
	webHandler := web.NewWebHandler(orderUseCase)

	mux := http.NewServeMux()

	s := &Server{
		orderUseCase: orderUseCase,
		webHandler:   webHandler,
		httpServer: &http.Server{
			Handler: mux,
		},
	}

	// API routes
	mux.HandleFunc("/order/", s.GetOrderHandler)
	mux.HandleFunc("/stats", s.StatsHandler)

	// Web routes
	mux.HandleFunc("/", s.webHandler.IndexHandler)
	mux.HandleFunc("/order", s.webHandler.OrderPageHandler)

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return s
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
	stats, err := s.orderUseCase.Stats()
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}

func (s *Server) Start(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
