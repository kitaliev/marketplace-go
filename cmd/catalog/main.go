package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kitaliev/marketplace-go/internal/repository/jsonfile"
)

func main() {
	fileName := flag.String("seed", "", "--seed=<path to file>")
	flag.Parse()
	if *fileName == "" {
		fmt.Println("no seed file was given")
		os.Exit(1)
	}

	service, err := jsonfile.Load(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idParam := r.PathValue("id")
		id, convErr := strconv.ParseInt(idParam, 10, 64)
		encoder := json.NewEncoder(w)
		if convErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeErr := encoder.Encode(map[string]string{
				"error": "invalid product id",
			})
			if encodeErr != nil {
				log.Printf("error responding error to product/%s", idParam)
			}
			return
		}
		product, getErr := service.GetByID(id)
		if getErr != nil {
			if errors.Is(getErr, jsonfile.ErrProductNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			getErr = encoder.Encode(map[string]string{
				"error": "product not found",
			})
			if getErr != nil {
				log.Printf("error responding error to product/%s", idParam)
			}
			return
		}
		encodeErr := encoder.Encode(product)
		if encodeErr != nil {
			log.Printf("error responding response to product/%s", idParam)
		}
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serveErr := server.ListenAndServe()

	if serveErr != nil {
		log.Printf("error during listen and serve %v", serveErr)
	}
}
