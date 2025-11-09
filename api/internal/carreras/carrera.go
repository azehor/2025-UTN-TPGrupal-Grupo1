package carreras

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"quepc/api/internal/carreras/model"
	"quepc/api/utils"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.Carreras, error)
	Create(*model.Carrera) (*model.Carrera, error)
	Read(string) (*model.Carrera, error)
	Update(*model.Carrera) (int64, error)
	Delete(string) (int64, error)
}

type Carreras struct {
	store Store
}

func New(s Store) *Carreras {
	return &Carreras{
		store: s,
	}
}

func (c *Carreras) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Traer listado de carreras de la base de datos
	carreras, err := c.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar carreras", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	if len(carreras) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacía de carreras", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	for _, cr := range carreras {
		cr.ImageURL = utils.CompletarImageURL(r.Host, cr.ImageURL)
	}

	if err := json.NewEncoder(w).Encode(carreras.ToDto()); err != nil {
		slog.Error("Error al codificar carreras a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (c *Carreras) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dto model.DTO
	r.ParseMultipartForm(10 << 20)

	// Decodificar el body para obtener el dto
	if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Create", "error", err)
		http.Error(w, `{"error":"JSON inválido o vacío"}`, http.StatusBadRequest)
		return
	}

	// Al ser create el id es generado por la base de datos
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
	dto.ImageURL, err = utils.GuardarImagen(file, handle, "carreras")
	if err != nil {
		slog.Error("Error al guardar imagen del request en Create", "error", err)
		http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
		return
	}

	// Convertir el dto a carrera
	car := &model.Carrera{
		Nombre:   dto.Nombre,
		ImageURL: dto.ImageURL,
	}

	created, err := c.store.Create(car)

	if err != nil {
		slog.Error("Error en la base de datos al crear carrera", "error", err)
		http.Error(w, `{"error":"No se pudo crear la carrera"}`, http.StatusInternalServerError)
		return
	}

	created.ImageURL = utils.CompletarImageURL(r.Host, created.ImageURL)

	w.WriteHeader(http.StatusCreated) // 201

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar carrera creada", "error", err)
	}
}

func (s *Carreras) Read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el id
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	// Buscar la carrera en la base de datos
	sw, err := s.store.Read(id)

	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"Carrera no encontrada"}`, http.StatusNotFound)
			return
		}
		slog.Error("Error en la base de datos al leer carrera", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	sw.ImageURL = utils.CompletarImageURL(r.Host, sw.ImageURL)

	if err := json.NewEncoder(w).Encode(sw.ToDto()); err != nil {
		slog.Error("Error al codificar carrera leída", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (c *Carreras) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el id
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}
	r.ParseMultipartForm(10 << 20)

	var dto model.DTO
	if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Update", "error", err)
		http.Error(w, `{"error":"JSON inválido o vacío"}`, http.StatusBadRequest)
		return
	}

	//Parseo y guardado de Imagen
	file, handle, err := r.FormFile("imagen")
	if err == nil {
		dto.ImageURL, err = utils.GuardarImagen(file, handle, "carreras")
		if err != nil {
			slog.Error("Error al guardar imagen del request en Create", "error", err)
			http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
			return
		}
	}

	dto.ID = id

	// Convertir dto a carrera
	car := &model.Carrera{
		ID:       dto.ID,
		Nombre:   dto.Nombre,
		ImageURL: dto.ImageURL,
	}

	// Realizar el update en la base de datos
	rows, err := c.store.Update(car)

	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar carrera", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar la carrera"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Carrera no encontrada"}`, http.StatusNotFound)
		return
	}

	// Devolver la carrera actualizada
	updated, err := c.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(car.ToDto()); err != nil {
			slog.Error("Error al codificar respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	updated.ImageURL = utils.CompletarImageURL(r.Host, updated.ImageURL)

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar carrera actualizada", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (c *Carreras) Delete(w http.ResponseWriter, r *http.Request) {
	// Obtener id
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	// Realizar la eliminacion en la base de datos
	rows, err := c.store.Delete(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al eliminar carrera", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar la carrera"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Carrera no encontrada "}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
