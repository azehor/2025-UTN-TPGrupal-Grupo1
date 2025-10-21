package model

type DTO struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Fabricante string `json:"fabricante"`
	Capacidad  int    `json:"capacidad"`
	Tipo       string `json:"tipo"`
	MSRP       string `json:"msrp"`
	ImageURL   string `json:"image_url"`
	Orden      int    `json:"orden"`
	Habilitado bool   `json:"habilitado"`
}

type Disco struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid"`
	Nombre     string `gorm:"type:varchar(100);not null"`
	Fabricante string `gorm:"type:varchar(20);not null"`
	Capacidad  int    `gorm:"default:0"`
	Tipo       string `gorm:"varchar(10)"`
	MSRP       string `gorm:"default:0"`
	ImageURL   string `gorm:"type:text"`
	Orden      int    `gorm:"default:0"`
	Habilitado bool   `gorm:"type:bool;default:false"`
}

type Discos []*Disco

func (d *Disco) ToDto() *DTO {
	return &DTO{
		ID:         d.ID,
		Nombre:     d.Nombre,
		Fabricante: d.Fabricante,
		Capacidad:  d.Capacidad,
		Tipo:       d.Tipo,
		MSRP:       d.MSRP,
		ImageURL:   d.ImageURL,
		Orden:      d.Orden,
		Habilitado: d.Habilitado,
	}
}

func (ds Discos) ToDto() []*DTO {
	dtos := make([]*DTO, len(ds))
	for i, v := range ds {
		dtos[i] = v.ToDto()
	}

	return dtos
}
