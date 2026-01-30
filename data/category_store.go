package data

import (
	"errors"
	"kasir-api/models"
	"sync"
)

type CategoryStore struct {
	categories []models.Category
	nextID     int
	mu         sync.RWMutex
}

func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		categories: []models.Category{},
		nextID:     1,
	}
}

func (s *CategoryStore) GetAll() []models.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.categories
}

func (s *CategoryStore) GetByID(id int) (*models.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, category := range s.categories {
		if category.ID == id {
			return &category, nil
		}
	}
	return nil, errors.New("category not found")
}

func (s *CategoryStore) Create(category models.Category) models.Category {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if category.ID == 0 {
		category.ID = s.nextID
		s.nextID++
	}
	s.categories = append(s.categories, category)
	return category
}

func (s *CategoryStore) Update(id int, category models.Category) (*models.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, c := range s.categories {
		if c.ID == id {
			category.ID = id
			s.categories[i] = category
			return &category, nil
		}
	}
	return nil, errors.New("category not found")
}

func (s *CategoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, category := range s.categories {
		if category.ID == id {
			s.categories = append(s.categories[:i], s.categories[i+1:]...)
			return nil
		}
	}
	return errors.New("category not found")
}