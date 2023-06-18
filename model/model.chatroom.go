package model

import "gorm.io/gorm"

type ChatRoom struct {
	UserID       uint
	DoctorID     uint
	StatusAccess string `json:"status_access" from:"status_access"`
}

func (Chatroom *ChatRoom) BeforeSave(db *gorm.DB) error {
	Chatroom.StatusAccess = "access"
	return nil
}
