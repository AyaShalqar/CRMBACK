package admin

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
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
