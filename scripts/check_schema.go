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
		fmt.Println("DATABASE_URL not set")
		return
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
	for rows.Next() {
		var colName string
		if err := rows.Scan(&colName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(colName)
	}
}
