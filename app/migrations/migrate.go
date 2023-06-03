package migrations

import (
	"github.com/fajaaro/dbo/app/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Customer{}, &models.Order{})
	if err != nil {
		return err
	}
	return nil
}
