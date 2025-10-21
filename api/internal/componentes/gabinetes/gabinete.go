package gabinetes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"quepc/api/internal/componentes/gabinetes/model"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.Gabinetes, error)
	Create(*model.Gabinete) (*model.Gabinete, error)
	Read(string) (*model.Gabinete, error)
	Update(*model.Gabinete) (int64, error)
	Delete(string) (int64, error)
}

type Gabinetes struct {
	store Store
}

func New(s Store) *Gabinetes {
	return &Gabinetes{
		store: s,
	}
}

func (g *Gabinetes) List(w http.ResponseWriter, r *http.Request) {
	gs, err := g.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar gabinetes", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
	}

	if len(gs) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacia de gabinetes", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(gs.ToDto()); err != nil {
		slog.Error("Error al codificar gabinetes a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (g *Gabinetes) Create(w http.ResponseWriter, r *http.Request) {
	var dto model.DTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Create", "error", err)
		http.Error(w, `{"error":"JSON invalido o vacio"}`, http.StatusBadRequest)
		return
	}

	dto.ID = ""

	if strings.TrimSpace(dto.Nombre) == "" {
		http.Error(w, `{"error":"El campo 'nombre' es requerido"}`, http.StatusBadRequest)
		return
	}

	gs := dto.ToModel()

	created, err := g.store.Create(gs)
	if err != nil {
		slog.Error("Error en la base de datos al crear gabinete", "error", err)
		http.Error(w, `{"error":"No se pudo crear el gabinete"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar gabinete creado", "error", err)
	}
}

func (g *Gabinetes) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	gs, err := g.store.Read(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"Gabinete no encontrado"`, http.StatusInternalServerError)
			return
		}
		slog.Error("Error en la base de datos al leer gabinete", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(gs.ToDto()); err != nil {
		slog.Error("Error al codificar el gabinete leido", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (g *Gabinetes) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	var dto model.DTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Update", "error", err)
		http.Error(w, `{"error":"JSON invalido o vacio"}`, http.StatusBadRequest)
		return
	}

	dto.ID = id

	gs := dto.ToModel()

	rows, err := g.store.Update(gs)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar gabinete", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar el gabinete"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Gabinete no encontrado"}`, http.StatusNotFound)
		return
	}

	updated, err := g.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(gs.ToDto()); err != nil {
			slog.Error("Error al codificar la respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar gabinete actualizado", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (g *Gabinetes) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	rows, err := g.store.Delete(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al eliminar gabinete", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar el gabinete"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Gabinete no encontrado"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
