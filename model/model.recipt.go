package model

import "gorm.io/gorm"

type Recipt struct {
	gorm.Model
	Drugs    []Drug `gorm:"many2many:recipt_drugs" json:"drugs" form:"drugs"`
	DoctorID uint   `json:"doctor_id" form:"doctor_id"`
	Doctor   Doctor `gorm:"foreignKey:DoctorID"`
}
