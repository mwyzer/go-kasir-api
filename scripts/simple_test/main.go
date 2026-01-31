package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbConn := "postgresql://postgres.obnibaapisdtgisagmus:NsZgVB8d1ndQ2gej@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require"
	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var one int
	err = db.QueryRow("SELECT 1").Scan(&one)
	if err != nil {
		log.Fatalf("QUERY ERROR: %v", err)
	}
	fmt.Printf("RESULT: %d\n", one)
}
