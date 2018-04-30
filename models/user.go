package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username     string         `gorm:"column:username"`
	Email        string         `gorm:"column:email;unique_index"`
	Bio          string         `gorm:"column:bio;size:1024"`
	Image        sql.NullString `gorm:"column:image"`
	PasswordHash string         `gorm:"column:password;not null"`
}
