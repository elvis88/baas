package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserService 用户表
type UserService struct {
}

// Register
func (srv *UserService) Register(router *gin.Engine, db *gorm.DB) {
}
