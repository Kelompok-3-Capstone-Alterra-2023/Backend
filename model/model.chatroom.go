package model

import "gorm.io/gorm"

type Chatroom struct {
	UserID       uint
	DoctorID     uint
	StatusAccess string `json:"status_access" from:"status_access"`
}

func (Chatroom *Chatroom) BeforeSave(db *gorm.DB) error {
	Chatroom.StatusAccess = "access"
	return nil
}
