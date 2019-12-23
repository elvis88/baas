package model

import (
	"github.com/elvis88/baas/common/log"
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
	ID   int64 `gorm:"primary_key;auto_increment"`
	Uid  uint  `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User User  `gorm:"ForeignKey:Uid"`
	Rid  uint  `sql:"type:integer REFERENCES t_role(id) on update no action on delete no action"`
	Role Role  `gorm:"ForeignKey:Rid"`
}

func (u UserRole) TableName() string {
	return "t_user_role"
}

// 区块链表结构
type Chain struct {
	Id 			int64 	`gorm:"primary_key;AUTO_INCREMENT"`
	Name 		string 	`gorm:"not null;unique"`
	Uid         uint    `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User     	User   	`gorm:"ForeignKey:Uid"`
	Description string	`gorm:"column:desc"`
}

func (c Chain) TableName() string {
	return "t_chain"
}

type ChainDeploy struct {
	ID    int64  `gorm:"primary_key;auto_increment"`
	Name  string `gorm:"type:varchar(100);not null;unique"`
	Uid         uint    `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User     	User   	`gorm:"ForeignKey:Uid"`
	Cid   uint   `sql:"type:integer REFERENCES t_chain(id) on update no action on delete no action"`
	Chain Chain  `gorm:"ForeignKey:Cid"`
	Desc  string `gorm:"type:varchar(255)"`
}

func (c ChainDeploy) TableName() string {
	return "t_chain_deploy"
}

type Browser struct {
	ID   int64  `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(100);not null;unique"`
	Uid         uint    `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	User     	User   	`gorm:"ForeignKey:Uid"`
	Desc string `gorm:"type:varchar(255)"`
}

func (b Browser) TableName() string {
	return "t_browser"
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

func (b BrowserDeploy) TableName() string {
	return "t_browser_deploy"
}

func ModelInit() {
	// 创建用户表
	if ok := db.DB.HasTable(User{}); !ok {
		if err := db.DB.CreateTable(&User{}).Error; err != nil {
			log.Log.Debug("Create %s table fail", User{}.TableName())
		}
	}

	// 创建角色表
	if ok := db.DB.HasTable(Role{}); !ok {
		if err := db.DB.CreateTable(&Role{}).Error; err != nil {
			log.Log.Debug("Create %s table fail", Role{}.TableName())
		}
	}

	// 用户角色关联表
	if ok := db.DB.HasTable(UserRole{}); !ok {
		if err := db.DB.CreateTable(&UserRole{}).Error; err != nil {
			log.Log.Debug("Create %s table fail", UserRole{}.TableName())
		}
	}

	// 创建浏览器表
	if ok := db.DB.HasTable(Browser{}); !ok {
		if err := db.DB.CreateTable(&Browser{}).Error; err != nil {
			log.Log.Debug("Create %s table fail", Browser{}.TableName())
		}
	}

	// 创建浏览器配置表
	if ok := db.DB.HasTable(BrowserDeploy{}); !ok {
		if err := db.DB.CreateTable(&BrowserDeploy{}).Error; err != nil {
			log.Log.Debug("Create %s table fail", BrowserDeploy{}.TableName())
		}
	}

	// 创建区块链表
	if ok := db.DB.HasTable(Chain{}); !ok {
		if err := db.DB.CreateTable(&Chain{}).Error; err != nil {
			log.Log.Debug("Create %s table fail", Chain{}.TableName())
		}
	}

	// 创建区块链配置表

	if ok := db.DB.HasTable(ChainDeploy{}); !ok {
		if err := db.DB.CreateTable(&ChainDeploy{}).Error; err != nil {
			log.Log.Debug("Create %s table fail", ChainDeploy{}.TableName())
		}
	}
}
