package web

import (
	"html/template"
	"net/http"
	"wb-tech-l0/internal/repository/database"
)

type WebHandler struct {
	db *database.DB
}

func NewWebHandler(db *database.DB) *WebHandler {
	return &WebHandler{db: db}
}

func (h *WebHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func (h *WebHandler) OrderPageHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("uid")
	if orderUID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	order, err := h.db.GetOrder(orderUID)
	if err != nil {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, map[string]string{
			"Error": "Заказ не найден: " + err.Error(),
		})
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/order.html"))
	tmpl.Execute(w, order)
}
