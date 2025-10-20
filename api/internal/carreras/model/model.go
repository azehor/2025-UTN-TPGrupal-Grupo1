package model

type DTO struct {
	ID       string `json:"id"`
	Nombre   string `json:"nombre"`
	ImageURL string `json:"image_url"`
}

type Carrera struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nombre   string `gorm:"type:varchar(100);not null"`
	ImageURL string `gorm:"type:text"`
}

type Carreras []*Carrera

func (c *Carrera) ToDto() *DTO {
	return &DTO{
		ID:       c.ID,
		Nombre:   c.Nombre,
		ImageURL: c.ImageURL,
	}
}

func (cs Carreras) ToDto() []*DTO {
	dtos := make([]*DTO, len(cs))
	for i, v := range cs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

//TODO: form-to-model func, si implementamos un panel de administracion
