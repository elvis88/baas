package core

import (
	srv "github.com/elvis88/baas/core/service"
	"github.com/gin-gonic/gin"
)

type service interface {
	Register(router *gin.Engine)
}

// Server 提供服务
func Server(router *gin.Engine) {
	services := []service{
		&srv.UserService{},
		&srv.RoleService{},
		&srv.ChainService{},
		&srv.ChainDeployService{},
		&srv.BrowserService{},
		&srv.BrowserDeployService{},
	}
	for _, service := range services {
		service.Register(router)
	}
}
