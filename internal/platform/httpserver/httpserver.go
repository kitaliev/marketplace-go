package httpserver

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Arguments struct {
	Storage         string
	Addr            string
	ShutdownTimeout time.Duration
}

func InitArguments(defStorage string, defAddr string, defTimeout time.Duration) (*Arguments, error) {
	file := flag.String("storage", defStorage, "--storage=<path to file>")
	addr := flag.String("addr", defAddr, "--addr=<path to file>")
	timeout := flag.Duration("shutdown-timeout", defTimeout, "--shutdown-timeout=5s")
	flag.Parse()
	if *file == "" {
		return nil, errors.New("no storage file was given")
	}
	if *addr == "" {
		return nil, errors.New("no address was given")
	}
	if *timeout <= 0 {
		return nil, errors.New("timeout should be more than 0s")
	}
	return &Arguments{
		Storage:         *file,
		Addr:            *addr,
		ShutdownTimeout: *timeout,
	}, nil
}

func InitServer(ctx context.Context, routes http.Handler, args *Arguments) error {
	// initializing server
	server := &http.Server{
		Addr:              args.Addr,
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// starting server and listen for error
	srvErr := make(chan error, 1) // buffered so on server stop goroutine can finish after putting error into channel
	go func() {
		var lc net.ListenConfig
		ln, err := lc.Listen(context.Background(), "tcp", args.Addr)
		if err != nil {
			srvErr <- err
			return
		}
		slog.Info("server started at port " + args.Addr)
		srvErr <- server.Serve(ln)
	}()

	// either system interrupts or error happens
	// if system interrupts then we're gracefully shutting down all connections and server
	select {
	case <-ctx.Done():
		slog.Info("system interrupted program execution")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), args.ShutdownTimeout)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("error during graceful shutdown", "err", err)
			return err
		}
		slog.Info("server stopped")
		return nil
	case err := <-srvErr:
		slog.Error("error happened during program execution", "err", err)
		return err
	}
}
