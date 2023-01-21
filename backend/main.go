package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	. "github.com/gerardrbentley/places/handler"
	. "github.com/gerardrbentley/places/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupHandlers(placeService PlaceService) *http.ServeMux {
	h := http.NewServeMux()
	h.Handle("/place", PlaceHandler(placeService))

	return h
}

func main() {
	log.Println("Starting Up....")

	log.Println("Connecting to postgres...")
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DB_CONNECTION"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	placeService := NewPgPlaceService(dbpool)

	h := setupHandlers(placeService)
	log.Fatal(http.ListenAndServe(":5000", h))
}
