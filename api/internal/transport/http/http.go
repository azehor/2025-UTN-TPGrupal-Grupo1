package http

import (
	"log/slog"
	"net/http"

	"quepc/api/internal/carreras"
	"quepc/api/internal/recomendaciones"
	"quepc/api/internal/softwares"

	"github.com/go-chi/chi/v5"
)

// Un Server incluye una referencia a cada modelo de los endpoints a exponer
// Estos modelos deben ser acompañados por una declaracion de la interfaz que cumplen
type Server struct {
	carreras        *carreras.Carreras
	softwares       *softwares.Softwares
	recomendaciones *recomendaciones.Recomendaciones
}

// Funcion de inicializacion del server, los argumentos seran los modelos a usar
func New(c *carreras.Carreras, s *softwares.Softwares) *Server {
	return &Server{
		carreras:  c,
		softwares: s,
	}
}

// Aqui se agregan todas las rutas asociandolas con la funcion que resuelve que se envia
func (s *Server) AddRoutes(r *chi.Mux) {
	//ejemplo de uso
	r.Get("/", s.getHomepage)
	r.Route("/v1", func(r chi.Router) {
		// Carreras
		r.Get("/carreras", s.carreras.List)
		r.Post("/carreras", s.carreras.Create)
		r.Get("/carreras/{id}", s.carreras.Read)
		r.Put("/carreras/{id}", s.carreras.Update)
		r.Delete("/carreras/{id}", s.carreras.Delete)

		// Softwares
		r.Get("/softwares", s.softwares.List)
		r.Post("/softwares", s.softwares.Create)
		r.Get("/softwares/{id}", s.softwares.Read)
		r.Put("/softwares/{id}", s.softwares.Update)
		r.Delete("/softwares/{id}", s.softwares.Delete)

		//Recomendaciones
		r.Post("/recomendaciones-softwares", r.recomendaciones.ListSoftwares)
		r.Post("/recomendaciones-carrera", r.recomendaciones.ListCarrera)
	})
}

func (s *Server) getHomepage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Homepage requested")
}
