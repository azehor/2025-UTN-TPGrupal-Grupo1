package psus

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"quepc/api/internal/componentes/psus/model"
	"quepc/api/utils"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.PSUs, error)
	Create(*model.PSU) (*model.PSU, error)
	Read(string) (*model.PSU, error)
	Update(*model.PSU) (int64, error)
	Delete(string) (int64, error)
}

type PSUs struct {
	store Store
}

func New(s Store) *PSUs {
	return &PSUs{
		store: s,
	}
}

func (p *PSUs) List(w http.ResponseWriter, r *http.Request) {
	ps, err := p.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar psus", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
	}

	if len(ps) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacia de psus", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	for _, pss := range ps {
		pss.ImageURL = utils.CompletarImageURL(r.Host, pss.ImageURL)
	}

	if err := json.NewEncoder(w).Encode(ps.ToDto()); err != nil {
		slog.Error("Error al codificar psus a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (p *PSUs) Create(w http.ResponseWriter, r *http.Request) {
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
	dto.ImageURL, err = utils.GuardarImagen(file, handle, "psus")
	if err != nil {
		slog.Error("Error al guardar imagen del request en Create", "error", err)
		http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
		return
	}

	ps := dto.ToModel()

	created, err := p.store.Create(ps)
	if err != nil {
		slog.Error("Error en la base de datos al crear psu", "error", err)
		http.Error(w, `{"error":"No se pudo crear la psu"}`, http.StatusInternalServerError)
		return
	}

	ps.ImageURL = utils.CompletarImageURL(r.Host, ps.ImageURL)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar psu creado", "error", err)
	}
}

func (p *PSUs) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	ps, err := p.store.Read(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"PSU no encontrada"`, http.StatusInternalServerError)
			return
		}
		slog.Error("Error en la base de datos al leer psu", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	ps.ImageURL = utils.CompletarImageURL(r.Host, ps.ImageURL)

	if err := json.NewEncoder(w).Encode(ps.ToDto()); err != nil {
		slog.Error("Error al codificar la psu leido", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (p *PSUs) Update(w http.ResponseWriter, r *http.Request) {
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

	dto.ID = id

	ps := dto.ToModel()

	//Parseo y guardado de Imagen
	file, handle, err := r.FormFile("imagen")
	if err == nil {
		dto.ImageURL, err = utils.GuardarImagen(file, handle, "psus")
		if err != nil {
			slog.Error("Error al guardar imagen del request en Create", "error", err)
			http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
			return
		}
	}

	rows, err := p.store.Update(ps)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar psu", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar la psu"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"PSU no encontrada"}`, http.StatusNotFound)
		return
	}

	updated, err := p.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(ps.ToDto()); err != nil {
			slog.Error("Error al codificar la respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
			return
		}
	}

	updated.ImageURL = utils.CompletarImageURL(r.Host, updated.ImageURL)

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar psu actualizada", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (p *PSUs) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	rows, err := p.store.Delete(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al eliminar psu", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar la psu"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"PSU no encontrada"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
