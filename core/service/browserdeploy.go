package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// BrowserDeployService 浏览器配置表
type BrowserDeployService struct {
	DB *gorm.DB
}

// BrowserDeployAdd 新增
func (srv *BrowserDeployService) BrowserDeployAdd(ctx *gin.Context) {

}

// BrowserDeployDelete 删除
func (srv *BrowserDeployService) BrowserDeployDelete(ctx *gin.Context) {

}

// BrowserDeployUpdate 修改
func (srv *BrowserDeployService) BrowserDeployUpdate(ctx *gin.Context) {

}

// Register ...
func (srv *BrowserDeployService) Register(api *gin.RouterGroup) {
	api.POST("/chaindeploy/add", srv.BrowserDeployAdd)
	api.POST("/chaindeploy/delete", srv.BrowserDeployDelete)
	api.POST("/chaindeploy/update", srv.BrowserDeployUpdate)
}
