package core

import (
	srv "github.com/elvis88/baas/core/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type service interface {
	Register(router *gin.Engine, db *gorm.DB)
}

// Server 提供服务
func Server(router *gin.Engine, db *gorm.DB) {
	services := []service{
		&srv.UserService{},
	}
	for _, service := range services {
		service.Register(router, db)
	}
}
