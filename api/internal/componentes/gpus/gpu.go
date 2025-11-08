package gpus

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"quepc/api/internal/componentes/gpus/model"
	"quepc/api/utils"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.GPUs, error)
	Create(*model.GPU) (*model.GPU, error)
	Read(string) (*model.GPU, error)
	Update(*model.GPU) (int64, error)
	Delete(string) (int64, error)
}

type GPUs struct {
	store Store
}

func New(s Store) *GPUs {
	return &GPUs{
		store: s,
	}
}

func (g *GPUs) List(w http.ResponseWriter, r *http.Request) {
	gp, err := g.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar gpus", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
	}

	if len(gp) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacia de gpus", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	for _, gpp := range gp {
		gpp.ImageURL = utils.CompletarImageURL(r.Host, gpp.ImageURL)
	}

	if err := json.NewEncoder(w).Encode(gp.ToDto()); err != nil {
		slog.Error("Error al codificar gpus a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (g *GPUs) Create(w http.ResponseWriter, r *http.Request) {
	var dto model.DTO
	r.ParseMultipartForm(10 << 20) //Limite de payload 10MB

	if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Create", "error", err)
		http.Error(w, `{"error":"JSON invalido o vacio"}`, http.StatusBadRequest)
		return
	}

	dto.ID = ""

	if strings.TrimSpace(dto.Nombre) == "" {
		http.Error(w, `{"error":"El campo 'nombre' es requerido"}`, http.StatusBadRequest)
		return
	}

	//Parseo y guardado de Imagen
	file, handle, err := r.FormFile("imagen")
	if err != nil {
		slog.Error("Error al leer imagen del request en Create", "error", err)
		http.Error(w, `{"error":"Campo 'imagen' invalido o vacio"}`, http.StatusBadRequest)
		return
	}
	dto.ImageURL, err = utils.GuardarImagen(file, handle, "gpus")
	if err != nil {
		slog.Error("Error al guardar imagen del request en Create", "error", err)
		http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
		return
	}

	gp := dto.ToModel()

	created, err := g.store.Create(gp)
	if err != nil {
		slog.Error("Error en la base de datos al crear gpu", "error", err)
		http.Error(w, `{"error":"No se pudo crear el gpu"}`, http.StatusInternalServerError)
		return
	}

	created.ImageURL = utils.CompletarImageURL(r.Host, created.ImageURL)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar gpu creado", "error", err)
	}
}

func (g *GPUs) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	gp, err := g.store.Read(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"GPU no encontrada"`, http.StatusInternalServerError)
			return
		}
		slog.Error("Error en la base de datos al leer gpu", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	gp.ImageURL = utils.CompletarImageURL(r.Host, gp.ImageURL)
	if err := json.NewEncoder(w).Encode(gp.ToDto()); err != nil {
		slog.Error("Error al codificar la gpu leida", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (g *GPUs) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}
	r.ParseMultipartForm(10 << 20)

	var dto model.DTO
	if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Update", "error", err)
		http.Error(w, `{"error":"JSON invalido o vacio"}`, http.StatusBadRequest)
		return
	}

	//Parseo y guardado de Imagen
	file, handle, err := r.FormFile("imagen")
	if err == nil {
		dto.ImageURL, err = utils.GuardarImagen(file, handle, "gpus")
		if err != nil {
			slog.Error("Error al guardar imagen del request en Create", "error", err)
			http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
			return
		}
	}

	dto.ID = id

	gp := dto.ToModel()

	rows, err := g.store.Update(gp)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar gpu", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar la gpu"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"GPU no encontrada"}`, http.StatusNotFound)
		return
	}

	updated, err := g.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(gp.ToDto()); err != nil {
			slog.Error("Error al codificar la respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar gpu actualizada", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (g *GPUs) Delete(w http.ResponseWriter, r *http.Request) {
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
		slog.Error("Error en la base de datos al eliminar gpu", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar la gpu"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"GPU no encontrada"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
