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
    
    // Check if we have a database connection string
    if dbConnection == "" {
        log.Fatal("ERROR: No database connection string found! Set DATABASE_URL or DB_CONN environment variable")
    }
    
    log.Println("Database connection string found, attempting to connect...")
    
    db, err := database.InitDB(dbConnection)
    if err != nil {
        log.Printf("DATABASE CONNECTION FAILED: %v", err)
        log.Fatal("Cannot start application without database connection")
    }
    defer db.Close()
    
    log.Println("✓ Database connected successfully")

    // Initialize repositories
    log.Println("Initializing repositories...")
    productRepo := repositories.NewProductRepository(db)
    categoryRepo := repositories.NewCategoryRepository(db)
    
    // Initialize services
    log.Println("Initializing services...")
    productService := services.NewProductService(productRepo)
    
    // Initialize handlers
    log.Println("Initializing handlers...")
    productHandler := handlers.NewProductHandler(productService)
    categoryHandler := handlers.NewCategoryHandler(categoryRepo)
    
    // Setup routes
    log.Println("Setting up routes...")
    http.HandleFunc("/health", handlers.HealthCheckHandler)
    
    // Product routes
    http.HandleFunc("/api/product", productHandler.HandleProducts)
    http.HandleFunc("/api/product/", productHandler.HandleProductByID)
    
    // Category routes
    http.HandleFunc("/api/category", categoryHandler.HandleCategories)
    http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)
    
    // Start server
    log.Printf("=== Server starting on port %s ===", config.Port)
    fmt.Printf("Server started at :%s\n", config.Port)
    fmt.Println("Health check available at: /health")
    fmt.Println("API endpoints available at: /api/product, /api/category")
    
    if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
        log.Printf("SERVER ERROR: %v", err)
        log.Fatal("Failed to start HTTP server")
    }
}