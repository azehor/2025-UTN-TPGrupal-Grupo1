package recomendaciones

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	carrerasoftware "quepc/api/internal/carreraSoftware"
	discosPkg "quepc/api/internal/componentes/discos"
	discosModel "quepc/api/internal/componentes/discos/model"
	gabPkg "quepc/api/internal/componentes/gabinetes"
	gabModel "quepc/api/internal/componentes/gabinetes/model"
	gpusPkg "quepc/api/internal/componentes/gpus"
	gpusModel "quepc/api/internal/componentes/gpus/model"
	mbPkg "quepc/api/internal/componentes/motherboards"
	mbModel "quepc/api/internal/componentes/motherboards/model"
	procsPkg "quepc/api/internal/componentes/procesadores"
	procsModel "quepc/api/internal/componentes/procesadores/model"
	psuPkg "quepc/api/internal/componentes/psus"
	psuModel "quepc/api/internal/componentes/psus/model"
	ramsPkg "quepc/api/internal/componentes/rams"
	ramsModel "quepc/api/internal/componentes/rams/model"
	"quepc/api/internal/recomendaciones/model"
	softwaresModel "quepc/api/internal/softwares/model"
	"quepc/api/utils"

	"github.com/go-chi/chi/v5"
)

// SoftStore representa las operaciones mínimas que necesitamos del store de softwares
type SoftStore interface {
	Read(string) (*softwaresModel.Software, error)
	List() (softwaresModel.Softwares, error)
}

type Recomendaciones struct {
	softwares  SoftStore
	csStore    carrerasoftware.Store
	gpuStore   gpusPkg.Store
	procStore  procsPkg.Store
	ramStore   ramsPkg.Store
	discoStore discosPkg.Store
	gabStore   gabPkg.Store
	mbStore    mbPkg.Store
	psuStore   psuPkg.Store
}

func New(s SoftStore, c carrerasoftware.Store, g gpusPkg.Store, p procsPkg.Store, r ramsPkg.Store, d discosPkg.Store, gb gabPkg.Store, mb mbPkg.Store, ps psuPkg.Store) *Recomendaciones {
	return &Recomendaciones{softwares: s, csStore: c, gpuStore: g, procStore: p, ramStore: r, discoStore: d, gabStore: gb, mbStore: mb, psuStore: ps}
}

// Utilizada para decodificar el json de ids en el body de la request
type SoftwareReq struct {
	IDs []string `json:"ids"`
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (rec *Recomendaciones) RecomendacionSoftware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req SoftwareReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"json inválido"}`, http.StatusBadRequest)
		return
	}
	if len(req.IDs) == 0 || len(req.IDs) > 50 {
		http.Error(w, `{"error":"cantidad de ids inválida"}`, http.StatusBadRequest)
		return
	}
	for _, id := range req.IDs {
		if !utils.ValidarUUID(id) {
			http.Error(w, `{"error":"ids deben ser UUIDs válidos"}`, http.StatusBadRequest)
			return
		}
	}

	reco, err := rec.buildRecommendation(req.IDs)
	if err != nil {
		slog.Error("Error al construir recomendación por software", "error", err)
		http.Error(w, `{"error":"Error interno al generar la recomendación"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(reco.ToDto()); err != nil {
		slog.Error("Error al codificar recomendacion por software a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
	}
}

func (rec *Recomendaciones) RecomendacionCarrera(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	carreraID := chi.URLParam(r, "id")
	if carreraID == "" {
		http.Error(w, `{"error":"id requerido"}`, http.StatusBadRequest)
		return
	}
	if !utils.ValidarUUID(carreraID) {
		http.Error(w, `{"error":"id inválido, debe ser UUID"}`, http.StatusBadRequest)
		return
	}

	rels, err := rec.csStore.ListByCarrera(carreraID)
	if err != nil {
		slog.Error("Error listando softwares por carrera", "carrera_id", carreraID, "error", err)
		http.Error(w, `{"error":"Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}
	if len(rels) == 0 {
		http.Error(w, `{"error":"La carrera no tiene softwares asociados"}`, http.StatusNotFound)
		return
	}

	ids := make([]string, 0, len(rels))
	for _, rel := range rels {
		if utils.ValidarUUID(rel.SoftwareID) {
			ids = append(ids, rel.SoftwareID)
		}
	}
	if len(ids) == 0 {
		http.Error(w, `{"error":"No hay softwares asociados a la carrera"}`, http.StatusNotFound)
		return
	}

	reco, err := rec.buildRecommendation(ids)
	if err != nil {
		slog.Error("Error al construir recomendación por carrera", "carrera_id", carreraID, "error", err)
		http.Error(w, `{"error":"Error interno al generar la recomendación"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(reco.ToDto()); err != nil {
		slog.Error("Error al codificar recomendacion por carrera a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
	}
}

// buildRecommendation calcula la recomendación a partir de una lista de UUIDs de softwares.
func (rec *Recomendaciones) buildRecommendation(ids []string) (model.Recomendacion, error) {
	// calcular máximos
	var maxGrafica, maxProcesador, maxRam int
	first := true

	for _, id := range ids {
		s, err := rec.softwares.Read(id)
		if err != nil {
			return nil, fmt.Errorf("leyendo software %s: %w", id, err)
		}
		if first {
			maxGrafica = s.OrdenGrafica
			maxProcesador = s.OrdenProcesador
			maxRam = s.OrdenRam
			first = false
			continue
		}
		if s.OrdenGrafica > maxGrafica {
			maxGrafica = s.OrdenGrafica
		}
		if s.OrdenProcesador > maxProcesador {
			maxProcesador = s.OrdenProcesador
		}
		if s.OrdenRam > maxRam {
			maxRam = s.OrdenRam
		}
	}

	var reco model.Recomendacion

	// GPUs
	if gpus, err := rec.gpuStore.List(); err == nil && len(gpus) > 0 {
		var chosen *gpusModel.GPU
		for _, g := range gpus {
			if g.Orden >= maxGrafica && (chosen == nil || g.Orden < chosen.Orden) {
				chosen = g
			}
		}
		if chosen == nil {
			for _, g := range gpus {
				if chosen == nil || g.Orden > chosen.Orden {
					chosen = g
				}
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "GPU", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	// Procesadores
	if procs, err := rec.procStore.List(); err == nil && len(procs) > 0 {
		var chosen *procsModel.Procesador
		for _, p := range procs {
			if p.Orden >= maxProcesador && (chosen == nil || p.Orden < chosen.Orden) {
				chosen = p
			}
		}
		if chosen == nil {
			for _, p := range procs {
				if chosen == nil || p.Orden > chosen.Orden {
					chosen = p
				}
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "Procesador", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	// RAMs
	if rams, err := rec.ramStore.List(); err == nil && len(rams) > 0 {
		var chosen *ramsModel.RAM
		for _, rm := range rams {
			if rm.Orden >= maxRam && (chosen == nil || rm.Orden < chosen.Orden) {
				chosen = rm
			}
		}
		if chosen == nil {
			for _, rm := range rams {
				if chosen == nil || rm.Orden > chosen.Orden {
					chosen = rm
				}
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "RAM", Precio: chosen.MSRP, ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	// promedio de los 3 principales
	promedio := (maxGrafica + maxProcesador + maxRam) / 3

	// En base a ese promedio se eligen los demás componentes
	// Discos
	if discos, err := rec.discoStore.List(); err == nil && len(discos) > 0 {
		var chosen *discosModel.Disco
		bestDiff := int(^uint(0) >> 1)
		for _, d := range discos {
			diff := absInt(d.Orden - promedio)
			if diff < bestDiff || (diff == bestDiff && (chosen == nil || d.Orden > chosen.Orden)) {
				chosen = d
				bestDiff = diff
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "Disco", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	// Gabinetes
	if gabs, err := rec.gabStore.List(); err == nil && len(gabs) > 0 {
		var chosen *gabModel.Gabinete
		bestDiff := int(^uint(0) >> 1)
		for _, g := range gabs {
			diff := absInt(g.Orden - promedio)
			if diff < bestDiff || (diff == bestDiff && (chosen == nil || g.Orden > chosen.Orden)) {
				chosen = g
				bestDiff = diff
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "Gabinete", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	// Motherboards
	if mbs, err := rec.mbStore.List(); err == nil && len(mbs) > 0 {
		var chosen *mbModel.Motherboard
		bestDiff := int(^uint(0) >> 1)
		for _, m := range mbs {
			diff := absInt(m.Orden - promedio)
			if diff < bestDiff || (diff == bestDiff && (chosen == nil || m.Orden > chosen.Orden)) {
				chosen = m
				bestDiff = diff
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Modelo, Tipo: "Motherboard", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	// PSUs
	if psus, err := rec.psuStore.List(); err == nil && len(psus) > 0 {
		var chosen *psuModel.PSU
		bestDiff := int(^uint(0) >> 1)
		for _, p := range psus {
			diff := absInt(p.Orden - promedio)
			if diff < bestDiff || (diff == bestDiff && (chosen == nil || p.Orden > chosen.Orden)) {
				chosen = p
				bestDiff = diff
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "PSU", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	return reco, nil
}
