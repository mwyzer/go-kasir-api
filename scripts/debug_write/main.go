package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Using the connection string from .env
	dbConn := "postgresql://postgres.obnibaapisdtgisagmus:NsZgVB8d1ndQ2gej@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require"
	
	fmt.Println("Connecting to database...")
	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		log.Fatalf("FAILED TO OPEN: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("FAILED TO PING: %v", err)
	}
	fmt.Println("Connected successfully!")

	// Test Insert
	fmt.Println("Attempting to insert a test product...")
	query := "INSERT INTO products (nama, harga, stok, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	var lastID int
	err = db.QueryRowContext(ctx, query, "Debug Product", 999, 5, 1).Scan(&lastID)
	if err != nil {
		fmt.Printf("INSERT ERROR: %v\n", err)
		
		// Let's try to see if the table exists and what columns it has
		fmt.Println("\nInspecting table structure...")
		rows, err := db.QueryContext(ctx, "SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'products'")
		if err != nil {
			fmt.Printf("SCHEMA INSPECTION ERROR: %v\n", err)
		} else {
			defer rows.Close()
			fmt.Println("Columns in 'products' table:")
			for rows.Next() {
				var colName, dataType string
				rows.Scan(&colName, &dataType)
				fmt.Printf("- %s (%s)\n", colName, dataType)
			}
		}
	} else {
		fmt.Printf("INSERT SUCCESS! ID: %d\n", lastID)
		
		// Clean up
		_, err = db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", lastID)
		if err != nil {
			fmt.Printf("CLEANUP ERROR: %v\n", err)
		} else {
			fmt.Println("Cleanup successful.")
		}
	}
}
