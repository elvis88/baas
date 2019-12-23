package model

import (
	"time"

	"github.com/elvis88/baas/db"
	"github.com/jinzhu/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Name       string `gorm:"type:varchar(100);not null;unique"`
	Pwd        string `gorm:"not null"`
	Nick       string `gorm:"not null"`
	Email      string
	Tele       string
	CreateTime time.Time `gorm:"column:create_time"`
	Role       []Role    `gorm:"many2many:user_role;"`
}

type Role struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null;unique"`
}

// 区块链表结构
type Chain struct {
	gorm.Model
	Name        string `gorm:"not null;unique"`
	UserID      uint
	Description string `gorm:"column:desc"`
}

type ChainDeploy struct {
	gorm.Model
	Name  string `gorm:"type:varchar(100);not null;unique"`
	Uid   uint   `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User  User   `gorm:"ForeignKey:Uid"`
	Cid   uint   `sql:"type:integer REFERENCES t_chain(id) on update no action on delete no action"`
	Chain Chain  `gorm:"ForeignKey:Cid"`
	Desc  string `gorm:"type:varchar(255)"`
}

type Browser struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null;unique"`
	Uid  uint   `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User User   `gorm:"ForeignKey:Uid"`
	Desc string `gorm:"type:varchar(255)"`
}

type BrowserDeploy struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(100);not null;unique"`
	Uid     uint    `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User    User    `gorm:"ForeignKey:Uid"`
	Bid     uint    `sql:"type:integer REFERENCES t_browser(id) on update no action on delete no action"`
	Browser Browser `gorm:"ForeignKey:Bid"`
	Desc    string  `gorm:"type:varchar(255)"`
}

func ModelInit() {
	db.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, Role{}, Browser{}, BrowserDeploy{}, Chain{}, ChainDeploy{})
}
