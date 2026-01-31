package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback for local testing if needed, though env var should be set
		// Attempting to read .env manually if os.Getenv is empty (simplified check)
		// For now assuming it runs in context with env var set or user provides it
		fmt.Println("DATABASE_URL not set in env")

        // Hardcoding based on previous view_file of .env just for this diagnostic if env var fails
        dbURL = "postgresql://postgres.obnibaapisdtgisagmus:NsZgVB8d1ndQ2gej@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require"
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_name = 'categories'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Columns in categories table:")
	found := false
	for rows.Next() {
		var colName string
		if err := rows.Scan(&colName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(colName)
		found = true
	}
    if !found {
        fmt.Println("No columns found or table 'categories' does not exist.")
    }
}
