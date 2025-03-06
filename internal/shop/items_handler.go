package shop

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemHandler struct {
	service *ItemService
}

func NewItemHandler(service *ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

// CreateItemHandler - POST /owner/shops/{shopID}/items
func (h *ItemHandler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	shopIDStr := chi.URLParam(r, "shopID")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		http.Error(w, "invalid shopID", http.StatusBadRequest)
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Важно: привязываем item.ShopID к значению из URL,
	// чтобы пользователь не мог «подменить» shop_id в теле JSON.
	item.ShopID = shopID

	if err := h.service.CreateItem(r.Context(), &item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданный item (с уже присвоенным ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

// GetItemsHandler - GET /owner/shops/{shopID}/items
func (h *ItemHandler) GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	shopIDStr := chi.URLParam(r, "shopID")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		http.Error(w, "invalid shopID", http.StatusBadRequest)
		return
	}

	items, err := h.service.GetItemsForShop(r.Context(), shopID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(items)
}

// GetItemHandler - GET /owner/shops/{shopID}/items/{itemID}
func (h *ItemHandler) GetItemHandler(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemID")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "invalid itemID", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetItemByID(r.Context(), itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}

// UpdateItemHandler - PUT /owner/shops/{shopID}/items/{itemID}
func (h *ItemHandler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	shopIDStr := chi.URLParam(r, "shopID")
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		http.Error(w, "invalid shopID", http.StatusBadRequest)
		return
	}

	itemIDStr := chi.URLParam(r, "itemID")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "invalid itemID", http.StatusBadRequest)
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Привязываем ID из URL, чтобы нельзя было подменить
	item.ID = itemID
	item.ShopID = shopID

	if err := h.service.UpdateItem(r.Context(), item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DeleteItemHandler - DELETE /owner/shops/{shopID}/items/{itemID}
func (h *ItemHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemID")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "invalid itemID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteItem(r.Context(), itemID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
