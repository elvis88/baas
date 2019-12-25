package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainDeployService 区块链配置表
type ChainDeployService struct {
	DB *gorm.DB
}

// Register ...
func (srv *ChainDeployService) Register(api *gin.RouterGroup) {
}
