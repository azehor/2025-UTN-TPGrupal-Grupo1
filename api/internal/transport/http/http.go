package http

import (
	"log/slog"
	"net/http"

	carrerasoftware "quepc/api/internal/carreraSoftware"
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
	carreraSoftware *carrerasoftware.CarreraSoftwares
}

// Funcion de inicializacion del server, los argumentos seran los modelos a usar
func New(c *carreras.Carreras, s *softwares.Softwares, r *recomendaciones.Recomendaciones, cs *carrerasoftware.CarreraSoftwares) *Server {
	return &Server{
		carreras:        c,
		softwares:       s,
		recomendaciones: r,
		carreraSoftware: cs,
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
		r.Post("/recomendaciones-softwares", s.recomendaciones.RecomendacionSoftware)
		r.Get("/recomendaciones-carrera/{id}", s.recomendaciones.RecomendacionCarrera)

		// Carrera - Softwares
		r.Get("/carreras/{id}/softwares", s.carreraSoftware.ListByCarrera)
		r.Post("/carreras/{id}/softwares/{software_id}", s.carreraSoftware.Create)
		r.Delete("/carreras/{id}/softwares/{software_id}", s.carreraSoftware.Delete)
	})
}

func (s *Server) getHomepage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Homepage requested")
}
