package procesadores

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"quepc/api/internal/componentes/procesadores/model"
	"quepc/api/utils"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.Procesadores, error)
	Create(*model.Procesador) (*model.Procesador, error)
	Read(string) (*model.Procesador, error)
	Update(*model.Procesador) (int64, error)
	Delete(string) (int64, error)
}

type Procesadores struct {
	store Store
}

func New(s Store) *Procesadores {
	return &Procesadores{
		store: s,
	}
}

func (p *Procesadores) List(w http.ResponseWriter, r *http.Request) {
	pr, err := p.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar procesadores", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
	}

	if len(pr) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacia de procesadores", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	for _, prs := range pr {
		prs.ImageURL = utils.CompletarImageURL(r.Host, prs.ImageURL)
	}

	if err := json.NewEncoder(w).Encode(pr.ToDto()); err != nil {
		slog.Error("Error al codificar procesadores a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (p *Procesadores) Create(w http.ResponseWriter, r *http.Request) {
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
	dto.ImageURL, err = utils.GuardarImagen(file, handle, "procesadores")
	if err != nil {
		slog.Error("Error al guardar imagen del request en Create", "error", err)
		http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
		return
	}

	pr := dto.ToModel()

	created, err := p.store.Create(pr)
	if err != nil {
		slog.Error("Error en la base de datos al crear procesador", "error", err)
		http.Error(w, `{"error":"No se pudo crear el procesador"}`, http.StatusInternalServerError)
		return
	}

	pr.ImageURL = utils.CompletarImageURL(r.Host, pr.ImageURL)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar procesador creado", "error", err)
	}
}

func (p *Procesadores) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	pr, err := p.store.Read(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"Procesador no encontrado"`, http.StatusInternalServerError)
			return
		}
		slog.Error("Error en la base de datos al leer procesador", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	pr.ImageURL = utils.CompletarImageURL(r.Host, pr.ImageURL)

	if err := json.NewEncoder(w).Encode(pr.ToDto()); err != nil {
		slog.Error("Error al codificar el procesador leido", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (p *Procesadores) Update(w http.ResponseWriter, r *http.Request) {
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
		dto.ImageURL, err = utils.GuardarImagen(file, handle, "procesadores")
		if err != nil {
			slog.Error("Error al guardar imagen del request en Create", "error", err)
			http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
			return
		}
	}

	dto.ID = id

	pr := dto.ToModel()

	rows, err := p.store.Update(pr)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar procesador", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar el procesador"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Procesador no encontrado"}`, http.StatusNotFound)
		return
	}

	updated, err := p.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(pr.ToDto()); err != nil {
			slog.Error("Error al codificar la respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
			return
		}
	}

	updated.ImageURL = utils.CompletarImageURL(r.Host, updated.ImageURL)

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar procesador actualizado", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (p *Procesadores) Delete(w http.ResponseWriter, r *http.Request) {
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
		slog.Error("Error en la base de datos al eliminar procesador", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar el procesador"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Procesador no encontrado"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
