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
	fileName := flag.String("seed", "data/seed.json", "--seed=<path to file>")
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
		id, err := strconv.ParseInt(idParam, 10, 64)
		encoder := json.NewEncoder(w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := encoder.Encode(map[string]string{
				"error": "invalid product id",
			})
			if err != nil {
				log.Printf("error responding error to product/%s", idParam)
			}
			return
		}
		product, err := service.GetByID(id)
		if err != nil && errors.Is(err, jsonfile.ErrProductNotFound) {
			w.WriteHeader(http.StatusNotFound)
			err := encoder.Encode(map[string]string{
				"error": "product not found",
			})
			if err != nil {
				log.Printf("error responding error to product/%s", idParam)
			}
			return
		}
		if err != nil && errors.Is(err, jsonfile.ErrIncorrectID) {
			w.WriteHeader(http.StatusBadRequest)
			err := encoder.Encode(map[string]string{
				"error": "invalid product id",
			})
			if err != nil {
				log.Printf("error responding error to product/%s", idParam)
			}
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := encoder.Encode(map[string]string{
				"error": "internal server error",
			})
			if err != nil {
				log.Printf("error responding error to product/%s", idParam)
			}
			return
		}
		err = encoder.Encode(product)
		if err != nil {
			log.Printf("error responding response to product/%s", idParam)
		}
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("error during listen and serve %v", err)
	}
}
