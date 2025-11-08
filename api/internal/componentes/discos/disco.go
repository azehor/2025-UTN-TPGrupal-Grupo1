package discos

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"quepc/api/internal/componentes/discos/model"
	"quepc/api/utils"

	"github.com/go-chi/chi/v5"
)

type Store interface {
	List() (model.Discos, error)
	Create(*model.Disco) (*model.Disco, error)
	Read(string) (*model.Disco, error)
	Update(*model.Disco) (int64, error)
	Delete(string) (int64, error)
}

type Discos struct {
	store Store
}

func New(s Store) *Discos {
	return &Discos{
		store: s,
	}
}

func (d *Discos) List(w http.ResponseWriter, r *http.Request) {
	ds, err := d.store.List()
	if err != nil {
		slog.Error("Error en la base de datos al listar discos", "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
	}

	if len(ds) == 0 {
		if err := json.NewEncoder(w).Encode([]any{}); err != nil {
			slog.Error("Error al codificar lista vacia de discos", "error", err)
			http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	for _, dss := range ds {
		dss.ImageURL = "http://" + filepath.Join(r.Host, dss.ImageURL)
	}

	if err := json.NewEncoder(w).Encode(ds.ToDto()); err != nil {
		slog.Error("Error al codificar discos a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (d *Discos) Create(w http.ResponseWriter, r *http.Request) {
	var dto model.DTO
	r.ParseMultipartForm(10 << 20) // Limite de tamaño de imagen: 10MB

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
	dto.ImageURL, err = utils.GuardarImagen(file, handle, "discos")
	if err != nil {
		slog.Error("Error al guardar imagen del request en Create", "error", err)
		http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
		return
	}

	pr := &model.Disco{
		Nombre:     dto.Nombre,
		Fabricante: dto.Fabricante,
		Capacidad:  dto.Capacidad,
		Tipo:       dto.Tipo,
		MSRP:       dto.MSRP,
		ImageURL:   dto.ImageURL,
		Orden:      dto.Orden,
		Habilitado: dto.Habilitado,
	}

	created, err := d.store.Create(pr)
	if err != nil {
		slog.Error("Error en la base de datos al crear disco", "error", err)
		http.Error(w, `{"error":"No se pudo crear el disco"}`, http.StatusInternalServerError)
		return
	}

	created.ImageURL = filepath.Join(r.Host, created.ImageURL)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar disco creado", "error", err)
	}
}

func (d *Discos) Read(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	ds, err := d.store.Read(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"Disco no encontrado"`, http.StatusInternalServerError)
			return
		}
		slog.Error("Error en la base de datos al leer disco", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(ds.ToDto()); err != nil {
		slog.Error("Error al codificar el disco leido", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (d *Discos) Update(w http.ResponseWriter, r *http.Request) {
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

	ds := &model.Disco{
		ID:         dto.ID,
		Nombre:     dto.Nombre,
		Fabricante: dto.Fabricante,
		Capacidad:  dto.Capacidad,
		Tipo:       dto.Tipo,
		MSRP:       dto.MSRP,
		ImageURL:   dto.ImageURL,
		Orden:      dto.Orden,
		Habilitado: dto.Habilitado,
	}

	rows, err := d.store.Update(ds)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar disco", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar el disco"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Disco no encontrado"}`, http.StatusNotFound)
		return
	}

	updated, err := d.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(ds.ToDto()); err != nil {
			slog.Error("Error al codificar la respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar disco actualizado", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (d *Discos) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
		return
	}

	rows, err := d.store.Delete(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID invalido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al eliminar disco", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar el disco"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Disco no encontrado"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
