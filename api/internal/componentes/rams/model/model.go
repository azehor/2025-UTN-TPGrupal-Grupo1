package model

type DTO struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Fabricante string `json:"fabricante"`
	Capacidad  int    `json:"capacidad"`
	Generacion string `json:"generacion"`
	MSRP       string `json:"msrp"`
	ImageURL   string `json:"image_url"`
	Orden      int    `json:"orden"`
	Habilitado bool   `json:"habilitado"`
}

type RAM struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nombre     string `gorm:"type:varchar(100);not null"`
	Fabricante string `gorm:"type:varchar(20);not null"`
	Capacidad  int    `gorm:"default:0"`
	Generacion string `gorm:"type:varchar(10)"`
	MSRP       string `gorm:"default:0"`
	ImageURL   string `gorm:"type:text"`
	Orden      int    `gorm:"default:0"`
	Habilitado bool   `gorm:"type:bool;default:false"`
}

type RAMs []*RAM

func (r *RAM) ToDto() *DTO {
	return &DTO{
		ID:         r.ID,
		Nombre:     r.Nombre,
		Fabricante: r.Fabricante,
		Capacidad:  r.Capacidad,
		Generacion: r.Generacion,
		MSRP:       r.MSRP,
		ImageURL:   r.ImageURL,
		Orden:      r.Orden,
		Habilitado: r.Habilitado,
	}
}

func (rs RAMs) ToDto() []*DTO {
	dtos := make([]*DTO, len(rs))
	for i, v := range rs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (dto *DTO) ToModel() *RAM {
	return &RAM{
		ID:         dto.ID,
		Nombre:     dto.Nombre,
		Fabricante: dto.Fabricante,
		Capacidad:  dto.Capacidad,
		Generacion: dto.Generacion,
		MSRP:       dto.MSRP,
		ImageURL:   dto.ImageURL,
		Orden:      dto.Orden,
		Habilitado: dto.Habilitado,
	}
}
