package handlers

import (
	"encoding/json"
	"errors"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.service.Create(&newProduct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	product, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	var updatedProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	updatedProduct.ID = id
	err = h.service.Update(&updatedProduct)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProduct)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.service.Delete(id); err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) extractID(r *http.Request) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	idStr = strings.Trim(idStr, "/")
	if idStr == "" {
		return 0, errors.New("Product ID is required in URL (e.g., /api/product/100)")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("Invalid Product ID format")
	}
	return id, nil
}

// Handler untuk routing /api/product
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListProducts(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler untuk routing /api/product/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	// Jika path setelah prefix kosong (misal: /api/product/), arahkan ke HandleProducts
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	idStr = strings.Trim(idStr, "/")
	if idStr == "" {
		h.HandleProducts(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetProductByID(w, r)
	case http.MethodPut, http.MethodPost:
		h.UpdateProduct(w, r)
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}