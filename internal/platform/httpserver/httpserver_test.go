package httpserver

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestInitServerSamePortErrors(t *testing.T) {
	var lc net.ListenConfig
	busy, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unable to occupy a port for test: %v", err)
	}
	defer func() { _ = busy.Close() }()

	args := Arguments{
		Addr:            busy.Addr().String(),
		ShutdownTimeout: 10 * time.Second,
	}
	err = InitServer(context.Background(), http.NotFoundHandler(), &args)
	if err == nil {
		t.Errorf("expected error on initializing server on the same port, no error received")
	}
}
