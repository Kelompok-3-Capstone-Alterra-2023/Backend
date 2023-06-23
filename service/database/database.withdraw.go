package database

import (
	"capstone/config"
	"capstone/model"
)

func UpdateBalanceDoctor(id string, newBalance float64) error {
	if err := config.DB.Table("doctors").Where("id = ?", id).Update("balance", newBalance).Error; err != nil {
		return err
	}

	return nil
}

func SaveWithdraw(data *model.Withdraw) error {
	return config.DB.Save(data).Error
}

func GetWithdraws() ([]model.Withdraw, error) {
	var withdraws []model.Withdraw
	if err := config.DB.Preload("Doctor").Where("status = \"queued\"").Find(&withdraws).Error; err != nil {
		return nil, err
	}

	return withdraws, nil
}

func GetWithdrawByID(id string) (model.Withdraw, error) {
	var withdraw model.Withdraw
	if err := config.DB.Where("id = ?", id).First(&withdraw).Error; err != nil {
		return withdraw, err
	}

	return withdraw, nil
}

func UpdateStatusWithdraw(id, status string) error {
	return config.DB.Table("withdraws").Where("id = ?", id).Update("status", status).Error
}

func SearchWithdraw(keyword string) ([]model.Withdraw, error) {
	var withdraw []model.Withdraw
	var doctorID string
	config.DB.Table("doctors").Select("id").Where("full_name LIKE ?", "%"+keyword+"%").Scan(&doctorID)
	if err := config.DB.Preload("Doctor").Where("doctor_id = ?", doctorID).Find(&withdraw).Error; err != nil {
		return nil, err
	}

	return withdraw, nil
}
