package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("=== STARTING ===")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("PORT=%s", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	log.Printf("Listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}