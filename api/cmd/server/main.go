package server

import (
	//stdlib imports
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	//internal project imports
	transport "quepc/api/internal/transport/http"

	//external library imports
	"github.com/go-chi/chi/v5"
)

func Start() {
	r := chi.NewRouter()
	httpServer := transport.New()
	httpServer.AddRoutes(r)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		shutdown()
		os.Exit(1)
	}()

	http.ListenAndServe(":8080", r)
	shutdown()
}

func shutdown() {
	fmt.Printf("Shutting down...\n")
}
