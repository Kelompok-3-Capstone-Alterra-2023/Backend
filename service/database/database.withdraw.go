package database

import "capstone/config"

func UpdateBalanceDoctor(id string, newBalance float64) error {
	if err := config.DB.Table("doctors").Where("id = ?", id).Update("balance", newBalance).Error; err != nil {
		return err
	}

	return nil
}
