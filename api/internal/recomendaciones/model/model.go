package model

type DTO struct {
	ID       string `json:"id"`
	Nombre   string `json:"nombre"`
	Tipo     string `json:"tipo"`
	Precio   string `json:"precio"`
	ImageURL string `json:"image_url"`
}

type Componente struct {
	ID       string
	Nombre   string
	Tipo     string
	Precio   string
	ImageURL string
}

type Recomendacion []*Componente

func (c *Componente) ToDto() *DTO {
	return &DTO{
		ID:       c.ID,
		Nombre:   c.Nombre,
		Tipo:     c.Tipo,
		Precio:   c.Precio,
		ImageURL: c.ImageURL,
	}
}

func (cs Recomendacion) ToDto() []*DTO {
	dtos := make([]*DTO, len(cs))
	for i, v := range cs {
		dtos[i] = v.ToDto()
	}

	return dtos
}
