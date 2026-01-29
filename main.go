package main

import (
	"fmt"
	"kasir-api/data"
	"kasir-api/handlers"
	"net/http"
)

func main() {
	// Initialize stores
	productStore := data.NewProductStore()
	categoryStore := data.NewCategoryStore()
	
	// Initialize handlers
	productHandler := handlers.NewProductHandler(productStore)
	categoryHandler := handlers.NewCategoryHandler(categoryStore)
	
	// Setup routes
	http.HandleFunc("/health", handlers.HealthCheckHandler)
	
	// Product routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)
	
	// Category routes
	http.HandleFunc("/api/category", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)
	
	// Start server
	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}