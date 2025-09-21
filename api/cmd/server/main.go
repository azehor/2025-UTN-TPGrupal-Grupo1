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

	"quepc/api/internal/softwares"
	sStore "quepc/api/internal/softwares/store"

	transport "quepc/api/internal/transport/http"

	//external library imports
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Start(port string) {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	carrerasStore := cStore.New()
	carreras := carreras.New(carrerasStore)

	softwaresStore := sStore.New()
	softwares := softwares.New(softwaresStore)

	httpServer := transport.New(carreras, softwares)
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
