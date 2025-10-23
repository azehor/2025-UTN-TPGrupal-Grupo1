package store

import (
	"fmt"
	"quepc/api/internal/componentes/rams/model"
	"quepc/api/utils"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.RAMs, error) {
	var list model.RAMs
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(ram *model.RAM) (*model.RAM, error) {
	if err := s.db.Create(ram).Error; err != nil {
		return nil, err
	}
	return ram, nil
}

func (s *Store) Read(id string) (*model.RAM, error) {
	var p model.RAM
	if !utils.ValidarUUID(id) {
		return nil, fmt.Errorf("id invalido")
	}

	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ram con id %s no encontrada", id)
		}
		return nil, err
	}

	return &p, nil
}

func (s *Store) Update(p *model.RAM) (int64, error) {
	if !utils.ValidarUUID(p.ID) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Model(&model.RAM{}).Where("id = ?", p.ID).Updates(p)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	if !utils.ValidarUUID(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.RAM{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
