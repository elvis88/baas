package model

import (
	"time"
)

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created,omitempty"`
	UpdatedAt time.Time  `json:"updated,omitempty"`
	DeletedAt *time.Time `json:"deleted,omitempty" sql:"index"`
}

// User 用户表
type User struct {
	Model

	Name      string `json:"name" gorm:"type:varchar(100);not null;unique"`
	Password  string `json:"pwd" gorm:"not null"`
	Nick      string `json:"nick" gorm:"not null"`
	Email     string `json:"email;unique"`
	Telephone string `json:"tel" gorm:"column:tel;unique"`

	Roles          []*Role          `json:"role,omitempty" gorm:"many2many:user_role;"`
	Chains         []*Chain         `json:"chain,omitempty"`
	ChainDeploys   []*ChainDeploy   `json:"chaindeploy,omitempty"`
	Browsers       []*Browser       `json:"browser,omitempty"`
	BrowserDeploys []*BrowserDeploy `json:"browserdeploy,omitempty"`
}

// Role 角色表
type Role struct {
	Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`
}

// Chain 区块链表
type Chain struct {
	Model

	Name        string `gorm:"not null;unique"`
	Description string `gorm:"column:desc"`

	UserID      uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	ChainDeploy []ChainDeploy
}

// ChainDeploy 区块链部署表
type ChainDeploy struct {
	Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`

	UserID  uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	ChainID uint `sql:"type:integer REFERENCES t_chain(id) on update no action on delete no action"`
}

// Browser 浏览器表
type Browser struct {
	Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`

	UserID        uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	BrowserDeploy []BrowserDeploy
}

// BrowserDeploy 浏览器部署表
type BrowserDeploy struct {
	Model

	Name        string `gorm:"type:varchar(100);not null;unique"`
	Description string `gorm:"column:desc"`

	UserID    uint `sql:"type:integer REFERENCES t_user(id) on update no action on delete no action"`
	BrowserID uint `sql:"type:integer REFERENCES t_browser(id) on update no action on delete no action"`
}
