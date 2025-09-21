package model

type DTO struct {
	ID       string `json:"id"`
	Nombre   string `json:"nombre"`
	Empresa  string `json:"empresa"`
	ImageURL string `json:"image_url"`
}

type Software struct {
	ID       string
	Nombre   string
	Empresa  string
	ImageURL string
}

type Softwares []*Software

func (s *Software) ToDto() *DTO {
	return &DTO{
		ID:       s.ID,
		Nombre:   s.Nombre,
		Empresa:  s.Empresa,
		ImageURL: s.ImageURL,
	}
}

func (ss Softwares) ToDto() []*DTO {
	dtos := make([]*DTO, len(ss))
	for i, v := range ss {
		dtos[i] = v.ToDto()
	}

	return dtos
}
