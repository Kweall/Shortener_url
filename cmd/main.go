package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"ozon/internal/handlers"
	"ozon/internal/storage"
)

func main() {
	storageType := flag.String("storage", "memory", "Storage type: memory or postgres")
	pgConnStr := flag.String("pg_conn_str", "", "PostgreSQL connection string")
	flag.Parse()

	var store storage.Storage
	var err error

	switch *storageType {
	case "memory":
		store = storage.NewMemoryStorage()
	case "postgres":
		store, err = storage.NewPostgresStorage(*pgConnStr)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
		}
	default:
		log.Fatal("Unknown storage type")
	}

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenURLHandler(w, r, store)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RedirectHandler(w, r, store)
	})

	fmt.Println("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
