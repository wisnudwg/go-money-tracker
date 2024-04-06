package initializers

import "go-money-tracker/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Entry{})
}
