package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainDeployService 区块链配置表
type ChainDeployService struct {
	DB *gorm.DB
}

// ChainDeployAdd 新增
func (srv *ChainDeployService) ChainDeployAdd(ctx *gin.Context) {

}

// ChainDeployDelete 删除
func (srv *ChainDeployService) ChainDeployDelete(ctx *gin.Context) {

}

// ChainDeployUpdate 修改
func (srv *ChainDeployService) ChainDeployUpdate(ctx *gin.Context) {

}

// Register ...
func (srv *ChainDeployService) Register(api *gin.RouterGroup) {
	api.POST("/chaindeploy/add", srv.ChainDeployAdd)
	api.POST("/chaindeploy/delete", srv.ChainDeployDelete)
	api.POST("/chaindeploy/update", srv.ChainDeployUpdate)
}
