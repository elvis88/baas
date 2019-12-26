package core

import (
	"github.com/elvis88/baas/common/ginutil"
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
	// 初始角色
	adminRole := &model.Role{
		Name:        "admin",
		Description: "超级管理员,拥有所有权限",
	}
	userRole := &model.Role{
		Name:        "user",
		Description: "普通用户",
	}
	roles := []*model.Role{
		adminRole,
		userRole,
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(role, &model.Role{
			Name: role.Name,
		}).Error; err != nil {
			return err
		}
	}
	// 初始Amdin用户
	adminUsr := &model.User{
		Name:     "admin",
		Password: "123456",
		Roles: []*model.Role{
			adminRole,
		},
	}
	usrs := []*model.User{
		adminUsr,
	}
	for _, usr := range usrs {
		if err := db.FirstOrCreate(usr, &model.User{
			Name: usr.Name,
		}).Error; err != nil {
			return err
		}
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

	ginutil.UseSession(router)
	apiv1 := router.Group("api/v1")
	for _, service := range services {
		service.Register(apiv1)
	}
	return nil
}
