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
	dbConn := "postgresql://postgres.obnibaapisdtgisagmus:NsZgVB8d1ndQ2gej@aws-1-ap-south-1.pooler.supabase.com:6543/postgres"
	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT * FROM products LIMIT 0")
	if err != nil {
		fmt.Printf("QUERY ERROR: %v\n", err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("COLS ERROR: %v\n", err)
		return
	}

	fmt.Printf("ACTUAL_COLS:%v\n", cols)
}
