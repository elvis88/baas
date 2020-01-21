package model

import (
	"github.com/jinzhu/gorm"
)

// User 用户表
type User struct {
	gorm.Model

	Name      string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	Nick      string `gorm:"not null"`
	Email     string `gorm:"unique"`
	Telephone string `gorm:"unique"`

	OwnerChains []*Chain `json:"-" gorm:"many2many:user_chain;"`
}

// Chain 区块链项目表
type Chain struct {
	gorm.Model

	Name        string `gorm:"not null;unique"`
	Description string `gorm:"column:desc"`
	URL         string
	Public      bool
	AncestorID  uint

	UserID uint
}

// ChainStatus 区块链项目状态表
type ChainStatus struct {
	gorm.Model

	Name        string `gorm:"not null;unique"`
	Description string `gorm:"column:desc"`
	RPC         string `gorm:"not null;unique"`

	ChainID uint
}

// ChainDeploy 区块链部署表
type ChainDeploy struct {
	gorm.Model

	Name        string `gorm:"not null;unique"`
	Description string `gorm:"column:desc"`

	ChainDeployNodes []*ChainDeployNode `json:"-"`

	ChainID uint
	UserID  uint
}

// Agent 监控进程表
type Agent struct {
	gorm.Model

	Name        string `gorm:"not null;unique"`
	Description string `gorm:"column:desc"`
	Status      string

	ChainDeployNodes []*ChainDeployNode `json:"-"`

	UserID uint
}

// ChainDeployNode 区块链节点表
type ChainDeployNode struct {
	gorm.Model

	Name        string `gorm:"not null;unique"`
	Description string `gorm:"column:desc"`
	Status      string

	AgentID       uint
	ChainDeployID uint

	UserID uint
}

// ChainDeployNodeStatus 区块链节点状态表
type ChainDeployNodeStatus struct {
	gorm.Model

	Value string

	ChainDeployNodeID uint
	ChainStatusID     uint
}
