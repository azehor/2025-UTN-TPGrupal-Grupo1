package store

import (
	"fmt"
	"quepc/api/internal/carreras/model"
	"quepc/api/utils"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.Carreras, error) {
	var list model.Carreras
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(carrera *model.Carrera) (*model.Carrera, error) {
	if err := s.db.Create(carrera).Error; err != nil {
		return nil, err
	}
	return carrera, nil
}

func (s *Store) Read(id string) (*model.Carrera, error) {
	var sw model.Carrera
	// Validar formato UUID
	if !utils.ValidarUUID(id) {
		return nil, fmt.Errorf("id invalido")
	}
	if err := s.db.First(&sw, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("carrera con id %s no encontrada", id)
		}
		return nil, err
	}
	return &sw, nil
}

func (s *Store) Update(carrera *model.Carrera) (int64, error) {
	// Validar formato UUID
	if !utils.ValidarUUID(carrera.ID) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Model(&model.Carrera{}).Where("id = ?", carrera.ID).Updates(carrera)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	// Validar formato UUID
	if !utils.ValidarUUID(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.Carrera{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
