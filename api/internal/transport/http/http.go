package http

import (
	"log/slog"
	"net/http"

	carrerasoftware "quepc/api/internal/carreraSoftware"
	"quepc/api/internal/carreras"
	"quepc/api/internal/componentes"
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
	componentes     *componentes.Componentes
}

// Funcion de inicializacion del server, los argumentos seran los modelos a usar
func New(c *carreras.Carreras,
	s *softwares.Softwares,
	r *recomendaciones.Recomendaciones,
	cs *carrerasoftware.CarreraSoftwares,
	cmp *componentes.Componentes) *Server {
	return &Server{
		carreras:        c,
		softwares:       s,
		recomendaciones: r,
		carreraSoftware: cs,
		componentes:     cmp,
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

		r.Route("/componentes", func(r chi.Router) {
			// Discos
			r.Get("/discos", s.componentes.Discos.List)
			r.Post("/discos", s.componentes.Discos.Create)
			r.Get("/discos/{id}", s.componentes.Discos.Read)
			r.Put("/discos/{id}", s.componentes.Discos.Update)
			r.Delete("/discos/{id}", s.componentes.Discos.Delete)

			// Gabinetes
			r.Get("/gabinetes", s.componentes.Gabinetes.List)
			r.Post("/gabinetes", s.componentes.Gabinetes.Create)
			r.Get("/gabinetes/{id}", s.componentes.Gabinetes.Read)
			r.Put("/gabinetes/{id}", s.componentes.Gabinetes.Update)
			r.Delete("/gabinetes/{id}", s.componentes.Gabinetes.Delete)

			// GPUs
			r.Get("/gpus", s.componentes.GPUs.List)
			r.Post("/gpus", s.componentes.GPUs.Create)
			r.Get("/gpus/{id}", s.componentes.GPUs.Read)
			r.Put("/gpus/{id}", s.componentes.GPUs.Update)
			r.Delete("/gpus/{id}", s.componentes.GPUs.Delete)

			// Motherboards
			r.Get("/motherboards", s.componentes.Motherboards.List)
			r.Post("/motherboards", s.componentes.Motherboards.Create)
			r.Get("/motherboards/{id}", s.componentes.Motherboards.Read)
			r.Put("/motherboards/{id}", s.componentes.Motherboards.Update)
			r.Delete("/motherboards/{id}", s.componentes.Motherboards.Delete)

			// Procesadores
			r.Get("/procesadores", s.componentes.Procesadores.List)
			r.Post("/procesadores", s.componentes.Procesadores.Create)
			r.Get("/procesadores/{id}", s.componentes.Procesadores.Read)
			r.Put("/procesadores/{id}", s.componentes.Procesadores.Update)
			r.Delete("/procesadores/{id}", s.componentes.Procesadores.Delete)

			// PSUs
			r.Get("/psus", s.componentes.PSUs.List)
			r.Post("/psus", s.componentes.PSUs.Create)
			r.Get("/psus/{id}", s.componentes.PSUs.Read)
			r.Put("/psus/{id}", s.componentes.PSUs.Update)
			r.Delete("/psus/{id}", s.componentes.PSUs.Delete)

			// RAMs
			r.Get("/rams", s.componentes.RAMs.List)
			r.Post("/rams", s.componentes.RAMs.Create)
			r.Get("/rams/{id}", s.componentes.RAMs.Read)
			r.Put("/rams/{id}", s.componentes.RAMs.Update)
			r.Delete("/rams/{id}", s.componentes.RAMs.Delete)

		})
	})
}

func (s *Server) getHomepage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Homepage requested")
}
