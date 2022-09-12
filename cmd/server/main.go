package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmvrus/json-api/internal/endpoint"
	"github.com/tmvrus/json-api/internal/storage/balance"
)

func main() {
	storage := balance.NewStorage()
	server := &http.Server{
		Addr:    ":8080",
		Handler: endpoint.NewEndpoint(storage),
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		s := <-sigs
		log.Printf("got external signal: %q, shoutdown", s.String())
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shoutdown server: %s", err.Error())
		}
	}()

	log.Printf("start serving on :8080")
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server shoutdown with error: %s", err.Error())
		}
	}

	log.Printf("shoutdown successful")
}
