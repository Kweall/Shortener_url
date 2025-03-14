package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"ozon/internal/handlers"
	"ozon/internal/service"
	"ozon/internal/storage/memory"
	"ozon/internal/storage/postgres"
)

func main() {
	storageType := flag.String("storage", "memory", "Storage type: memory or postgres")
	pgConnStr := flag.String("pg_conn_str", "", "PostgreSQL connection string")
	flag.Parse()

	var (
		storage service.Storage
		err     error
	)

	switch *storageType {
	case "memory":
		storage = memory.NewMemoryStorage()
	case "postgres":
		storage, err = postgres.NewPostgresStorage(*pgConnStr)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
		}
	default:
		log.Fatal("Unknown storage type")
	}

	service := service.NewService(storage)

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenURLHandler(w, r, service)
	})

	http.HandleFunc("/redirect/", func(w http.ResponseWriter, r *http.Request) {
		handlers.RedirectHandler(w, r, service)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.OriginalURLHandler(w, r, service)
	})

	fmt.Println("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
