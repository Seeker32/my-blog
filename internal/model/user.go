package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
	Role     int    `gorm:"default: 2" json:"role"` // 1.admin 2.user
}
