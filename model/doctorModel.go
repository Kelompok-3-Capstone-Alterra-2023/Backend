package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Doctor struct {
	gorm.Model
	Email         string `json:"email" form:"email" gorm:"type:varchar(255)unique;not null"`
	Full_Name        string `json:"full_name" form:"full_name" gorm:"type:varchar(255)"`
	Display_Name        string `json:"display_name" form:"display_name" gorm:"type:varchar(50)"`
	Password      string `json:"password" form:"password" gorm:"not null"`
	Alumnus        string `json:"alumnus" form:"alumnus" gorm:"type:varchar(255)"`
	Work        string `json:"work" form:"work" gorm:"type:varchar(100)"`
	Date_of_Entry time.Time `json:"date_of_entry" form:"date_of_entry" gorm:"type:date"`
	Date_of_Out time.Time `json:"date_of_out" form:"date_of_out" gorm:"type:date"`
	Practice_Adress        string `json:"practice_address" form:"practice_address" gorm:"type:text"`
	Price float64 `json:"price" form:"price" gorm:"type:double"`
	Balance float64 `json:"balance" form:"balance" gorm:"type:double"`
	Photo        string `json:"photo" form:"photo" gorm:"type:text"`
	Status_Online bool   `json:"status_online" form:"status_online" gorm:"type:boolean"`
}