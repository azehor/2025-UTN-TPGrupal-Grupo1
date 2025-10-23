package model

type DTO struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Fabricante string `json:"fabricante"`
	Capacidad  int    `json:"capacidad"`
	MSRP       string `json:"msrp"`
	ImageURL   string `json:"image_url"`
	Orden      int    `json:"orden"`
	Habilitado bool   `json:"habilitado"`
}

type PSU struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nombre     string `gorm:"type:varchar(100);not null"`
	Fabricante string `gorm:"type:varchar(20);not null"`
	Capacidad  int    `gorm:"default:0"`
	MSRP       string `gorm:"default:0"`
	ImageURL   string `gorm:"type:text"`
	Orden      int    `gorm:"default:0"`
	Habilitado bool   `gorm:"type:bool;default:false"`
}

type PSUs []*PSU

func (p *PSU) ToDto() *DTO {
	return &DTO{
		ID:         p.ID,
		Nombre:     p.Nombre,
		Fabricante: p.Fabricante,
		Capacidad:  p.Capacidad,
		MSRP:       p.MSRP,
		ImageURL:   p.ImageURL,
		Orden:      p.Orden,
		Habilitado: p.Habilitado,
	}
}

func (ps PSUs) ToDto() []*DTO {
	dtos := make([]*DTO, len(ps))
	for i, v := range ps {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (dto *DTO) ToModel() *PSU {
	return &PSU{
		ID:         dto.ID,
		Nombre:     dto.Nombre,
		Fabricante: dto.Fabricante,
		Capacidad:  dto.Capacidad,
		MSRP:       dto.MSRP,
		ImageURL:   dto.ImageURL,
		Orden:      dto.Orden,
		Habilitado: dto.Habilitado,
	}
}
