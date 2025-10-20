package carrerasoftware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"quepc/api/internal/carreraSoftware/model"
	softwaresModel "quepc/api/internal/softwares/model"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	CreateIfNotExists(*model.CarreraSoftware) (*model.CarreraSoftware, error)
	ListByCarrera(string) (model.CarreraSoftwares, error)
	DeleteByCarreraAndSoftware(string, string) (int64, error)
}

type CarreraSoftwares struct {
	store Store
}

func New(s Store) *CarreraSoftwares {
	return &CarreraSoftwares{store: s}
}

func (cs *CarreraSoftwares) ListByCarrera(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	list, err := cs.store.ListByCarrera(id)
	if err != nil {
		slog.Error("Error en la base de datos al listar softwares por carrera", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacía de softwares por carrera", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	// Se devuelve solo el listado de softwares
	softwares := make([]*softwaresModel.DTO, 0, len(list))
	for _, rel := range list {
		if rel.Software != nil {
			softwares = append(softwares, rel.Software.ToDto())
		}
	}

	if err := json.NewEncoder(w).Encode(softwares); err != nil {
		slog.Error("Error al codificar softwares vinculados a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (cs *CarreraSoftwares) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	carreraID := chi.URLParam(r, "id")

	if strings.TrimSpace(carreraID) == "" {
		http.Error(w, `{"error":"ID de carrera inválido"}`, http.StatusBadRequest)
		return
	}

	softwareID := chi.URLParam(r, "software_id")

	if strings.TrimSpace(softwareID) == "" {
		http.Error(w, `{"error":"ID de software inválido"}`, http.StatusBadRequest)
		return
	}

	rel := &model.CarreraSoftware{
		CarreraID:  carreraID,
		SoftwareID: softwareID,
	}

	created, err := cs.store.CreateIfNotExists(rel)
	if err != nil {
		slog.Error("Error en la base de datos al crear relacion carrera-software", "error", err)
		http.Error(w, `{"error":"No se pudo vincular el software"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar relacion carrera-software", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
	}
}

func (cs *CarreraSoftwares) Delete(w http.ResponseWriter, r *http.Request) {
	carreraID := chi.URLParam(r, "id")
	softwareID := chi.URLParam(r, "software_id")
	if strings.TrimSpace(carreraID) == "" || strings.TrimSpace(softwareID) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	rows, err := cs.store.DeleteByCarreraAndSoftware(carreraID, softwareID)
	if err != nil {
		slog.Error("Error en la base de datos al eliminar relacion carrera-software", "error", err)
		http.Error(w, `{"error":"No se pudo eliminar la relacion carrera-software"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Relacion carrera-software no encontrada"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
