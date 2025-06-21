package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env")
	}

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL"),
	)

	// start DB
	InitDB(dsn)

	mux := http.NewServeMux()

	// start server
	mux.HandleFunc("/ratings", GetRatingsHandler)
	mux.HandleFunc("/sync", FetchAndStoreRatings)
	mux.HandleFunc("/ratings/best", GetBestRatingHandler)

	handler := cors.Default().Handler(mux)
	fmt.Println("Server run in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
