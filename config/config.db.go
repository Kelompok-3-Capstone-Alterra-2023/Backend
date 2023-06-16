package config

import (
	"capstone/model"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Open() error {
	//load env file
	errenv := godotenv.Load()

	if errenv != nil {
		log.Fatal("error load env file")
	}

	//connect db
	dbUsername := os.Getenv("DBUSERNAME")
	dbPassword := os.Getenv("DBPASSWORD")
	dbHost := os.Getenv("DBHOST")
	dbName := os.Getenv("DBNAME")
	dbPort := os.Getenv("DBPORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	InitMigrate()
	return nil
}


func InitMigrate(){
  	DB.AutoMigrate(model.Doctor{})
	DB.AutoMigrate(model.Article{})
	DB.AutoMigrate(model.User{})
	DB.AutoMigrate(model.UserOTP{})
	DB.AutoMigrate(model.Admin{})
	DB.AutoMigrate(model.Recipt{})
	DB.AutoMigrate(model.Drug{})
	DB.AutoMigrate(model.Order{}, model.ConsultationSchedule{}, model.Payment{})
}
