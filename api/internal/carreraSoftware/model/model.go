package model

import (
	carrerasModel "quepc/api/internal/carreras/model"
	softwaresModel "quepc/api/internal/softwares/model"
)

type DTO struct {
	ID         string `json:"id"`
	CarreraID  string `json:"carrera_id"`
	SoftwareID string `json:"software_id"`
	Carrera    *carrerasModel.DTO   `json:"carrera,omitempty"`
	Software   *softwaresModel.DTO  `json:"software,omitempty"`
}

type CarreraSoftware struct {
	ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	// Se configura para que solo pueda haber una relacion de la misma carrera con el mismo software
	CarreraID  string `gorm:"type:uuid;not null;uniqueIndex:idx_carrera_software" json:"carrera_id"`
	SoftwareID string `gorm:"type:uuid;not null;uniqueIndex:idx_carrera_software" json:"software_id"`

	// Relaciones opcionales para consultas con Preload
	Carrera  *carrerasModel.Carrera   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:CarreraID;references:ID" json:"carrera,omitempty"`
	Software *softwaresModel.Software `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SoftwareID;references:ID" json:"software,omitempty"`
}

type CarreraSoftwares []*CarreraSoftware

func (cs *CarreraSoftware) ToDto() *DTO {
	dto := &DTO{
		ID:         cs.ID,
		CarreraID:  cs.CarreraID,
		SoftwareID: cs.SoftwareID,
	}
	if cs.Carrera != nil {
		dto.Carrera = cs.Carrera.ToDto()
	}
	if cs.Software != nil {
		dto.Software = cs.Software.ToDto()
	}
	return dto
}

func (css CarreraSoftwares) ToDto() []*DTO {
	dtos := make([]*DTO, len(css))
	for i, v := range css {
		dtos[i] = v.ToDto()
	}
	return dtos
}
