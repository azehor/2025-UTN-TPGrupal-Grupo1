package motherboards

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"quepc/api/internal/componentes/motherboards/model"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.Motherboards, error)
	Create(*model.Motherboard) (*model.Motherboard, error)
	Read(string) (*model.Motherboard, error)
	Update(*model.Motherboard) (int64, error)
	Delete(string) (int64, error)
}

type Motherboards struct {
	store Store
}

func New(s Store) *Motherboards {
	return &Motherboards{
		store: s,
	}
}

func (m *Motherboards) List(w http.ResponseWriter, r *http.Request) {
	mb, err := m.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar motherboards", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
	}

	if len(mb) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacia de motherboards", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(mb.ToDto()); err != nil {
		slog.Error("Error al codificar motherboards a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (m *Motherboards) Create(w http.ResponseWriter, r *http.Request) {
	var dto model.DTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Create", "error", err)
		http.Error(w, `{"error":"JSON invalido o vacio"}`, http.StatusBadRequest)
		return
	}

	dto.ID = ""

	if strings.TrimSpace(dto.Modelo) == "" {
		http.Error(w, `{"error":"El campo 'modelo' es requerido"}`, http.StatusBadRequest)
		return
	}

	mb := dto.ToModel()

	created, err := m.store.Create(mb)
	if err != nil {
		slog.Error("Error en la base de datos al crear motherboard", "error", err)
		http.Error(w, `{"error":"No se pudo crear la motherboard"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar motherboard creada", "error", err)
	}
}

func (m *Motherboards) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	mb, err := m.store.Read(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"Motherboard no encontrada"`, http.StatusInternalServerError)
			return
		}
		slog.Error("Error en la base de datos al leer motherboard", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(mb.ToDto()); err != nil {
		slog.Error("Error al codificar la motherboard leida", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (m *Motherboards) Update(w http.ResponseWriter, r *http.Request) {
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

	mb := dto.ToModel()

	rows, err := m.store.Update(mb)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar motherboard", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar la motherboard"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Motherboard no encontrada"}`, http.StatusNotFound)
		return
	}

	updated, err := m.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(mb.ToDto()); err != nil {
			slog.Error("Error al codificar la respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar motherboard actualizada", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (m *Motherboards) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	rows, err := m.store.Delete(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al eliminar motherboard", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar la motherboard"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Motherboard no encontrada"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
