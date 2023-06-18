package model

import (
	"gorm.io/gorm"
)

// Struct untuk menyimpan informasi pengguna
type Admin struct {
	gorm.Model
	Email   string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}