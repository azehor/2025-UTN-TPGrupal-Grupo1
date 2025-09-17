package server

import (
	//stdlib imports
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	//internal project imports
	"quepc/api/internal/carreras"
	cStore "quepc/api/internal/carreras/store"
	transport "quepc/api/internal/transport/http"

	//external library imports
	"github.com/go-chi/chi/v5"
)

func Start(port string) {
	r := chi.NewRouter()

	carrerasStore := cStore.New()
	carreras := carreras.New(carrerasStore)

	httpServer := transport.New(carreras)
	httpServer.AddRoutes(r)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		shutdown()
		os.Exit(1)
	}()

	fmt.Printf("Running server on port: %v", port)
	http.ListenAndServe(port, r)
	shutdown()
}

func shutdown() {
	fmt.Printf("Shutting down...\n")
}
