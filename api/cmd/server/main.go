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
	"quepc/api/internal/recomendaciones"

	"quepc/api/internal/softwares"
	sStore "quepc/api/internal/softwares/store"

	transport "quepc/api/internal/transport/http"

	//external library imports
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	dbPostgresql "quepc/api/internal/db"
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

	dbPostgresql.InitDB()

	carrerasStore := cStore.New(dbPostgresql.DB)
	carreras := carreras.New(carrerasStore)

	softwaresStore := sStore.New(dbPostgresql.DB)
	softwares := softwares.New(softwaresStore)

	recomendaciones := recomendaciones.New()

	httpServer := transport.New(carreras, softwares, recomendaciones)
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
