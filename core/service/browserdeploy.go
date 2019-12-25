package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// BrowserDeployService 浏览器配置表
type BrowserDeployService struct {
	DB *gorm.DB
}

// Register ...
func (srv *BrowserDeployService) Register(api *gin.RouterGroup) {
}
