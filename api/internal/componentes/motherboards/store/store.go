package store

import (
	"fmt"
	"quepc/api/internal/componentes/motherboards/model"
	"quepc/api/utils"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.Motherboards, error) {
	var list model.Motherboards
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(motherboard *model.Motherboard) (*model.Motherboard, error) {
	if err := s.db.Create(motherboard).Error; err != nil {
		return nil, err
	}
	return motherboard, nil
}

func (s *Store) Read(id string) (*model.Motherboard, error) {
	var p model.Motherboard
	if !utils.ValidarUUID(id) {
		return nil, fmt.Errorf("id invalido")
	}

	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("motherboard con id %s no encontrado", id)
		}
		return nil, err
	}

	return &p, nil
}

func (s *Store) Update(p *model.Motherboard) (int64, error) {
	if !utils.ValidarUUID(p.ID) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Model(&model.Motherboard{}).Where("id = ?", p.ID).Updates(p)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	if !utils.ValidarUUID(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.Motherboard{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
