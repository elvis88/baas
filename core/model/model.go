package model

import (
	"time"

	"github.com/elvis88/baas/db"
)

// User 用户表
type User struct {
	// gorm.Model
	ID         int64  `gorm:"primary_key"`
	Name       string `gorm:"type:varchar(100);not null;unique"`
	Pwd        string `gorm:"not null"`
	Nick       string `gorm:"not null"`
	Email      string
	Tele       string
	CreateTime time.Time `gorm:"column:create_time"`
}

func (u User) TableName() string {
	return "t_user"
}

func init() {
	table := db.DB.HasTable(User{})
	if !table {
		db.DB.CreateTable(User{})
	}
}
