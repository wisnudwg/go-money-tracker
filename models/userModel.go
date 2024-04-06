package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	ID         int    `gorm:"primaryKey" json:"id"`
	Email      string `gorm:"unique" json:"email"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	Phone      string `gorm:"unique" json:"phone"`
}
