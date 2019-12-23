package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainService 区块链配置表
type ChainService struct {
}

// Register
func (srv *ChainService) Register(router *gin.Engine, db *gorm.DB) {
}
