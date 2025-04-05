package auth

import (
	"encoding/json"
	"net/http"
)

// Handler — хендлер auth
type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// LoginRequest — входящие данные из JSON
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse — ответ при успехе
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user,omitempty"`
}

// Login — POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, token, err := h.service.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Успех
	resp := LoginResponse{
		Token: token,
		User:  user, // Если хотим сразу вернуть user
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Me godoc
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} map[string]interface{} "User data with id, email, and role"
// @Failure 401 {string} string "No user in context"
// @Router /auth/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "No user in context", http.StatusUnauthorized)
		return
	}

	// Можно вернуть claims напрямую, либо сформировать объект
	userData := map[string]interface{}{
		"id":    claims.ID,
		"email": claims.Email,
		"role":  claims.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}
