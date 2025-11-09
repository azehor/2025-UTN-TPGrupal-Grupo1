package model

type DTO struct {
	ID               string `json:"id"`
	Modelo           string `json:"modelo"`
	Fabricante       string `json:"fabricante"`
	FormFactor       string `json:"form_factor"`
	SocketProcesador string `json:"socket_procesador"`
	SocketRam        string `json:"socket_ram"`
	MSRP             int    `json:"msrp"`
	ImageURL         string `json:"image_url"`
	Orden            int    `json:"orden"`
	Habilitado       bool   `json:"habilitado"`
}

type Motherboard struct {
	ID               string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Modelo           string `gorm:"type:varchar(100);not null"`
	Fabricante       string `gorm:"type:varchar(20);not null"`
	FormFactor       string `gorm:"type:varchar(10);not null"`
	SocketProcesador string `gorm:"type:varchar(10);not null"`
	SocketRam        string `gorm:"type:varchar(10);not null"`
	MSRP             int    `gorm:"default:0"`
	ImageURL         string `gorm:"type:text"`
	Orden            int    `gorm:"default:0"`
	Habilitado       bool   `gorm:"type:bool;default:false"`
}

type Motherboards []*Motherboard

func (m *Motherboard) ToDto() *DTO {
	return &DTO{
		ID:               m.ID,
		Modelo:           m.Modelo,
		Fabricante:       m.Fabricante,
		FormFactor:       m.FormFactor,
		SocketProcesador: m.SocketProcesador,
		SocketRam:        m.SocketRam,
		MSRP:             m.MSRP,
		ImageURL:         m.ImageURL,
		Orden:            m.Orden,
		Habilitado:       m.Habilitado,
	}
}

func (mbs Motherboards) ToDto() []*DTO {
	dtos := make([]*DTO, len(mbs))
	for i, v := range mbs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (dto *DTO) ToModel() *Motherboard {
	return &Motherboard{
		ID:               dto.ID,
		Modelo:           dto.Modelo,
		Fabricante:       dto.Fabricante,
		FormFactor:       dto.FormFactor,
		SocketProcesador: dto.SocketProcesador,
		SocketRam:        dto.SocketRam,
		MSRP:             dto.MSRP,
		ImageURL:         dto.ImageURL,
		Orden:            dto.Orden,
		Habilitado:       dto.Habilitado,
	}
}
