package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// BrowserService 浏览器配置服务
type BrowserService struct {
	DB *gorm.DB
}

// BrowserAdd 新增
func (srv *BrowserService) BrowserAdd(ctx *gin.Context) {

}

// BrowserDelete 删除
func (srv *BrowserService) BrowserDelete(ctx *gin.Context) {

}

// BrowserUpdate 修改
func (srv *BrowserService) BrowserUpdate(ctx *gin.Context) {

}

// Register ...
func (srv *BrowserService) Register(api *gin.RouterGroup) {
	api.POST("/browser/add", srv.BrowserAdd)
	api.POST("/browser/delete", srv.BrowserDelete)
	api.POST("/browser/update", srv.BrowserUpdate)
}
