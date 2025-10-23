package store

import (
	"fmt"
	"quepc/api/internal/componentes/gabinetes/model"
	"quepc/api/utils"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.Gabinetes, error) {
	var list model.Gabinetes
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(gabinete *model.Gabinete) (*model.Gabinete, error) {
	if err := s.db.Create(gabinete).Error; err != nil {
		return nil, err
	}
	return gabinete, nil
}

func (s *Store) Read(id string) (*model.Gabinete, error) {
	var p model.Gabinete
	if !utils.ValidarUUID(id) {
		return nil, fmt.Errorf("id invalido")
	}

	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("gabinete con id %s no encontrado", id)
		}
		return nil, err
	}

	return &p, nil
}

func (s *Store) Update(p *model.Gabinete) (int64, error) {
	if !utils.ValidarUUID(p.ID) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Model(&model.Gabinete{}).Where("id = ?", p.ID).Updates(p)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	if !utils.ValidarUUID(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.Gabinete{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
