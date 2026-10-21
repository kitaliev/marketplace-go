package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitaliev/marketplace-go/internal/catalog/handler"
	"github.com/kitaliev/marketplace-go/internal/catalog/storage"
	"github.com/kitaliev/marketplace-go/internal/platform/httpserver"
	"github.com/kitaliev/marketplace-go/internal/platform/storage/memory"
)

func main() {
	ok := initServer()
	if !ok {
		os.Exit(1)
	}
	os.Exit(0)
}

func initServer() bool {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// parsing arguments
	args, err := httpserver.InitArguments("data/catalog.json", ":8080", 10*time.Second)
	if err != nil {
		slog.Error("error initializing program", "err", err)
		return false
	}

	// initializing database
	data, err := storage.Read(args.Storage)
	if err != nil {
		slog.Error("error initializing database", "err", err)
		return false
	}
	products := memory.New(data.Products)
	categories := memory.New(data.Categories)

	h := handler.New(products, categories)
	err = httpserver.InitServer(ctx, h.Routes(), args)
	if err != nil {
		slog.Error("error initializing server", "err", err)
		return false
	}
	return true
}
