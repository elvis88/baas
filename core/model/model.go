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

type Role struct {
	ID   int64  `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(100);not null;unique"`
}

type UserRole struct {
	ID   int64 `gorm:"primary_key;auto_increment"`
	Uid  uint  `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User User  `gorm:"ForeignKey:Uid"`
	Rid  uint  `sql:"type:integer REFERENCES t_role(id) on update no action on delete no action"`
	Role Role  `gorm:"ForeignKey:Rid"`
}

// 区块链表结构
type Chain struct {
	Id          int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Name        string `gorm:"not null;unique"`
	Uid         uint   `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User        User   `gorm:"ForeignKey:Uid"`
	Description string `gorm:"column:desc"`
}

type ChainDeploy struct {
	ID    int64  `gorm:"primary_key;auto_increment"`
	Name  string `gorm:"type:varchar(100);not null;unique"`
	Uid   uint   `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User  User   `gorm:"ForeignKey:Uid"`
	Cid   uint   `sql:"type:integer REFERENCES t_chain(id) on update no action on delete no action"`
	Chain Chain  `gorm:"ForeignKey:Cid"`
	Desc  string `gorm:"type:varchar(255)"`
}

type Browser struct {
	ID   int64  `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(100);not null;unique"`
	Uid  uint   `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User User   `gorm:"ForeignKey:Uid"`
	Desc string `gorm:"type:varchar(255)"`
}

type BrowserDeploy struct {
	ID      int64   `gorm:"primary_key;auto_increment"`
	Name    string  `gorm:"type:varchar(100);not null;unique"`
	Uid     uint    `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User    User    `gorm:"ForeignKey:Uid"`
	Bid     uint    `sql:"type:integer REFERENCES t_browser(id) on update no action on delete no action"`
	Browser Browser `gorm:"ForeignKey:Bid"`
	Desc    string  `gorm:"type:varchar(255)"`
}

func ModelInit() {
	db.DB.AutoMigrate(&User{}, Role{}, UserRole{}, Browser{}, BrowserDeploy{}, Chain{}, ChainDeploy{})
}
