package model

import (
	"time"

	"github.com/elvis88/baas/db"
	"github.com/jinzhu/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Name          string `gorm:"type:varchar(100);not null;unique"`
	Pwd           string `gorm:"not null"`
	Nick          string `gorm:"not null"`
	Email         string
	Tele          string
	CreateTime    time.Time       `gorm:"column:create_time"`
	Role          []Role          `gorm:"many2many:user_role;"`
	Chain         []Chain         `gorm:"foreignkey:UserID"`
	ChainDeploy   []ChainDeploy   `gorm:"foreignkey:UserID"`
	Browser       []Browser       `gorm:"foreignkey:UserID"`
	BrowserDeploy []BrowserDeploy `gorm:"foreignkey:UserID"`
}

type Role struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null;unique"`
}

// 区块链表结构
type Chain struct {
	gorm.Model
	Name        string        `gorm:"not null;unique"`
	UserID      uint          `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	Description string        `gorm:"column:desc"`
	ChainDeploy []ChainDeploy `gorm:"foreignkey:ChainID"`
}

type ChainDeploy struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null;unique"`
	UserID  uint   `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	ChainID uint   `sql:"type:integer REFERENCES t_chain(id) on update no action on delete no action"`
	Desc    string `gorm:"type:varchar(255)"`
}

type Browser struct {
	gorm.Model
	Name          string          `gorm:"type:varchar(100);not null;unique"`
	UserID        uint            `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	Desc          string          `gorm:"type:varchar(255)"`
	BrowserDeploy []BrowserDeploy `gorm:"foreignkey:BrowserID"`
}

type BrowserDeploy struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);not null;unique"`
	UserID    uint   `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	BrowserID uint   `sql:"type:integer REFERENCES t_browser(id) on update no action on delete no action"`
	Desc      string `gorm:"type:varchar(255)"`
}

func ModelInit() {
	db.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&User{}, Role{}, Browser{}, BrowserDeploy{}, Chain{}, ChainDeploy{})
}
