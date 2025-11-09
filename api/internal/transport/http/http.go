package http

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"quepc/api/internal/auth"
	carrerasoftware "quepc/api/internal/carreraSoftware"
	"quepc/api/internal/carreras"
	"quepc/api/internal/componentes"
	"quepc/api/internal/middleware"
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
	auth            *auth.Handler
}

// Funcion de inicializacion del server, los argumentos seran los modelos a usar
func New(c *carreras.Carreras,
	s *softwares.Softwares,
	r *recomendaciones.Recomendaciones,
	cs *carrerasoftware.CarreraSoftwares,
	cmp *componentes.Componentes,
	auth *auth.Handler) *Server {
	return &Server{
		carreras:        c,
		softwares:       s,
		recomendaciones: r,
		carreraSoftware: cs,
		componentes:     cmp,
		auth:            auth,
	}
}

// Aqui se agregan todas las rutas asociandolas con la funcion que resuelve que se envia
func (s *Server) AddRoutes(r *chi.Mux) {
	r.Get("/", s.getHomepage)

	workdir, _ := os.Getwd()
	// Servir toda la carpeta /static
	staticDir := http.Dir(filepath.Join(workdir, "static"))
	FileServer(r, "/static", staticDir)
	// Compatibilidad con rutas antiguas /static/images
	imgDir := http.Dir(filepath.Join(workdir, "static/images"))
	FileServer(r, "/static/images", imgDir)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeMiddleware)

		// PÚBLICAS
		// AUTH
		r.Post("/auth/login", s.auth.Login)

		// ADMIN SOLO PUEDE CREAR USUARIOS
		r.Route("/admin", func(ar chi.Router) {
			ar.Use(auth.AdminOnly)
			ar.Post("/users", s.auth.AdminCreateUser)
		})

		// Carreras (solo listado público)
		r.Get("/carreras", s.carreras.List)

		// Softwares (solo listado público)
		r.Get("/softwares", s.softwares.List)

		// Carrera-Softwares (solo listado público)
		r.Get("/carreras/{id}/softwares", s.carreraSoftware.ListByCarrera)

		// Recomendaciones (públicas)
		r.Post("/recomendaciones-softwares", s.recomendaciones.RecomendacionSoftware)
		r.Get("/recomendaciones-carrera/{id}", s.recomendaciones.RecomendacionCarrera)

		// PROTEGIDAS
		r.Group(func(pr chi.Router) {
			pr.Use(auth.AuthMiddleware)

			// Carreras
			pr.Post("/carreras", s.carreras.Create)
			pr.Get("/carreras/{id}", s.carreras.Read)
			pr.Put("/carreras/{id}", s.carreras.Update)
			pr.Delete("/carreras/{id}", s.carreras.Delete)

			// Softwares
			pr.Post("/softwares", s.softwares.Create)
			pr.Get("/softwares/{id}", s.softwares.Read)
			pr.Put("/softwares/{id}", s.softwares.Update)
			pr.Delete("/softwares/{id}", s.softwares.Delete)

			// Carrera - Softwares (acciones)
			pr.Post("/carreras/{id}/softwares/{software_id}", s.carreraSoftware.Create)
			pr.Delete("/carreras/{id}/softwares/{software_id}", s.carreraSoftware.Delete)

			// auth/me protegido
			pr.Get("/auth/me", s.auth.Me)
		})

		// Componentes PUBLICAS y PROTEGIDAS
		r.Route("/componentes", func(rr chi.Router) {
			// públicas
			rr.Get("/discos", s.componentes.Discos.List)
			rr.Get("/gabinetes", s.componentes.Gabinetes.List)
			rr.Get("/gpus", s.componentes.GPUs.List)
			rr.Get("/motherboards", s.componentes.Motherboards.List)
			rr.Get("/procesadores", s.componentes.Procesadores.List)
			rr.Get("/psus", s.componentes.PSUs.List)
			rr.Get("/rams", s.componentes.RAMs.List)

			// protegidas
			rr.Group(func(pr chi.Router) {
				pr.Use(auth.AuthMiddleware)

				// Discos
				pr.Post("/discos", s.componentes.Discos.Create)
				pr.Get("/discos/{id}", s.componentes.Discos.Read)
				pr.Put("/discos/{id}", s.componentes.Discos.Update)
				pr.Delete("/discos/{id}", s.componentes.Discos.Delete)

				// Gabinetes
				pr.Post("/gabinetes", s.componentes.Gabinetes.Create)
				pr.Get("/gabinetes/{id}", s.componentes.Gabinetes.Read)
				pr.Put("/gabinetes/{id}", s.componentes.Gabinetes.Update)
				pr.Delete("/gabinetes/{id}", s.componentes.Gabinetes.Delete)

				// GPUs
				pr.Post("/gpus", s.componentes.GPUs.Create)
				pr.Get("/gpus/{id}", s.componentes.GPUs.Read)
				pr.Put("/gpus/{id}", s.componentes.GPUs.Update)
				pr.Delete("/gpus/{id}", s.componentes.GPUs.Delete)

				// Motherboards
				pr.Post("/motherboards", s.componentes.Motherboards.Create)
				pr.Get("/motherboards/{id}", s.componentes.Motherboards.Read)
				pr.Put("/motherboards/{id}", s.componentes.Motherboards.Update)
				pr.Delete("/motherboards/{id}", s.componentes.Motherboards.Delete)

				// Procesadores
				pr.Post("/procesadores", s.componentes.Procesadores.Create)
				pr.Get("/procesadores/{id}", s.componentes.Procesadores.Read)
				pr.Put("/procesadores/{id}", s.componentes.Procesadores.Update)
				pr.Delete("/procesadores/{id}", s.componentes.Procesadores.Delete)

				// PSUs
				pr.Post("/psus", s.componentes.PSUs.Create)
				pr.Get("/psus/{id}", s.componentes.PSUs.Read)
				pr.Put("/psus/{id}", s.componentes.PSUs.Update)
				pr.Delete("/psus/{id}", s.componentes.PSUs.Delete)

				// RAMs
				pr.Post("/rams", s.componentes.RAMs.Create)
				pr.Get("/rams/{id}", s.componentes.RAMs.Read)
				pr.Put("/rams/{id}", s.componentes.RAMs.Update)
				pr.Delete("/rams/{id}", s.componentes.RAMs.Delete)
			})
		})
	})
}

func (s *Server) getHomepage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Homepage requested")
}
