package recomendaciones

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"quepc/api/internal/recomendaciones/model"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Recomendaciones struct {
}

func New() *Recomendaciones {
	return &Recomendaciones{}
}

func (rec *Recomendaciones) RecomendacionSoftware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener listado de ids de software del body
	var req SoftwareReq

	// Validaciones
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"json inválido"}`, http.StatusBadRequest)
		return
	}
	// Validaciones de seguridad
	if len(req.IDs) == 0 || len(req.IDs) > 50 { // No mas de 50 ids
		http.Error(w, `{"error":"cantidad de ids inválida"}`, http.StatusBadRequest)
		return
	}
	// Evitar numeros negativos
	for _, id := range req.IDs {
		if id <= 0 {
			http.Error(w, `{"error":"ids deben ser enteros positivos"}`, http.StatusBadRequest)
			return
		}
	}

	fmt.Println("Se esta generando recomendacion del siguiente listado de ids de softwares: ", req)

	recomendacion := hardcodedData()

	if err := json.NewEncoder(w).Encode(recomendacion.ToDto()); err != nil {
		slog.Error("Error al codificar recomendacion por software a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}
}

// Utilizada para decodificar el json de ids en el body de la request
type SoftwareReq struct {
	IDs []int `json:"ids"`
}

func (rec *Recomendaciones) RecomendacionCarrera(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener parametro id de la url
	idStr := chi.URLParam(r, "id")

	// Validaciones
	if idStr == "" {
		http.Error(w, `{"error":"id requerido"}`, http.StatusBadRequest)
		return
	}

	// Convertir a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"id inválido, debe ser numérico"}`, http.StatusBadRequest)
		return
	}

	// Validar que sea positivo
	if id <= 0 {
		http.Error(w, `{"error":"id inválido, debe ser mayor a 0"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("Se esta generando recomendacion de carrera de id: ", id)

	recomendacion := hardcodedData()

	if err := json.NewEncoder(w).Encode(recomendacion.ToDto()); err != nil {
		slog.Error("Error al codificar recomendacion por carrera a JSON", "error", err)
		http.Error(w, `{"error":"Error interno al generar la respuesta"}`, http.StatusInternalServerError)
		return
	}

}

func hardcodedData() model.Recomendacion {
	return model.Recomendacion{
		&model.Componente{
			ID:       "1",
			Nombre:   "AMD Ryzen 5 5600X",
			Tipo:     "Procesador",
			Precio:   "250",
			ImageURL: "https://example.com/ryzen5600x.png",
		},
		&model.Componente{
			ID:       "2",
			Nombre:   "ASUS TUF Gaming B550-PLUS",
			Tipo:     "Placa Madre",
			Precio:   "180",
			ImageURL: "https://example.com/asus-b550.png",
		},
		&model.Componente{
			ID:       "3",
			Nombre:   "Corsair Vengeance LPX 16GB (2x8GB) 3200MHz",
			Tipo:     "Memoria RAM",
			Precio:   "90",
			ImageURL: "https://example.com/corsair-ram.png",
		},
		&model.Componente{
			ID:       "4",
			Nombre:   "NVIDIA GeForce RTX 3060 Ti",
			Tipo:     "Tarjeta Gráfica",
			Precio:   "400",
			ImageURL: "https://example.com/rtx3060ti.png",
		},
		&model.Componente{
			ID:       "5",
			Nombre:   "Samsung 970 EVO Plus 1TB NVMe",
			Tipo:     "Disco SSD",
			Precio:   "120",
			ImageURL: "https://example.com/samsung970.png",
		},
		&model.Componente{
			ID:       "6",
			Nombre:   "Seagate Barracuda 2TB",
			Tipo:     "Disco HDD",
			Precio:   "60",
			ImageURL: "https://example.com/seagate2tb.png",
		},
		&model.Componente{
			ID:       "7",
			Nombre:   "EVGA 650W 80+ Bronze",
			Tipo:     "Fuente de Poder",
			Precio:   "75",
			ImageURL: "https://example.com/evga650w.png",
		},
		&model.Componente{
			ID:       "8",
			Nombre:   "NZXT H510",
			Tipo:     "Gabinete",
			Precio:   "80",
			ImageURL: "https://example.com/nzxt-h510.png",
		},
	}
}
