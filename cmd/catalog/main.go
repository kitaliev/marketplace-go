package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kitaliev/marketplace-go/internal/handler"
	"github.com/kitaliev/marketplace-go/internal/repository/memory"
	"github.com/kitaliev/marketplace-go/internal/seed"
)

func main() {
	fileName := flag.String("seed", "data/seed.json", "--seed=<path to file>")
	flag.Parse()
	if *fileName == "" {
		fmt.Println("no seed file was given")
		os.Exit(1)
	}
	data, err := seed.Read(*fileName)
	if err != nil {
		slog.Error("error initializing memory", "err", err)
		log.Fatal(err)
	}

	products := memory.New(data.Products)
	categories := memory.New(data.Categories)
	h := handler.New(products, categories)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           h.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
