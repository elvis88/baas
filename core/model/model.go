package model

import (
	"time"

	"github.com/elvis88/baas/db"
)

// User 用户表
type User struct {
	// gorm.Model
	ID         int64  `gorm:"primary_key;auto_increment"`
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

type Role struct {
	ID   int64  `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(100);not null;unique"`
}

func (r Role) TableName() string {
	return "t_role"
}

type UserRole struct {
	ID  int64 `gorm:"primary_key;auto_increment"`
	Uid User  `gorm:"ForeignKey:ID`
	Rid Role  `gorm:"ForeignKey:ID`
}

func (u UserRole) TableName() string {
	return "t_user_role"
}

type Browser struct {
	ID   int64  `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(100);not null;unique"`
	Uid  User   `gorm:"ForeignKey:ID`
	Desc string `gorm:"type:varchar(255)"`
}

func (b Browser) TableName() string {
	return "t_browser"
}

type BrowserDeploy struct {
	ID   int64   `gorm:"primary_key;auto_increment"`
	Name string  `gorm:"type:varchar(100);not null;unique"`
	Uid  User    `gorm:"ForeignKey:ID`
	Desc string  `gorm:"type:varchar(255)"`
	Bid  Browser `gorm:"ForeignKey:ID`
}

func (b BrowserDeploy) TableName() string {
	return "t_browser_deploy"
}

func init() {
	if table := db.DB.HasTable(User{}); !table {
		db.DB.CreateTable(&User{})
	}
	if table := db.DB.HasTable(Role{}); !table {
		db.DB.CreateTable(&Role{})
	}
	if table := db.DB.HasTable(UserRole{}); !table {
		db.DB.CreateTable(&UserRole{})
	}
	if table := db.DB.HasTable(Browser{}); !table {
		db.DB.CreateTable(&Browser{})
	}
	if table := db.DB.HasTable(BrowserDeploy{}); !table {
		db.DB.CreateTable(&BrowserDeploy{})
	}
}
