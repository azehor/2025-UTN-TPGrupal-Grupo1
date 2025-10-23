package model

type DTO struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Fabricante string `json:"fabricante"`
	Nucleos    int    `json:"nucleos"`
	Threads    int    `json:"threads"`
	Generacion string `json:"generacion"`
	Consumo    int    `json:"consumo"`
	MSRP       string `json:"msrp"`
	ImageURL   string `json:"image_url"`
	Orden      int    `json:"orden"`
	Habilitado bool   `json:"habilitado"`
}

type Procesador struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nombre     string `gorm:"type:varchar(100);not null"`
	Fabricante string `gorm:"type:varchar(20);not null"`
	Nucleos    int    `gorm:"default:1"`
	Threads    int    `gorm:"default:1"`
	Generacion string `gorm:"type:varchar(10)"`
	Consumo    int    `gorm:"default:0"`
	MSRP       string `gorm:"default:0"`
	ImageURL   string `gorm:"type:text"`
	Orden      int    `gorm:"default:0"`
	Habilitado bool   `gorm:"type:bool;default:false"`
}

type Procesadores []*Procesador

func (p *Procesador) ToDto() *DTO {
	return &DTO{
		ID:         p.ID,
		Nombre:     p.Nombre,
		Fabricante: p.Fabricante,
		Nucleos:    p.Nucleos,
		Threads:    p.Threads,
		Generacion: p.Generacion,
		Consumo:    p.Consumo,
		MSRP:       p.MSRP,
		ImageURL:   p.ImageURL,
		Orden:      p.Orden,
		Habilitado: p.Habilitado,
	}
}

func (ps Procesadores) ToDto() []*DTO {
	dtos := make([]*DTO, len(ps))
	for i, v := range ps {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (dto *DTO) ToModel() *Procesador {
	return &Procesador{
		ID:         dto.ID,
		Nombre:     dto.Nombre,
		Fabricante: dto.Fabricante,
		Nucleos:    dto.Nucleos,
		Threads:    dto.Threads,
		Generacion: dto.Generacion,
		Consumo:    dto.Consumo,
		MSRP:       dto.MSRP,
		ImageURL:   dto.ImageURL,
		Orden:      dto.Orden,
		Habilitado: dto.Habilitado,
	}
}
