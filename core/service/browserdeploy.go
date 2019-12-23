package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// BrowserDeployService 浏览器配置表
type BrowserDeployService struct {
}

// Register
func (srv *BrowserDeployService) Register(router *gin.Engine, db *gorm.DB) {
}
