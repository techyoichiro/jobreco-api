package model

import (
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name      string `gorm:"size:100;not null"`
	LoginID   string `gorm:"size:50;unique;not null"`
	Password  string `gorm:"size:255;not null"`
	RoleID    int    `gorm:"not null"`
	HourlyPay int    `gorm:"not null"`
}
