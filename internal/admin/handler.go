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

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "не удалось получить пользователей: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

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
	token, err := auth.GenerateJWT(user.Email, user.Role)
	if err != nil {
		http.Error(w, "ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
