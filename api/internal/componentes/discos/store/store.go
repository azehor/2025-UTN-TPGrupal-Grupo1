package store

import (
	"fmt"
	"quepc/api/internal/componentes/discos/model"
	"quepc/api/utils"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.Discos, error) {
	var list model.Discos
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(disco *model.Disco) (*model.Disco, error) {
	if err := s.db.Create(disco).Error; err != nil {
		return nil, err
	}
	return disco, nil
}

func (s *Store) Read(id string) (*model.Disco, error) {
	var p model.Disco
	if !utils.ValidarUUID(id) {
		return nil, fmt.Errorf("id invalido")
	}

	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("disco con id %s no encontrado", id)
		}
		return nil, err
	}

	return &p, nil
}

func (s *Store) Update(p *model.Disco) (int64, error) {
	if !utils.ValidarUUID(p.ID) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Model(&model.Disco{}).Where("id = ?", p.ID).Updates(p)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	if !utils.ValidarUUID(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.Disco{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
