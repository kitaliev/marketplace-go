package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitaliev/marketplace-go/internal/handler"
	"github.com/kitaliev/marketplace-go/internal/repository/memory"
	"github.com/kitaliev/marketplace-go/internal/seed"
)

type arguments struct {
	seed            string
	addr            string
	shutdownTimeout time.Duration
}

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
	args, err := initArguments()
	if err != nil {
		slog.Error("error initializing program", "err", err)
		return false
	}

	// initializing database
	data, err := seed.Read(args.seed)
	if err != nil {
		slog.Error("error initializing database", "err", err)
		return false
	}
	products := memory.New(data.Products)
	categories := memory.New(data.Categories)

	// initializing server
	h := handler.New(products, categories)
	server := &http.Server{
		Addr:              args.addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// starting server and listen for error
	srvErr := make(chan error, 1) // buffered so on server stop goroutine can finish after putting error into channel
	go func() {
		slog.Info("server started at port " + args.addr)
		srvErr <- server.ListenAndServe()
	}()

	// ether system interrupts or error happens
	// if system interrupts then we're gracefully shutting down all connections and server
	select {
	case <-ctx.Done():
		slog.Info("system interrupted program execution")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), args.shutdownTimeout)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("error during graceful shutdown", "err", err)
			return false
		}
		slog.Info("server stopped")
		return true
	case err := <-srvErr:
		slog.Error("error happened during program execution", "err", err)
		return false
	}
}

func initArguments() (*arguments, error) {
	file := flag.String("seed", "data/seed.json", "--seed=<path to file>")
	addr := flag.String("addr", ":8080", "--addr=<path to file>")
	timeout := flag.Duration("shutdown-timeout", 10*time.Second, "--shutdown-timeout=5s")
	flag.Parse()
	if *file == "" {
		return nil, errors.New("no seed file was given")
	}
	if *addr == "" {
		return nil, errors.New("no address was given")
	}
	if *timeout <= 0 {
		return nil, errors.New("timeout should be more than 0s")
	}
	return &arguments{
		seed:            *file,
		addr:            *addr,
		shutdownTimeout: *timeout,
	}, nil
}
