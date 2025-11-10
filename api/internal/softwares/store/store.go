package store

import (
	"fmt"
	"quepc/api/internal/softwares/model"
	"quepc/api/utils"
	"strings"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.Softwares, error) {
	var list model.Softwares
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListByTipo lista todos los softwares, o filtra por tipo cuando se pasa un valor
func (s *Store) ListByTipo(tipo string) (model.Softwares, error) {
	var list model.Softwares
	var err error
	if strings.TrimSpace(tipo) == "" {
		err = s.db.Find(&list).Error
	} else {
		err = s.db.Where("tipo = ?", tipo).Find(&list).Error
	}
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(software *model.Software) (*model.Software, error) {
	if err := s.db.Create(software).Error; err != nil {
		return nil, err
	}
	return software, nil
}

func (s *Store) Read(id string) (*model.Software, error) {
	var sw model.Software
	// Validar formato UUID
	if !utils.ValidarUUID(id) {
		return nil, fmt.Errorf("id invalido")
	}
	if err := s.db.First(&sw, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("software con id %s no encontrado", id)
		}
		return nil, err
	}
	return &sw, nil
}

func (s *Store) Update(software *model.Software) (int64, error) {
	// Validar formato UUID
	if !utils.ValidarUUID(software.ID) {
		return 0, fmt.Errorf("id invalido")
	}

	// usar map para incluir ceros y strings vacíos
	values := map[string]interface{}{
		"tipo":             software.Tipo,
		"nombre":           software.Nombre,
		"empresa":          software.Empresa,
		"image_url":        software.ImageURL,
		"orden_grafica":    software.OrdenGrafica,
		"orden_procesador": software.OrdenProcesador,
		"orden_ram":        software.OrdenRam,
	}

	res := s.db.Model(&model.Software{}).Where("id = ?", software.ID).Updates(values)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	// Validar formato UUID
	if !utils.ValidarUUID(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.Software{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
