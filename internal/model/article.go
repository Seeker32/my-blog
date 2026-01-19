package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title      string   `gorm:"type:varchar(100);not null" json:"title"`
	Desc       string   `gorm:"type:varchar(255)" json:"desc"`
	Content    string   `gorm:"type:longtext;not null" json:"content"`
	Status     int      `gorm:"default:0" json:"status"`
	UserID     uint     `gorm:"not null" json:"user_id"`
	User       User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags       []Tag    `gorm:"many2many:article_tags" json:"tags,omitempty"`
}
