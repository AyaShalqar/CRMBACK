package shop

import (
	"encoding/json"
	"net/http"

	"crm-backend/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateShopHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	if claims == nil || claims.Role != "superadmin" {
		http.Error(w, "доступ запрещён", http.StatusForbidden)
		return
	}
	var shop Shop
	if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
		http.Error(w, "неправильный формат данных", http.StatusBadRequest)
		return
	}
	if shop.OwnerID == 0 {
		http.Error(w, "нужно указать владельца магазина (owner_id)", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateShop(r.Context(), shop); err != nil {
		http.Error(w, "ошибка создания магазина", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Магазин создан для пользователя"))
}

func (h *Handler) GetShopsHandler(w http.ResponseWriter, r *http.Request) {
	shops, err := h.service.GetShops(r.Context())
	if err != nil {
		http.Error(w, "ошибка получения магазинов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shops)
}

func (h *Handler) GetShopsByOwner(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}
	shops, err := h.service.GetShopsByOwner(r.Context(), claims.ID)
	if err != nil {
		http.Error(w, "ошибка получения магазинов", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shops)
}
