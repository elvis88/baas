package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// BrowserService 浏览器配置服务
type BrowserService struct {
}

// Register
func (srv *BrowserService) Register(router *gin.Engine, db *gorm.DB) {
}
