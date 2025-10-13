package softwares

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"quepc/api/internal/softwares/model"
	"strings"

	"github.com/go-chi/chi/v5"
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

	// Traer listado de software de la base de datos
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

func (s *Softwares) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dto model.DTO

	// Decodificar el body para obtener el dto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
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

	// Convertir el dto a software
	sw := &model.Software{
		Nombre:          dto.Nombre,
		Tipo:            dto.Tipo,
		Empresa:         dto.Empresa,
		ImageURL:        dto.ImageURL,
		OrdenGrafica:    dto.OrdenGrafica,
		OrdenProcesador: dto.OrdenProcesador,
		OrdenRam:        dto.OrdenRam,
	}
	// Crear nueva entidad en base de datos
	created, err := s.store.Create(sw)
	if err != nil {
		slog.Error("Error en la base de datos al crear software", "error", err)
		http.Error(w, `{"error":"No se pudo crear el software"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) // Codigo 201

	if err := json.NewEncoder(w).Encode(created.ToDto()); err != nil {
		slog.Error("Error al codificar software creado", "error", err)
	}

}

func (s *Softwares) Read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el id
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	// Buscar el software en la base de datos
	sw, err := s.store.Read(id)

	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(lower, "no encontrado") {
			http.Error(w, `{"error":"Software no encontrado"}`, http.StatusNotFound)
			return
		}
		slog.Error("Error en la base de datos al leer software", "id", id, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(sw.ToDto()); err != nil {
		slog.Error("Error al codificar software leído", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (s *Softwares) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el id
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	var dto model.DTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		slog.Error("Error al decodificar body en Update", "error", err)
		http.Error(w, `{"error":"JSON inválido o vacío"}`, http.StatusBadRequest)
		return
	}

	dto.ID = id

	// Convertir dto a software
	sw := &model.Software{
		ID:              dto.ID,
		Tipo:            dto.Tipo,
		Nombre:          dto.Nombre,
		Empresa:         dto.Empresa,
		ImageURL:        dto.ImageURL,
		OrdenGrafica:    dto.OrdenGrafica,
		OrdenProcesador: dto.OrdenProcesador,
		OrdenRam:        dto.OrdenRam,
	}

	// Realizar el update en la base de datos
	rows, err := s.store.Update(sw)

	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al actualizar software", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo actualizar el software"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Software no encontrado"}`, http.StatusNotFound)
		return
	}

	// Devolver el software actualizado
	updated, err := s.store.Read(id)
	if err != nil {
		slog.Warn("Fallo lectura post-actualización, devolviendo payload", "id", id, "error", err)
		if err := json.NewEncoder(w).Encode(sw.ToDto()); err != nil {
			slog.Error("Error al codificar respuesta", "error", err)
			http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(updated.ToDto()); err != nil {
		slog.Error("Error al codificar software actualizado", "error", err)
		http.Error(w, `{"error":"Error al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

func (s *Softwares) Delete(w http.ResponseWriter, r *http.Request) {
	// Obtener id
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	// Realizar la eliminacion en la base de datos
	rows, err := s.store.Delete(id)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "id invalido") {
			http.Error(w, `{"error":"ID inválido"}`, http.StatusBadRequest)
			return
		}
		slog.Error("Error en la base de datos al eliminar software", "id", id, "error", err)
		http.Error(w, `{"error":"No se pudo eliminar el software"}`, http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, `{"error":"Software no encontrado"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
