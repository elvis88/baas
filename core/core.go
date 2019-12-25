package core

import (
	"github.com/elvis88/baas/core/model"
	srv "github.com/elvis88/baas/core/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type service interface {
	Register(router *gin.RouterGroup)
}

// Server 提供服务
func Server(router *gin.Engine, db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{}, model.Role{},
		model.Chain{}, model.ChainDeploy{},
		model.Browser{}, model.BrowserDeploy{}).Error; err != nil {
		return err
	}
	services := []service{
		&srv.UserService{
			DB: db,
		},
		&srv.RoleService{
			DB: db,
		},
		&srv.ChainService{
			DB: db,
		},
		&srv.ChainDeployService{
			DB: db,
		},
		&srv.BrowserService{
			DB: db,
		},
		&srv.BrowserDeployService{
			DB: db,
		},
	}

	apiv1 := router.Group("api/v1")
	for _, service := range services {
		service.Register(apiv1)
	}
	return nil
}
