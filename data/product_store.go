package data

import (
	"errors"
	"kasir-api/models"
	"sync"
)

type ProductStore struct {
	products []models.Product
	nextID   int
	mu       sync.RWMutex
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: []models.Product{
			{ID: 1, Nama: "Produk 1", Harga: 10000, Stok: 10, CategoryID: 1},
			{ID: 2, Nama: "Produk 2", Harga: 20000, Stok: 20, CategoryID: 2},
			{ID: 3, Nama: "Produk 3", Harga: 30000, Stok: 30, CategoryID: 3},
		},
		nextID: 4,
	}
}

func (s *ProductStore) GetAll() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.products
}

func (s *ProductStore) GetByID(id int) (*models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, product := range s.products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, errors.New("product not found")
}

func (s *ProductStore) GetByCategory(categoryID int) []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var result []models.Product
	for _, product := range s.products {
		if product.CategoryID == categoryID {
			result = append(result, product)
		}
	}
	return result
}

func (s *ProductStore) Create(product models.Product) models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if product.ID == 0 {
		product.ID = s.nextID
		s.nextID++
	}
	s.products = append(s.products, product)
	return product
}

func (s *ProductStore) Update(id int, product models.Product) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, p := range s.products {
		if p.ID == id {
			product.ID = id
			s.products[i] = product
			return &product, nil
		}
	}
	return nil, errors.New("product not found")
}

func (s *ProductStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, product := range s.products {
		if product.ID == id {
			s.products = append(s.products[:i], s.products[i+1:]...)
			return nil
		}
	}
	return errors.New("product not found")
}