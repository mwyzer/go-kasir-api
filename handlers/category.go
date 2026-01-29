package handlers

import (
	"encoding/json"
	"kasir-api/data"
	"kasir-api/models"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	store *data.CategoryStore
}

func NewCategoryHandler(store *data.CategoryStore) *CategoryHandler {
	return &CategoryHandler{store: store}
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	category := h.store.Create(newCategory)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractID(r)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	
	category, err := h.store.GetByID(id)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractID(r)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	
	var updatedCategory models.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	category, err := h.store.Update(id, updatedCategory)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractID(r)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	
	if err := h.store.Delete(id); err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) extractID(r *http.Request) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	return strconv.Atoi(idStr)
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListCategories(w, r)
	case http.MethodPost:
		h.CreateCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetCategoryByID(w, r)
	case http.MethodPut, http.MethodPost:
		h.UpdateCategory(w, r)
	case http.MethodDelete:
		h.DeleteCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}