package carreras

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"quepc/api/internal/carreras/model"
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
	carreras, err := c.store.List()
	if err != nil {
		slog.Error(err.Error()) //TODO: explicitar que es error de acceso a la DB
		return
	}

	if len(carreras) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(carreras.ToDto()); err != nil {
		slog.Error(err.Error()) //TODO: explicitar que es error de JSON encoding
		return
	}
}

func (c *Carreras) Create(w http.ResponseWriter, r *http.Request) {}

func (c *Carreras) Read(w http.ResponseWriter, r *http.Request) {}

func (c *Carreras) Update(w http.ResponseWriter, r *http.Request) {}

func (c *Carreras) Delete(w http.ResponseWriter, r *http.Request) {}
