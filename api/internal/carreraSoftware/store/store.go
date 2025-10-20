package store

import (
	"fmt"
	"quepc/api/internal/carreraSoftware/model"
	"regexp"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	db *gorm.DB
}

// Expresion regex para validar formato de uuid
var uuidRegexp = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List() (model.CarreraSoftwares, error) {
	var list model.CarreraSoftwares
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) Create(carrera *model.CarreraSoftware) (*model.CarreraSoftware, error) {
	if err := s.db.Create(carrera).Error; err != nil {
		return nil, err
	}
	return carrera, nil
}

// Retorna la entidad creada o la existente (preloaded) sin error.
func (s *Store) CreateIfNotExists(rel *model.CarreraSoftware) (*model.CarreraSoftware, error) {
	// Validar formato UUID de las llaves foraneas
	if !uuidRegexp.MatchString(rel.CarreraID) || !uuidRegexp.MatchString(rel.SoftwareID) {
		return nil, fmt.Errorf("carrera_id o software_id invalidos")
	}

	res := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "carrera_id"}, {Name: "software_id"}},
		DoNothing: true,
	}).Create(rel)

	if res.Error != nil {
		return nil, res.Error
	}

	// si RowsAffected == 0, la fila ya existe; recuperarla
	if res.RowsAffected == 0 {
		var existing model.CarreraSoftware
		if err := s.db.Preload("Carrera").Preload("Software").First(&existing, "carrera_id = ? AND software_id = ?", rel.CarreraID, rel.SoftwareID).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	// Caso creado: recargar con relaciones
	if err := s.db.Preload("Carrera").Preload("Software").First(rel, "id = ?", rel.ID).Error; err != nil {
		return nil, err
	}
	return rel, nil
}

func (s *Store) Read(id string) (*model.CarreraSoftware, error) {
	var sw model.CarreraSoftware
	// Validar formato UUID
	if !uuidRegexp.MatchString(id) {
		return nil, fmt.Errorf("id invalido")
	}
	if err := s.db.Preload("Carrera").Preload("Software").First(&sw, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("carreraSoftware con id %s no encontrada", id)
		}
		return nil, err
	}
	return &sw, nil
}

func (s *Store) ListByCarrera(carreraID string) (model.CarreraSoftwares, error) {
	var list model.CarreraSoftwares
	if !uuidRegexp.MatchString(carreraID) {
		return nil, fmt.Errorf("id invalido")
	}
	if err := s.db.Preload("Software").Where("carrera_id = ?", carreraID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) ListBySoftware(softwareID string) (model.CarreraSoftwares, error) {
	var list model.CarreraSoftwares
	if !uuidRegexp.MatchString(softwareID) {
		return nil, fmt.Errorf("id invalido")
	}
	if err := s.db.Preload("Carrera").Where("software_id = ?", softwareID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) DeleteByCarreraAndSoftware(carreraID, softwareID string) (int64, error) {
	if !uuidRegexp.MatchString(carreraID) || !uuidRegexp.MatchString(softwareID) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.CarreraSoftware{}, "carrera_id = ? AND software_id = ?", carreraID, softwareID)
	return res.RowsAffected, res.Error
}

func (s *Store) Delete(id string) (int64, error) {
	// Validar formato UUID
	if !uuidRegexp.MatchString(id) {
		return 0, fmt.Errorf("id invalido")
	}
	res := s.db.Delete(&model.CarreraSoftware{}, "id = ?", id)
	return res.RowsAffected, res.Error
}
