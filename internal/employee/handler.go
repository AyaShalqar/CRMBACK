package employee

import (
	"crm-backend/internal/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

// NewHandler — конструктор для хендлера
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// AddEmployee godoc
// @Summary Add new employee
// @Description Создаёт в таблице users пользователя с role='employee' и привязывает к магазину.
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Shop ID"
// @Param Authorization header string true "Bearer <token>"
// @Param data body AddEmployeeRequest true "Данные сотрудника"
// @Success 201 {string} string "Сотрудник успешно добавлен"
// @Failure 400 {string} string "неправильный формат данных/ID магазина"
// @Failure 401 {string} string "не авторизован"
// @Failure 403 {string} string "доступ запрещён"
// @Failure 500 {string} string "ошибка создания пользователя/ошибка привязки"
// @Router /owner/shops/{id}/employees [post]
func (h *Handler) AddEmployee(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	// Разрешаем только владельцу или супер-админу
	if claims.Role != "owner" && claims.Role != "superadmin" {
		http.Error(w, "доступ запрещён", http.StatusForbidden)
		return
	}

	shopIDStr := chi.URLParam(r, "id")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		http.Error(w, "неправильный ID магазина", http.StatusBadRequest)
		return
	}

	var req AddEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неправильный формат данных", http.StatusBadRequest)
		return
	}

	// 1) Создаём запись в таблице users (role='employee')
	userID, err := h.service.CreateUserForEmployee(r.Context(), req)
	if err != nil {
		http.Error(w, "ошибка создания пользователя-сотрудника: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) Привязываем user_id к shop_id (employees)
	err = h.service.AddEmployeeLink(r.Context(), userID, shopID, req.Position, claims.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка привязки сотрудника: %v", err), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Сотрудник успешно добавлен"))
}

// GetEmployeesByShop godoc
// @Summary Get employees of a shop
// @Description Возвращает список сотрудников конкретного магазина (требует роль owner или superadmin).
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Shop ID"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {array} Employee
// @Failure 400 {string} string "неправильный ID магазина"
// @Failure 401 {string} string "не авторизован"
// @Failure 403 {string} string "доступ запрещён"
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
	_ = json.NewEncoder(w).Encode(employees)
}

// RemoveEmployee godoc
// @Summary Remove employee
// @Description Удаляет сотрудника из магазина (требует роль owner или superadmin).
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Shop ID"
// @Param employee_id path int true "Employee ID"
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {string} string "Сотрудник удалён"
// @Failure 400 {string} string "неправильный ID сотрудника"
// @Failure 401 {string} string "не авторизован"
// @Failure 403 {string} string "доступ запрещён"
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
