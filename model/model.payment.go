package model

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	OrderID         uint    `json:"order_id" form:"order_id"`
	TotalPrice      float64 `json:"total_price" form:"total_price"`
	PaymentMethod   string  `json:"payment_method" form:"payment_method"`
	TransferStatus  string  `json:"transfer_status" form:"transfer_status"`
	TransactionTime string  `json:"transaction_time" form:"transaction_time"`
	Order           Order   `gorm:"foreignKey:OrderID"`
}

type Notification struct {
	OrderID         string `json:"order_id"`
	PaymentType     string `json:"payment_type"`
	PaymentStatus   string `json:"transaction_status"`
	TransactionTime string `json:"transaction_time"`
}
