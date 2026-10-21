package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitaliev/marketplace-go/internal/platform/httpserver"
	"github.com/kitaliev/marketplace-go/internal/platform/storage/memory"
	"github.com/kitaliev/marketplace-go/internal/prices/handler"
	"github.com/kitaliev/marketplace-go/internal/prices/storage"
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
	args, err := httpserver.InitArguments("data/prices.json", ":8081", 10*time.Second)
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
	prices := memory.New(data.Prices)

	h := handler.New(prices)
	err = httpserver.InitServer(ctx, h.Routes(), args)
	if err != nil {
		slog.Error("error initializing server", "err", err)
		return false
	}
	return true
}
