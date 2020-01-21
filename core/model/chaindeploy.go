package model

import (
	"fmt"

	"github.com/elvis88/baas/core/generate"
	"github.com/jinzhu/gorm"
)

// Add 新建
func (c *ChainDeploy) Add(db *gorm.DB) error {
	spec, err := c.Spec(db)
	if err != nil {
		return err
	}

	err = spec.Build()
	if err != nil {
		return err
	}

	// defer func() {
	// 	if err != nil {
	// 		spec.Remove()
	// 	}
	// }()

	err = db.Create(c).Error
	if err != nil {
		return err
	}
	return nil
}

// Remove 移除
func (c *ChainDeploy) Remove(db *gorm.DB) error {
	if c.ID <= 0 {
		return fmt.Errorf("not support ")
	}

	spec, err := c.Spec(db)
	if err != nil {
		return err
	}

	var chainDeployNodes []*ChainDeployNode
	if err := db.Where(&ChainDeployNode{
		ChainDeployID: c.ID,
	}).Find(&chainDeployNodes).Error; err != nil {
		return err
	}
	for _, chainDeployNode := range chainDeployNodes {
		chainDeployNode.Remove(db)
	}

	err = db.Unscoped().Delete(c).Error
	if err != nil {
		return err
	}

	spec.Remove()
	return nil
}

// GetScriptPath 脚本文件
func (c *ChainDeploy) GetScriptPath(db *gorm.DB) (string, error) {
	if c.ID <= 0 {
		return "", fmt.Errorf("not support ")
	}

	spec, err := c.Spec(db)
	if err != nil {
		return "", err
	}

	return spec.GetScriptPath(), nil
}

func (c *ChainDeploy) Spec(db *gorm.DB) (generate.AppDeploySpec, error) {
	user := &User{}
	if err := db.First(user, c.UserID).Error; err != nil {
		return nil, err
	}

	chain := &Chain{}
	if err := db.First(chain, c.ChainID).Error; err != nil {
		return nil, err
	}

	ancestorChain := &Chain{}
	if err := db.First(ancestorChain, chain.AncestorID).Error; err != nil {
		return nil, err
	}

	spec := generate.NewAppDeploySpec(user.Name, c.Name, ancestorChain.Name, chain.Name)
	if spec == nil {
		return nil, fmt.Errorf("not support ancestor chain deploy %s", ancestorChain.Name)
	}
	return spec, nil
}
