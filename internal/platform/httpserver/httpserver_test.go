package httpserver

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestInitServerSamePortErrors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	args := Arguments{
		Addr:            ":8080",
		ShutdownTimeout: 10 * time.Second,
	}
	err := InitServer(context.Background(), mux, &args)
	if err != nil {
		t.Fatalf("unable to start server for test %v:", err)
	}
	err = InitServer(context.Background(), mux, &args)
	if err == nil {
		t.Errorf("expected error on initializing server on the same port, no error received")
	}
}
