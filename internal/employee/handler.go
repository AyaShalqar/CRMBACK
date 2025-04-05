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

// AddEmployee godoc
// @Summary Add employee to shop
// @Description Add a new employee to a shop
// @Tags employees
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "Shop ID"
// @Param employee body Employee true "Employee data"
// @Success 201 {string} string "Сотрудник добавлен"
// @Failure 400 {string} string "неправильный ID магазина"
// @Failure 400 {string} string "неправильный формат данных"
// @Failure 401 {string} string "не авторизован"
// @Failure 403 {string} string "error message"
// @Router /owner/shops/{id}/employees [post]
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

// GetEmployeesByShop godoc
// @Summary Get shop employees
// @Description Get all employees of a specific shop
// @Tags employees
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "Shop ID"
// @Success 200 {array} Employee
// @Failure 400 {string} string "неправильный ID магазина"
// @Failure 401 {string} string "не авторизован"
// @Failure 403 {string} string "error message"
// @Router /owner/shops/{id}/employees [get]
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

// RemoveEmployee godoc
// @Summary Remove employee
// @Description Remove an employee from a shop
// @Tags employees
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "Shop ID"
// @Param employee_id path int true "Employee ID"
// @Success 200 {string} string "Сотрудник удалён"
// @Failure 400 {string} string "неправильный ID сотрудника"
// @Failure 401 {string} string "не авторизован"
// @Failure 403 {string} string "error message"
// @Router /owner/shops/{id}/employees/{employee_id} [delete]
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
