package model

import (
	"github.com/jinzhu/gorm"
)

// Add 加入
func (c *ChainDeployNodeStatus) Add(db *gorm.DB) error {
	chaindeploynode := &ChainDeployNode{}
	if err := db.First(chaindeploynode, c.ChainDeployNodeID).Error; err != nil {
		return err
	}

	chainstatus := &ChainStatus{}
	if err := db.First(chainstatus, c.ChainStatusID).Error; err != nil {
		return err
	}
	return db.Create(c).Error
}
