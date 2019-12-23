package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RoleService 角色表
type RoleService struct {
}

// Register
func (srv *RoleService) Register(router *gin.Engine, db *gorm.DB) {
}
