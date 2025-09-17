package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Un Server incluye una referencia a cada modelo de los endpoints a exponer
// Estos modelos deben ser acompañados por una declaracion de la interfaz que cumplen
type Server struct {
}

// Funcion de inicializacion del server, los argumentos seran los modelos a usar
func New() *Server {
	return &Server{}
}

// Aqui se agregan todas las rutas asociandolas con la funcion que resuelve que se envia
func (s *Server) AddRoutes(r *chi.Mux) {
	//ejemplo de uso
	r.Get("/", s.getHomepage)
}

func (s *Server) getHomepage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Homepage requested")
}
