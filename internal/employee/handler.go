package employee

import (
	"encoding/json"
	"net/http"
	"strconv"

	"crm-backend/internal/auth"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddEmployee(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	shopIDStr := chi.URLParam(r, "id")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		http.Error(w, "неправильный ID магазина", http.StatusBadRequest)
		return
	}

	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "неправильный формат данных", http.StatusBadRequest)
		return
	}
	employee.ShopID = shopID

	err = h.service.AddEmployee(r.Context(), claims.ID, employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Сотрудник добавлен"))
}

func (h *Handler) GetEmployeesByShop(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	shopIDStr := chi.URLParam(r, "id")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		http.Error(w, "неправильный ID магазина", http.StatusBadRequest)
		return
	}

	employees, err := h.service.GetEmployeesByShop(r.Context(), claims.ID, shopID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func (h *Handler) RemoveEmployee(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	employeeIDStr := chi.URLParam(r, "employee_id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		http.Error(w, "неправильный ID сотрудника", http.StatusBadRequest)
		return
	}

	err = h.service.RemoveEmployee(r.Context(), claims.ID, employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сотрудник удалён"))
}
func (h *Handler) CreateItems(w http.ResponseWriter, r *http.Request) {

}
