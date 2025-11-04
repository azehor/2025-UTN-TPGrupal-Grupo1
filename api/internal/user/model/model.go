package model

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	Name      string `json:"name"`
	Role      string `gorm:"type:VARCHAR(20);not null;default:'EDITOR'" json:"role"` // "ADMIN","EDITOR"
	CreatedAt time.Time
	UpdatedAt time.Time
}
