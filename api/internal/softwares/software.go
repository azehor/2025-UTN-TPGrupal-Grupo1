package softwares

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"quepc/api/internal/softwares/model"
)

type Store interface {
	List() (model.Softwares, error)
	Create(*model.Software) (*model.Software, error)
	Read(string) (*model.Software, error)
	Update(*model.Software) (int64, error)
	Delete(string) (int64, error)
}

type Softwares struct {
	store Store
}

func New(s Store) *Softwares {
	return &Softwares{
		store: s,
	}
}

func (s *Softwares) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	softwares, err := s.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar softwares", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	if len(softwares) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacía de softwares", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(softwares.ToDto()); err != nil {
		slog.Error("Error al codificar softwares a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (s *Softwares) Create(w http.ResponseWriter, r *http.Request) {}

func (s *Softwares) Read(w http.ResponseWriter, r *http.Request) {}

func (s *Softwares) Update(w http.ResponseWriter, r *http.Request) {}

func (s *Softwares) Delete(w http.ResponseWriter, r *http.Request) {}
