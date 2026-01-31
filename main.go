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
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// Immediate startup log
	log.Println("=== Kasir API Starting ===")

	// Try to read .env file if it exists, but don't fail if it doesn't
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found (this is normal for Zeabur): %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Config port
	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("PORT")
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("Server will use port: %s", port)

	config := Config{
		Port:   port,
		DBConn: viper.GetString("DB_CONN"),
	}

	// Config database
	dbConnection := os.Getenv("DATABASE_URL")
	if dbConnection == "" {
		dbConnection = config.DBConn
	}

	// Setup basic health check first (works without DB)
	log.Println("Setting up health check endpoint...")
	http.HandleFunc("/health", handlers.HealthCheckHandler)

	// Try to initialize database
	var db *sql.DB
	var productRepo *repositories.ProductRepository
	var categoryRepo *repositories.CategoryRepository

	if dbConnection == "" {
		log.Println("WARNING: No database connection string found!")
		log.Println("Service will start but API endpoints will not work")
		log.Println("Set DATABASE_URL environment variable in Zeabur")
	} else {
		log.Println("Database connection string found, attempting to connect...")

		var err error
		db, err = database.InitDB(dbConnection)
		if err != nil {
			log.Printf("WARNING: Database connection failed: %v", err)
			log.Println("Service will start but API endpoints will not work")
		} else {
			defer db.Close()
			log.Println("✓ Database connected successfully")

			// Initialize repositories only if DB is connected
			log.Println("Initializing repositories...")
			productRepo = repositories.NewProductRepository(db)
			categoryRepo = repositories.NewCategoryRepository(db)

			// Initialize services
			log.Println("Initializing services...")
			productService := services.NewProductService(productRepo)

			// Initialize handlers
			log.Println("Initializing handlers...")
			productHandler := handlers.NewProductHandler(productService)
			categoryHandler := handlers.NewCategoryHandler(categoryRepo)

			// Setup API routes only if DB is connected
			log.Println("Setting up API routes...")
			http.HandleFunc("/api/product", productHandler.HandleProducts)
			http.HandleFunc("/api/product/", productHandler.HandleProductByID)
			http.HandleFunc("/api/category", categoryHandler.HandleCategories)
			http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)
		}
	}

	// Start server (this will always run)
	log.Printf("=== Server starting on port %s ===", config.Port)
	fmt.Printf("Server started at :%s\n", config.Port)
	fmt.Println("Health check available at: /health")
	if db != nil {
		fmt.Println("API endpoints available at: /api/product, /api/category")
	} else {
		fmt.Println("API endpoints NOT available - database not connected")
	}

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Printf("SERVER ERROR: %v", err)
		log.Fatal("Failed to start HTTP server")
	}
}