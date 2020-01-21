package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Add 加入
func (c *Agent) Add(db *gorm.DB) error {
	user := &User{}
	if err := db.First(user, c.UserID).Error; err != nil {
		return err
	}
	return db.Create(c).Error
}

// Remove 移除
func (c *Agent) Remove(db *gorm.DB) error {
	if c.ID <= 0 {
		return fmt.Errorf("not support ")
	}

	var chainDeployNodes []*ChainDeployNode
	if err := db.Where(&ChainDeployNode{
		AgentID: c.ID,
	}).Find(&chainDeployNodes).Error; err != nil {
		return err
	}
	for _, chainDeployNode := range chainDeployNodes {
		if err := chainDeployNode.Remove(db); err != nil {
			return err
		}
	}

	return db.Unscoped().Delete(c).Error
}
