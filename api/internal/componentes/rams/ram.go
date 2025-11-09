package rams

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"quepc/api/internal/componentes/rams/model"
	"quepc/api/utils"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.RAMs, error)
	Create(*model.RAM) (*model.RAM, error)
	Read(string) (*model.RAM, error)
	Update(*model.RAM) (int64, error)
	Delete(string) (int64, error)
}

type RAMs struct {
	store Store
}

func New(s Store) *RAMs {
	return &RAMs{
		store: s,
	}
}

func (ram *RAMs) List(w http.ResponseWriter, r *http.Request) {
	rs, err := ram.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar rams", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
	}

	if len(rs) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacia de rams", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	for _, rss := range rs {
		rss.ImageURL = utils.CompletarImageURL(r.Host, rss.ImageURL)
	}

	if err := json.NewEncoder(w).Encode(rs.ToDto()); err != nil {
		slog.Error("Error al codificar rams a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (ram *RAMs) Create(w http.ResponseWriter, r *http.Request) {
	var dto model.DTO
	r.ParseMultipartForm(10 << 20)

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
	dto.ImageURL, err = utils.GuardarImagen(file, handle, "rams")
	if err != nil {
		slog.Error("Error al guardar imagen del request en Create", "error", err)
		http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
		return
	}

	rs := dto.ToModel()

	created, err := ram.store.Create(rs)
	if err != nil {
		slog.Error("Error en la base de datos al crear ram", "error", err)
		http.Error(w, `{"error":"No se pudo crear la ram"}`, http.StatusInternalServerError)
		return
	}

	created.ImageURL = utils.CompletarImageURL(r.Host, created.ImageURL)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar ram creada", "error", err)
	}
}

func (ram *RAMs) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	rs, err := ram.store.Read(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"RAM no encontrada"`, http.StatusInternalServerError)
			return
		}
		slog.Error("Error en la base de datos al leer ram", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	rs.ImageURL = utils.CompletarImageURL(r.Host, rs.ImageURL)

	if err := json.NewEncoder(w).Encode(rs.ToDto()); err != nil {
		slog.Error("Error al codificar la ram leido", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (ram *RAMs) Update(w http.ResponseWriter, r *http.Request) {
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
		dto.ImageURL, err = utils.GuardarImagen(file, handle, "rams")
		if err != nil {
			slog.Error("Error al guardar imagen del request en Create", "error", err)
			http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
			return
		}
	}

	dto.ID = id

	rs := dto.ToModel()

	rows, err := ram.store.Update(rs)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar ram", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar la ram"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"RAM no encontrada"}`, http.StatusNotFound)
		return
	}

	updated, err := ram.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(rs.ToDto()); err != nil {
			slog.Error("Error al codificar la respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
			return
		}
	}

	updated.ImageURL = utils.CompletarImageURL(r.Host, updated.ImageURL)

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar ram actualizada", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (ram *RAMs) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	rows, err := ram.store.Delete(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al eliminar ram", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar la ram"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"RAM no encontrada"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
