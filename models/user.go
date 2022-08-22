package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint         `json:"id" gorm:"primaryKey"`
	Username string       `json:"username" binding:"required" gorm:"unique"`
	Password string       `json:"password" binding:"required"`
	Fullname string       `json:"fullname"`
	Provider string       `json:"provider" gorm:"default:local"`
	IsVerify sql.NullBool `json:"isVerify" gorm:"default:true"`
}
