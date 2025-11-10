package recomendaciones

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
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
	"strconv"

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
	IDs       []string `json:"ids"`
	MinBudget float64  `json:"minBudget,omitempty"` // opcional
	MaxBudget float64  `json:"maxBudget,omitempty"` // opcional
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

	var reco model.Recomendacion
	var err error
	if req.MinBudget > 0 || req.MaxBudget > 0 {
		reco, err = rec.buildRecommendationWithBudget(req.IDs, req.MinBudget, req.MaxBudget)
	} else {
		reco, err = rec.buildRecommendation(req.IDs)
	}

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

	var minBudget, maxBudget float64
	if v := r.URL.Query().Get("minBudget"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			minBudget = f
		}
	}
	if v := r.URL.Query().Get("maxBudget"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			maxBudget = f
		}
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

	// 👇 usar el builder correcto
	var reco model.Recomendacion
	if minBudget > 0 || maxBudget > 0 {
		reco, err = rec.buildRecommendationWithBudget(ids, minBudget, maxBudget)
	} else {
		reco, err = rec.buildRecommendation(ids)
	}

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
	var maxGrafica, maxProcesador, maxRam, consumoTotal int
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

	var procesadorAniadido *procsModel.Procesador
	var ramAniadida *ramsModel.RAM

	// GPUs
	if maxGrafica > 0 {
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
				consumoTotal += chosen.Consumo
				reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "GPU", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
			}
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
			procesadorAniadido = chosen
			consumoTotal += chosen.Consumo
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
			ramAniadida = chosen
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "RAM", Precio: chosen.MSRP, ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	// promedio de los 3 principales
	var promedio int
	if maxGrafica > 0 {
		promedio = (maxGrafica + maxProcesador + maxRam) / 3
	} else {
		promedio = (maxProcesador + maxRam) / 2
	}

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

	// Motherboards - Buscamos una que sea compatible con el procesador y RAM ya seleccionados
	if mbs, err := rec.mbStore.List(); err == nil && len(mbs) > 0 {
		var chosen *mbModel.Motherboard
		bestDiff := int(^uint(0) >> 1)
		for _, m := range mbs {
			// Verificar que la motherboard sea compatible con los componentes seleccionados
			if procesadorAniadido != nil && m.SocketProcesador != procesadorAniadido.Socket {
				continue
			}
			if ramAniadida != nil && m.SocketRam != ramAniadida.Generacion {
				continue
			}

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

	// PSUs - seleccionar según consumo total + margen de seguridad
	if psus, err := rec.psuStore.List(); err == nil && len(psus) > 0 {
		// Definir consumo requerido con un margen del 25% y mínimo de 100W
		required := consumoTotal
		if required < 100 {
			required = 100
		}
		safety := int(float64(required) * 1.25)

		// Buscar la PSU con la menor capacidad que sea >= safety
		var chosen *psuModel.PSU
		var bestCap = 0
		for _, p := range psus {
			if p.Capacidad >= safety {
				if chosen == nil || p.Capacidad < bestCap {
					chosen = p
					bestCap = p.Capacidad
				}
			}
		}

		// Si no hay ninguna PSU que cumpla, elegimos la más potente disponible
		if chosen == nil {
			for _, p := range psus {
				if chosen == nil || p.Capacidad > bestCap {
					chosen = p
					bestCap = p.Capacidad
				}
			}
		}
		if chosen != nil {
			reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "PSU", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
		}
	}

	return reco, nil
}

// buildRecommendationWithBudget mantiene la esencia de buildRecommendation
// pero elige, dentro de los que cumplen el orden mínimo, el componente
// cuyo MSRP esté más cerca del sub-presupuesto asignado al rubro.
func (rec *Recomendaciones) buildRecommendationWithBudget(ids []string, minBudget, maxBudget float64) (model.Recomendacion, error) {
	// 1) calcular máximos
	var maxGrafica, maxProcesador, maxRam, consumoTotal int
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
	hasGPU := maxGrafica > 0

	// 2) si no viene tope, usamos min como referencia
	if maxBudget <= 0 {
		maxBudget = minBudget
		if maxBudget <= 0 {
			maxBudget = 1e12
		}
	}

	// 3) reparto del presupuesto por componente
	type bucket string
	const (
		bGPU bucket = "GPU"
		bCPU bucket = "CPU"
		bRAM bucket = "RAM"
		bDIS bucket = "DISCO"
		bMB  bucket = "MB"
		bPSU bucket = "PSU"
		bGAB bucket = "GAB"
	)
	alloc := map[bucket]float64{}
	if hasGPU {
		alloc[bGPU] = 0.35
		alloc[bCPU] = 0.30
		alloc[bRAM] = 0.10
		alloc[bDIS] = 0.10
		alloc[bMB] = 0.10
		alloc[bPSU] = 0.03
		alloc[bGAB] = 0.02
	} else {
		alloc[bCPU] = 0.45
		alloc[bRAM] = 0.20
		alloc[bDIS] = 0.15
		alloc[bMB] = 0.12
		alloc[bPSU] = 0.05
		alloc[bGAB] = 0.03
	}

	var reco model.Recomendacion
	var procesadorAniadido *procsModel.Procesador
	var ramAniadida *ramsModel.RAM

	// GPU
	if hasGPU {
		if gpus, err := rec.gpuStore.List(); err == nil && len(gpus) > 0 {
			target := alloc[bGPU] * maxBudget
			if idx, ok := closestIndex(gpus, maxGrafica,
				target,
				func(g *gpusModel.GPU) int { return g.Orden },
				func(g *gpusModel.GPU) float64 { return float64(g.MSRP) },
			); ok {
				g := gpus[idx]
				consumoTotal += g.Consumo
				reco = append(reco, &model.Componente{ID: g.ID, Nombre: g.Nombre, Tipo: "GPU", Precio: float64(g.MSRP), ImageURL: g.ImageURL, Orden: g.Orden})
			}
		}
	}

	// Procesador
	if procs, err := rec.procStore.List(); err == nil && len(procs) > 0 {
		target := alloc[bCPU] * maxBudget
		if idx, ok := closestIndex(procs, maxProcesador,
			target,
			func(p *procsModel.Procesador) int { return p.Orden },
			func(p *procsModel.Procesador) float64 { return float64(p.MSRP) },
		); ok {
			p := procs[idx]
			procesadorAniadido = p
			consumoTotal += p.Consumo
			reco = append(reco, &model.Componente{ID: p.ID, Nombre: p.Nombre, Tipo: "Procesador", Precio: float64(p.MSRP), ImageURL: p.ImageURL, Orden: p.Orden})
		}
	}

	// RAM
	if rams, err := rec.ramStore.List(); err == nil && len(rams) > 0 {
		target := alloc[bRAM] * maxBudget
		if idx, ok := closestIndex(rams, maxRam,
			target,
			func(rm *ramsModel.RAM) int { return rm.Orden },
			func(rm *ramsModel.RAM) float64 { return rm.MSRP },
		); ok {
			rm := rams[idx]
			ramAniadida = rm
			reco = append(reco, &model.Componente{ID: rm.ID, Nombre: rm.Nombre, Tipo: "RAM", Precio: rm.MSRP, ImageURL: rm.ImageURL, Orden: rm.Orden})
		}
	}

	// Calcular promedio de orden en base a los componentes principales
	var promedio int
	if hasGPU {
		promedio = (maxGrafica + maxProcesador + maxRam) / 3
	} else {
		promedio = (maxProcesador + maxRam) / 2
	}

	// Disco
	if discos, err := rec.discoStore.List(); err == nil && len(discos) > 0 {
		target := alloc[bDIS] * maxBudget
		if idx, ok := closestIndex(discos, promedio,
			target,
			func(d *discosModel.Disco) int { return d.Orden },
			func(d *discosModel.Disco) float64 { return float64(d.MSRP) },
		); ok {
			d := discos[idx]
			reco = append(reco, &model.Componente{ID: d.ID, Nombre: d.Nombre, Tipo: "Disco", Precio: float64(d.MSRP), ImageURL: d.ImageURL, Orden: d.Orden})
		}
	}

	// Gabinete
	if gabs, err := rec.gabStore.List(); err == nil && len(gabs) > 0 {
		target := alloc[bGAB] * maxBudget
		if idx, ok := closestIndex(gabs, promedio,
			target,
			func(g *gabModel.Gabinete) int { return g.Orden },
			func(g *gabModel.Gabinete) float64 { return float64(g.MSRP) },
		); ok {
			g := gabs[idx]
			reco = append(reco, &model.Componente{ID: g.ID, Nombre: g.Nombre, Tipo: "Gabinete", Precio: float64(g.MSRP), ImageURL: g.ImageURL, Orden: g.Orden})
		}
	}

	// Motherboard (conpatibilidad con CPU y RAM)
	if mbs, err := rec.mbStore.List(); err == nil && len(mbs) > 0 {
		compatibles := mbs[:0]
		for _, m := range mbs {
			if procesadorAniadido != nil && m.SocketProcesador != procesadorAniadido.Socket {
				continue
			}
			if ramAniadida != nil && m.SocketRam != ramAniadida.Generacion {
				continue
			}
			compatibles = append(compatibles, m)
		}
		if len(compatibles) == 0 {
			compatibles = mbs
		}
		target := alloc[bMB] * maxBudget
		if idx, ok := closestIndex(compatibles, promedio,
			target,
			func(m *mbModel.Motherboard) int { return m.Orden },
			func(m *mbModel.Motherboard) float64 { return float64(m.MSRP) },
		); ok {
			m := compatibles[idx]
			reco = append(reco, &model.Componente{ID: m.ID, Nombre: m.Modelo, Tipo: "Motherboard", Precio: float64(m.MSRP), ImageURL: m.ImageURL, Orden: m.Orden})
		}
	}

	// PSU
	if psus, err := rec.psuStore.List(); err == nil && len(psus) > 0 {
		required := consumoTotal
		if required < 100 {
			required = 100
		}
		safety := int(float64(required) * 1.25)

		target := alloc[bPSU] * maxBudget
		// primero las que cumplen potencia
		cands := psus[:0]
		for _, p := range psus {
			if p.Capacidad >= safety {
				cands = append(cands, p)
			}
		}
		if len(cands) == 0 {
			// si ninguna cumple, elegimos la más potente
			var chosen *psuModel.PSU
			var bestCap = 0
			for _, p := range psus {
				if chosen == nil || p.Capacidad > bestCap {
					chosen, bestCap = p, p.Capacidad
				}
			}
			if chosen != nil {
				reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "PSU", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
			}
		} else {
			// si hay candidatas, elegir la más cercana al target
			if idx, ok := closestIndex(cands, safety,
				target,
				func(p *psuModel.PSU) int { return p.Orden },
				func(p *psuModel.PSU) float64 { return float64(p.MSRP) },
			); ok {
				chosen := cands[idx]
				reco = append(reco, &model.Componente{ID: chosen.ID, Nombre: chosen.Nombre, Tipo: "PSU", Precio: float64(chosen.MSRP), ImageURL: chosen.ImageURL, Orden: chosen.Orden})
			}
		}
	}

	return reco, nil
}

// closestIndex elige el índice del elemento en 'arr' cuyo MSRP esté más cerca de  'target' entre los que tengan orden >= minOrden.
// Si ninguno cumple, devuelve el si el slice está vacío.
func closestIndex[T any](arr []T, minOrden int, target float64, getOrden func(T) int, getMSRP func(T) float64) (int, bool) {
	best := -1
	ok := false
	bestDiff := math.MaxFloat64
	bestOrden := math.MinInt

	for i, v := range arr {
		ord := getOrden(v)
		if ord < minOrden {
			continue
		}
		diff := math.Abs(getMSRP(v) - target)
		if diff < bestDiff || (diff == bestDiff && ord > bestOrden) {
			best = i
			ok = true
			bestDiff = diff
			bestOrden = ord
		}
	}

	if ok {
		return best, true
	}

	// devolver el de mayor orden si ninguno cumplió minOrden
	for i, v := range arr {
		ord := getOrden(v)
		if ord > bestOrden {
			bestOrden = ord
			best = i
			ok = true
		}
	}

	return best, ok
}
