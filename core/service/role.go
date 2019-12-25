package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RoleService 角色表
type RoleService struct {
	DB *gorm.DB
}

// RoleAdd 新增
func (srv *RoleService) RoleAdd(ctx *gin.Context) {

}

// RoleDelete 删除
func (srv *RoleService) RoleDelete(ctx *gin.Context) {

}

// RoleUpdate 修改
func (srv *RoleService) RoleUpdate(ctx *gin.Context) {

}

// Register ...
func (srv *RoleService) Register(api *gin.RouterGroup) {
	api.POST("/role/add", srv.RoleAdd)
	api.POST("/role/delete", srv.RoleDelete)
	api.POST("/role/update", srv.RoleUpdate)
}
