package models

import (
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model `json:"-"`
	ID         int     `gorm:"primaryKey" json:"id"`
	UID        int     `gorm:"foreignKey" json:"uid"`
	Operation  string  `json:"operation"`
	Amount     float64 `json:"amount"`
	Source     string  `json:"source"`
	Target     string  `json:"target"`
	Category   string  `json:"category"`
	Note       string  `json:"note"`
	Datestring string  `json:"datestring"`
	Timestamp  int     `json:"timestamp"`
}
