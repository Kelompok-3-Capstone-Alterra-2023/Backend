package model

import "gorm.io/gorm"

type Withdraw struct {
	gorm.Model
	ReferenceNumber string  `json:"reference_number"`
	DoctorID        uint    `json:"doctor_id"`
	Method          string  `json:"method"`
	Bank            string  `json:"bank"`
	AccountName     string  `json:"account_name"`
	AccountNumber   string  `json:"account_number"`
	Amount          float64 `json:"amount" gorm:"type:double"`
	TransactionFee  float64 `json:"transaction_fee" gorm:"type:double"`
	Total           float64 `json:"total" gorm:"type:double"`
	Status          string  `json:"status"`
	Notes           string  `json:"notes"`
	Doctor          Doctor  `gorm:"foreignKey:DoctorID"`
}

type WithdrawResponse struct {
	ReferenceNumber string  `json:"reference_number"`
	Method          string  `json:"method"`
	Bank            string  `json:"bank"`
	AccountName     string  `json:"account_name"`
	AccountNumber   string  `json:"account_number"`
	Amount          float64 `json:"amount"`
	TransactionFee  float64 `json:"transaction_fee"`
	Total           float64 `json:"total"`
	Date            string  `json:"date"`
}

type WithdrawsResponse struct {
	ReferenceNumber string  `json:"reference_number"`
	Method          string  `json:"method"`
	Bank            string  `json:"bank"`
	DoctorName      string  `json:"doctor_name"`
	DoctorEmail     string  `json:"doctor_email"`
	AccountNumber   string  `json:"account_number"`
	Amount          float64 `json:"amount"`
	TransactionFee  float64 `json:"transaction_fee"`
	Total           float64 `json:"total"`
	Date            string  `json:"date"`
}

// type PayoutRequest struct {
// 	BeneficiaryName    string  `json:"beneficiary_name"`
// 	BeneficiaryAccount string  `json:"beneficiary_account"`
// 	BeneficiaryEmail   string  `json:"beneficiary_email"`
// 	BeneficiaryBank    string  `json:"beneficiary_bank"`
// 	Amount             float64 `json:"amount"`
// 	Notes              string  `json:"notes"`
// }
