package model

type DTO struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Fabricante  string `json:"fabricante"`
	FormFactor  string `json:"form_factor"`
	MaxLargoGPU int    `json:"max_largo_gpu"`
	MSRP        int    `json:"msrp"`
	ImageURL    string `json:"image_url"`
	Orden       int    `json:"orden"`
	Habilitado  bool   `json:"habilitado"`
}

type Gabinete struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nombre      string `gorm:"type:varchar(100);not null"`
	Fabricante  string `gorm:"type:varchar(20);not null"`
	FormFactor  string `gorm:"type:varchar(10);not null"`
	MaxLargoGPU int    `gorm:"default:0"`
	MSRP        int    `gorm:"default:0"`
	ImageURL    string `gorm:"type:text"`
	Orden       int    `gorm:"default:0"`
	Habilitado  bool   `gorm:"type:bool;default:false"`
}

type Gabinetes []*Gabinete

func (g *Gabinete) ToDto() *DTO {
	return &DTO{
		ID:          g.ID,
		Nombre:      g.Nombre,
		Fabricante:  g.Fabricante,
		FormFactor:  g.FormFactor,
		MaxLargoGPU: g.MaxLargoGPU,
		MSRP:        g.MSRP,
		ImageURL:    g.ImageURL,
		Orden:       g.Orden,
		Habilitado:  g.Habilitado,
	}
}

func (gs Gabinetes) ToDto() []*DTO {
	dtos := make([]*DTO, len(gs))
	for i, v := range gs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (d *DTO) ToModel() *Gabinete {
	return &Gabinete{
		ID:          d.ID,
		Nombre:      d.Nombre,
		Fabricante:  d.Fabricante,
		FormFactor:  d.FormFactor,
		MaxLargoGPU: d.MaxLargoGPU,
		MSRP:        d.MSRP,
		ImageURL:    d.ImageURL,
		Orden:       d.Orden,
		Habilitado:  d.Habilitado,
	}
}
