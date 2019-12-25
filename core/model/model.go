package model

import (
	"github.com/jinzhu/gorm"
)

// User 用户表
type User struct {
	gorm.Model

	Name      string `json:"name" gorm:"type:varchar(100);not null;unique"`
	Password  string `json:"pwd" gorm:"not null"`
	Nick      string `json:"nick" gorm:"not null"`
	Email     string `json:"email"`
	Telephone string `json:"tel" gorm:"column:tel"`

	Role          []Role          `json:"-" gorm:"many2many:user_role;"`
	Chain         []Chain         `json:"-"`
	ChainDeploy   []ChainDeploy   `json:"-"`
	Browser       []Browser       `json:"-"`
	BrowserDeploy []BrowserDeploy `json:"-"`
}

// Role 角色表
type Role struct {
	gorm.Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`
}

// Chain 区块链表
type Chain struct {
	gorm.Model

	Name        string `gorm:"not null;unique"`
	Description string `gorm:"column:desc"`

	UserID      uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	ChainDeploy []ChainDeploy
}

// ChainDeploy 区块链部署表
type ChainDeploy struct {
	gorm.Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`

	UserID  uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	ChainID uint `sql:"type:integer REFERENCES t_chain(id) on update no action on delete no action"`
}

// Browser 浏览器表
type Browser struct {
	gorm.Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`

	UserID        uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	BrowserDeploy []BrowserDeploy
}

// BrowserDeploy 浏览器部署表
type BrowserDeploy struct {
	gorm.Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`

	UserID    uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	BrowserID uint `sql:"type:integer REFERENCES t_browser(id) on update no action on delete no action"`
}
