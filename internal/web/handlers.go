package web

import (
	"html/template"
	"net/http"

	"wb-tech-l0/internal/application/ports"
)

// WebHandler is an HTTP adapter that talks only to the use case layer.
type WebHandler struct {
	orderUseCase ports.OrderUseCase
}

func NewWebHandler(orderUseCase ports.OrderUseCase) *WebHandler {
	return &WebHandler{orderUseCase: orderUseCase}
}

func (h *WebHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	_ = tmpl.Execute(w, nil)
}

func (h *WebHandler) OrderPageHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("uid")
	if orderUID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	order, err := h.orderUseCase.GetOrder(orderUID)
	if err != nil {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		_ = tmpl.Execute(w, map[string]string{
			"Error": "Заказ не найден: " + err.Error(),
		})
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/order.html"))
	_ = tmpl.Execute(w, order)
}
