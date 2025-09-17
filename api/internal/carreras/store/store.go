package store

import (
	"quepc/api/internal/carreras/model"
)

type Store struct {
	//GORM db
}

func New() *Store {
	return &Store{}
}

func (s *Store) List() (model.Carreras, error) {
	//carreras := make([]*model.Carrera, 0)
	carreras := hardcodedData()
	//TODO:Recuperar de la db las carreras
	return carreras, nil
}

func (s *Store) Create(carrera *model.Carrera) (*model.Carrera, error) {
	//TODO: Insertar en la db
	return carrera, nil
}

func (s *Store) Read(id string) (*model.Carrera, error) {
	carrera := &model.Carrera{}
	//TODO: Recuperar carrera con id = id

	return carrera, nil
}

func (s *Store) Update(carrera *model.Carrera) (int64, error) {
	//TODO: Modificar la carrera de id = carrera.ID

	return 0, nil //(result.RowsAffected, result.Error)
}

func (s *Store) Delete(id string) (int64, error) {
	//TODO: Eliminar la carrera de id = id

	return 0, nil //(result.RowsAffected, result.Error)
}

func hardcodedData() []*model.Carrera {
	c := make([]*model.Carrera, 0)
	c = append(c, &model.Carrera{
		ID:       "1",
		Nombre:   "Ing. en Sistemas de Informacion",
		ImageURL: "",
	})
	c = append(c, &model.Carrera{
		ID:       "2",
		Nombre:   "Licenciado en Medicina",
		ImageURL: "",
	})
	c = append(c, &model.Carrera{
		ID:       "3",
		Nombre:   "Tecnicatura en Enfermeria",
		ImageURL: "",
	})
	c = append(c, &model.Carrera{
		ID:       "4",
		Nombre:   "Ing. en Quimica",
		ImageURL: "",
	})
	return c
}
