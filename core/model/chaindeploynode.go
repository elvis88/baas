package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Add 加入
func (c *ChainDeployNode) Add(db *gorm.DB) error {
	agent := &Agent{}
	if err := db.First(agent, c.AgentID).Error; err != nil {
		return err
	}

	chainDeploy := &ChainDeploy{}
	if err := db.First(chainDeploy, c.ChainDeployID).Error; err != nil {
		return err
	}
	return db.Create(c).Error
}

// Remove 移除
func (c *ChainDeployNode) Remove(db *gorm.DB) error {
	if c.ID <= 0 {
		return fmt.Errorf("not support ")
	}

	var err error
	err = db.Unscoped().Where(&ChainDeployNodeStatus{
		ChainDeployNodeID: c.ID,
	}).Delete(&ChainDeployNodeStatus{}).Error
	if err != nil {
		return err
	}

	err = db.Unscoped().Delete(c).Error
	return err
}
