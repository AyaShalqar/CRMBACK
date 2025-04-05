package admin

import (
	"crm-backend/internal/auth"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *Service
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with admin rights
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param user body CreateUserDto true "User data"
// @Success 201 {string} string "Пользователь создан"
// @Failure 400 {string} string "неправильный формат данных"
// @Failure 500 {string} string "не удалось создать пользователя"
// @Router /admin/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto CreateUserDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "неправильный формат данных", http.StatusBadRequest)
		return
	}
	err := h.service.CreateUserByAdmin(r.Context(), dto)
	if err != nil {
		http.Error(w, "не удалось создать пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь создан"))
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {array} User
// @Failure 500 {string} string "не удалось получить пользователей"
// @Router /admin/users [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "не удалось получить пользователей: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "User ID"
// @Success 200 {string} string "пользователь удален"
// @Failure 400 {string} string "неправильный ID"
// @Failure 500 {string} string "не удалось удалить пользователя"
// @Router /admin/users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неправильный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, "не удалось удалить пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("пользователь удален"))
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's information
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "User ID"
// @Param user body CreateUserDto true "User data"
// @Success 200 {string} string "Пользователь обновлён"
// @Failure 400 {string} string "неправильный ID"
// @Failure 400 {string} string "неправильный формат данных"
// @Failure 500 {string} string "не удалось обновить пользователя"
// @Router /admin/users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неправильный ID", http.StatusBadRequest)
		return
	}
	var dto CreateUserDto

	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "неправильный формат данных", http.StatusBadRequest)
		return
	}

	user := User{
		ID:        id,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Role:      dto.Role,
	}
	err = h.service.UpdateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "не удалось обновить пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь обновлён"))
}

// Login godoc
// @Summary Authenticate user
// @Description Log in with username and password to get access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {string} string "неправильный формат данных"
// @Failure 401 {string} string "пользователь не найден"
// @Failure 401 {string} string "неверный пароль"
// @Failure 500 {string} string "ошибка генерации токена"
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неправильный формат данных", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "пользователь не найден", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	// ✅ Теперь передаём user.ID в токен
	token, err := auth.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
