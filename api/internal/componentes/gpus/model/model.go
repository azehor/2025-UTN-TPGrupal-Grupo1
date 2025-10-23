package model

type DTO struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Fabricante string `json:"fabricante"`
	Modelo     string `json:"modelo"`
	VRAM       int    `json:"vram"`
	Largo      int    `json:"largo"`
	Consumo    int    `json:"consumo"`
	MSRP       string `json:"msrp"`
	ImageURL   string `json:"image_url"`
	Orden      int    `json:"orden"`
	Habilitado bool   `json:"habilitado"`
}

type GPU struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid"`
	Nombre     string `gorm:"type:varchar(100);not null"`
	Fabricante string `gorm:"type:varchar(20);not null"`
	Modelo     string `gorm:"type:varchar(100);not null"`
	VRAM       int    `gorm:"default:0"`
	Largo      int    `gorm:"default:0"`
	Consumo    int    `gorm:"default:0"`
	MSRP       string `gorm:"default:0"`
	ImageURL   string `gorm:"type:text"`
	Orden      int    `gorm:"default:0"`
	Habilitado bool   `gorm:"type:bool;default:false"`
}

type GPUs []*GPU

func (g *GPU) ToDto() *DTO {
	return &DTO{
		ID:         g.ID,
		Nombre:     g.Nombre,
		Fabricante: g.Fabricante,
		Modelo:     g.Modelo,
		VRAM:       g.VRAM,
		Largo:      g.Largo,
		Consumo:    g.Consumo,
		MSRP:       g.MSRP,
		ImageURL:   g.ImageURL,
		Orden:      g.Orden,
		Habilitado: g.Habilitado,
	}
}

func (gs GPUs) ToDto() []*DTO {
	dtos := make([]*DTO, len(gs))
	for i, v := range gs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (dto DTO) ToModel() *GPU {
	return &GPU{
		ID:         dto.ID,
		Nombre:     dto.Nombre,
		Fabricante: dto.Fabricante,
		Modelo:     dto.Modelo,
		VRAM:       dto.VRAM,
		Largo:      dto.Largo,
		Consumo:    dto.Consumo,
		MSRP:       dto.MSRP,
		ImageURL:   dto.ImageURL,
		Orden:      dto.Orden,
		Habilitado: dto.Habilitado,
	}
}
