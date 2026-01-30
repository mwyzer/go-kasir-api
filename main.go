package main

import (
	"fmt"
	"kasir-api/data"
	"kasir-api/database"
	"kasir-api/handlers"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//config port
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	//config database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


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
	fmt.Println("Server started at :"+config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		fmt.Println(err)
	}
}