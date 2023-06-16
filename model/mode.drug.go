package model

import "gorm.io/gorm"

type Drug struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"nama" form:"nama"`
}
