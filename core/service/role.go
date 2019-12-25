package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RoleService 角色表
type RoleService struct {
	DB *gorm.DB
}

// Register ...
func (srv *RoleService) Register(api *gin.RouterGroup) {
}
