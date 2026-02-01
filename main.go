package main

import (
	"database/sql"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func main() {
	log.Println("=== Kasir API Starting ===")

	// Try to read .env file if it exists (local dev only)
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found (normal on Zeabur)")
	} else {
		log.Println(".env file loaded successfully")
	}
	viper.AutomaticEnv()

	// Port — Zeabur injects this automatically
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Using port: %s", port)

	// Database — Zeabur injects DATABASE_URL automatically
	dbConnection := os.Getenv("DATABASE_URL")
	if dbConnection == "" {
		log.Println("WARNING: DATABASE_URL not found!")
	} else {
		log.Println("DATABASE_URL found")
	}

	// Health check always available
	http.HandleFunc("/health", handlers.HealthCheckHandler)

	// Try database connection
	var db *sql.DB
	if dbConnection != "" {
		log.Println("Connecting to database...")
		var err error
		db, err = database.InitDB(dbConnection)
		if err != nil {
			log.Printf("DATABASE CONNECTION FAILED: %v", err)
		} else {
			defer db.Close()
			log.Println("✓ Database connected successfully")
		}
	}

	// Setup API routes only if DB is connected
	if db != nil {
		log.Println("Setting up API routes...")
		productRepo := repositories.NewProductRepository(db)
		categoryRepo := repositories.NewCategoryRepository(db)
		productService := services.NewProductService(productRepo)
		productHandler := handlers.NewProductHandler(productService)
		categoryHandler := handlers.NewCategoryHandler(categoryRepo)

		http.HandleFunc("/api/product", productHandler.HandleProducts)
		http.HandleFunc("/api/product/", productHandler.HandleProductByID)
		http.HandleFunc("/api/category", categoryHandler.HandleCategories)
		http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)
		log.Println("✓ API routes registered")
	} else {
		log.Println("WARNING: API routes NOT available — no database")
	}

	// Start server
	log.Printf("=== Server starting on port %s ===", port)
	fmt.Printf("Server started at :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("SERVER ERROR: %v", err)
	}
}