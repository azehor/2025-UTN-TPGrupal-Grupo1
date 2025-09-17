package http

import (
	"log/slog"
	"net/http"

	"quepc/api/internal/carreras"

	"github.com/go-chi/chi/v5"
)

// Un Server incluye una referencia a cada modelo de los endpoints a exponer
// Estos modelos deben ser acompañados por una declaracion de la interfaz que cumplen
type Server struct {
	carreras *carreras.Carreras
}

// Funcion de inicializacion del server, los argumentos seran los modelos a usar
func New(c *carreras.Carreras) *Server {
	return &Server{
		carreras: c,
	}
}

// Aqui se agregan todas las rutas asociandolas con la funcion que resuelve que se envia
func (s *Server) AddRoutes(r *chi.Mux) {
	//ejemplo de uso
	r.Get("/", s.getHomepage)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/carreras", s.carreras.List)
		r.Post("/carreras", s.carreras.Create)
		r.Get("/carreras/{id}", s.carreras.Read)
		r.Put("/carreras/{id}", s.carreras.Update)
		r.Delete("/carreras/{id}", s.carreras.Delete)
	})
}

func (s *Server) getHomepage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Homepage requested")
}
