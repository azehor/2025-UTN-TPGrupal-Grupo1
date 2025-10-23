package store

import (
	"fmt"
	"quepc/api/internal/componentes/gpus/model"
	"quepc/api/utils"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.GPUs, error) {
	var list model.GPUs
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(gpu *model.GPU) (*model.GPU, error) {
	if err := s.db.Create(gpu).Error; err != nil {
		return nil, err
	}
	return gpu, nil
}

func (s *Store) Read(id string) (*model.GPU, error) {
	var p model.GPU
	if !utils.ValidarUUID(id) {
		return nil, fmt.Errorf("id invalido")
	}

	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gpu con id %s no encontrado", id)
		}
		return nil, err
	}

	return &p, nil
}

func (s *Store) Update(g *model.GPU) (int64, error) {
	if !utils.ValidarUUID(g.ID) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Model(&model.GPU{}).Where("id = ?", g.ID).Updates(g)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	if !utils.ValidarUUID(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.GPU{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
