package database

import (
	"capstone/config"
	"capstone/model"
	"errors"
)

func GetPassword(email, table string) (interface{}, error) {
	if table == "admins" {
		var passForLogin model.Admin
		if err := config.DB.Table(table).Select("password").Where("email = ?", email).First(&passForLogin).Error; err != nil {
			return passForLogin, err
		}

		return passForLogin, nil
	} else if table == "users" {
		var passForLogin model.User
		if err := config.DB.Table(table).Select("password").Where("email = ?", email).First(&passForLogin).Error; err != nil {
			return passForLogin, err
		}

		return passForLogin, nil
	} else if table == "doctors" {
		var passForLogin model.Doctor
		if err := config.DB.Table(table).Select("password").Where("email = ?", email).First(&passForLogin).Error; err != nil {
			return passForLogin, err
		}

		return passForLogin, nil
	}

	return "", errors.New("table not found")

}
