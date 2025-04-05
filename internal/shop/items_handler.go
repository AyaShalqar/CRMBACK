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

// CreateItemHandler godoc
// @Summary Create new item
// @Description Create a new item in a shop
// @Tags items
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param shopID path int true "Shop ID"
// @Param item body Item true "Item details"
// @Success 201 {object} Item
// @Failure 400 {string} string "invalid shopID"
// @Failure 400 {string} string "invalid request body"
// @Failure 500 {string} string "error message"
// @Router /owner/shops/{shopID}/items [post]
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

// GetItemsHandler godoc
// @Summary Get all items in a shop
// @Description Get all items for a specific shop
// @Tags items
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param shopID path int true "Shop ID"
// @Success 200 {array} Item
// @Failure 400 {string} string "invalid shopID"
// @Failure 500 {string} string "error message"
// @Router /owner/shops/{shopID}/items [get]
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

// GetItemHandler godoc
// @Summary Get specific item
// @Description Get a specific item by ID
// @Tags items
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param shopID path int true "Shop ID"
// @Param itemID path int true "Item ID"
// @Success 200 {object} Item
// @Failure 400 {string} string "invalid itemID"
// @Failure 404 {string} string "error message"
// @Router /owner/shops/{shopID}/items/{itemID} [get]
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

// UpdateItemHandler godoc
// @Summary Update item
// @Description Update an existing item
// @Tags items
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param shopID path int true "Shop ID"
// @Param itemID path int true "Item ID"
// @Param item body Item true "Item details"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid shopID"
// @Failure 400 {string} string "invalid itemID"
// @Failure 400 {string} string "invalid request body"
// @Failure 500 {string} string "error message"
// @Router /owner/shops/{shopID}/items/{itemID} [put]
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

// DeleteItemHandler godoc
// @Summary Delete item
// @Description Delete an item by ID
// @Tags items
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param shopID path int true "Shop ID"
// @Param itemID path int true "Item ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid itemID"
// @Failure 500 {string} string "error message"
// @Router /owner/shops/{shopID}/items/{itemID} [delete]
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
