package main

import (
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	// Try to read .env file if it exists, but don't fail if it doesn't  
	viper.SetConfigFile(".env")  
	viper.ReadInConfig() // Ignore error if file doesn't exist  
	viper.AutomaticEnv()  

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//config port
	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("PORT")
	}
	if port == "" {
		port = "8080"
	}

	config := Config{
		Port:   port,
		DBConn: viper.GetString("DB_CONN"),
	}

	//config database
	dbConnection := os.Getenv("DATABASE_URL")
	if dbConnection == "" {
		dbConnection = config.DBConn
	}
	db, err := database.InitDB(dbConnection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	// Initialize stores
	// productStore := data.NewProductStore() // No longer used for products
	// categoryStore := data.NewCategoryStore()
	
	// Initialize repositories
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	
	// Initialize services
	productService := services.NewProductService(productRepo)
	
	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	
	// Setup routes
	http.HandleFunc("/health", handlers.HealthCheckHandler)
	
	// Product routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)
	
	// Category routes
	http.HandleFunc("/api/category", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)
	
	// Start server
	fmt.Println("Server started at :"+config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		fmt.Println(err)
	}
}