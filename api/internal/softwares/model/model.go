package model

type DTO struct {
	ID              string `json:"id"`
	Nombre          string `json:"nombre"`
	Tipo            string `json:"tipo"`
	Empresa         string `json:"empresa"`
	ImageURL        string `json:"image_url"`
	OrdenGrafica    int    `json:"orden_grafica"`
	OrdenProcesador int    `json:"orden_procesador"`
	OrdenRam        int    `json:"orden_ram"`
}

type Software struct {
	ID              string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nombre          string `gorm:"type:varchar(100);not null"`
	Tipo            string `gorm:"type:varchar(50)"`
	Empresa         string `gorm:"type:varchar(100)"`
	ImageURL        string `gorm:"type:text"`
	OrdenGrafica    int    `gorm:"default:0"`
	OrdenProcesador int    `gorm:"default:0"`
	OrdenRam        int    `gorm:"default:0"`
}

type Softwares []*Software

func (s *Software) ToDto() *DTO {
	return &DTO{
		ID:              s.ID,
		Nombre:          s.Nombre,
		Tipo:            s.Tipo,
		Empresa:         s.Empresa,
		ImageURL:        s.ImageURL,
		OrdenGrafica:    s.OrdenGrafica,
		OrdenProcesador: s.OrdenProcesador,
		OrdenRam:        s.OrdenRam,
	}
}

func (ss Softwares) ToDto() []*DTO {
	dtos := make([]*DTO, len(ss))
	for i, v := range ss {
		dtos[i] = v.ToDto()
	}

	return dtos
}
