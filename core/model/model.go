package model

import (
	"github.com/elvis88/baas/common/log"
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

// 区块链表结构
type Chain struct {
	Id 			int64 	`gorm:"primary_key;AUTO_INCREMENT"`
	Name 		string 	`gorm:"not null;unique"`
	Userid     	User   	`gorm:"ForeignKey:userID;column:uid"`
	Description string	`gorm:"column:desc"`
}

func (c Chain) TableName() string {
	return "t_chain"
}

func init() {

	// 创建用户表
	table := db.DB.HasTable(User{})
	if !table {
		db.DB.CreateTable(User{})
	}

	// 创建区块链表
	if ok := db.DB.HasTable(Chain{}); !ok {
		if err := db.DB.CreateTable(Chain{}).Error; err != nil {
			log.Log.Debug("Create chain table fail")
		}
	}

}
