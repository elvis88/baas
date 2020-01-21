package model

import (
	"fmt"

	"github.com/elvis88/baas/core/generate"
	"github.com/jinzhu/gorm"
)

// Add 新建
func (c *Chain) Add(db *gorm.DB) error {
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
	err = db.Model(&User{
		Model: gorm.Model{
			ID: c.UserID,
		},
	}).Association("OwnerChains").Append(c).Error
	if err != nil {
		return err
	}
	return nil
}

// Remove 移除
func (c *Chain) Remove(db *gorm.DB) error {
	if c.ID <= 0 {
		return fmt.Errorf("not support ")
	}

	spec, err := c.Spec(db)
	if err != nil {
		return err
	}

	var chainDeploys []*ChainDeploy
	if err := db.Where(&ChainDeploy{
		ChainID: c.ID,
	}).Find(&chainDeploys).Error; err != nil {
		return err
	}
	for _, chainDeploy := range chainDeploys {
		chainDeploy.Remove(db)
	}

	err = db.Unscoped().Delete(c).Error
	if err != nil {
		return err
	}

	err = db.Unscoped().Model(&User{
		Model: gorm.Model{
			ID: c.UserID,
		},
	}).Association("OwnerChains").Delete(c).Error
	if err != nil {
		return err
	}
	return spec.Remove()
}

// Join 加入
func (c *Chain) Join(db *gorm.DB, user *User) error {
	spec, err := c.Spec(db)
	if err != nil {
		return err
	}

	err = spec.Join(user.Name)
	if err != nil {
		return err
	}

	// defer func() {
	// 	if err != nil {
	// 		spec.Unjoin(user.Name)
	// 	}
	// }()

	if err = db.Model(&User{
		Model: gorm.Model{
			ID: user.ID,
		},
	}).Association("OwnerChains").Append(c).Error; nil != err {
		return err
	}
	return nil
}

// Unjoin 加入
func (c *Chain) Unjoin(db *gorm.DB, user *User) error {
	spec, err := c.Spec(db)
	if err != nil {
		return err
	}

	var chainDeploys []*ChainDeploy
	if err := db.Where(&ChainDeploy{
		ChainID: c.ID,
	}).Find(&chainDeploys).Error; err != nil {
		return err
	}
	for _, chainDeploy := range chainDeploys {
		chainDeploy.Remove(db)
	}

	if err = db.Unscoped().Model(&User{
		Model: gorm.Model{
			ID: user.ID,
		},
	}).Association("OwnerChains").Delete(c).Error; nil != err {
		return err
	}
	return spec.Unjoin(user.Name)
}

func (c *Chain) Spec(db *gorm.DB) (generate.AppSpec, error) {
	user := &User{}
	if err := db.First(user, c.UserID).Error; err != nil {
		return nil, err
	}

	ancestorChain := &Chain{}
	if err := db.First(ancestorChain, c.AncestorID).Error; err != nil {
		return nil, err
	}

	spec := generate.NewAppSpec(user.Name, c.Name, ancestorChain.Name)
	if spec == nil {
		return nil, fmt.Errorf("not support ancestor chain %s", ancestorChain.Name)
	}
	return spec, nil
}
