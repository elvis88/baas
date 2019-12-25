package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserService 用户表
type UserService struct {
	DB *gorm.DB
}

// UserLogin 登陆
func (srv *UserService) UserLogin(ctx *gin.Context) {

}

// UserLogout 退出登陆
func (srv *UserService) UserLogout(ctx *gin.Context) {

}

// UserAuthorize 用户验证
func (srv *UserService) UserAuthorize(ctx *gin.Context) {

}

// UserInfo 用户信息
func (srv *UserService) UserInfo(ctx *gin.Context) {

}

// UserList 用户列表
func (srv *UserService) UserList(ctx *gin.Context) {

}

// UserAdd 新增用户
func (srv *UserService) UserAdd(ctx *gin.Context) {

}

// UserDelete 删除用户
func (srv *UserService) UserDelete(ctx *gin.Context) {

}

// UserUpdate 修改
func (srv *UserService) UserUpdate(ctx *gin.Context) {

}

// Register ...
func (srv *UserService) Register(api *gin.RouterGroup) {
	api.POST("/user/login", srv.UserLogin)
	api.POST("/user/logout", srv.UserLogout)
	//认证校验
	api.Use(srv.UserAuthorize)
	api.GET("/user/info", srv.UserInfo)
	api.GET("/user/list", srv.UserList)
	api.POST("/user/add", srv.UserAdd)
	api.POST("/user/delete", srv.UserDelete)
	api.POST("/user/update", srv.UserUpdate)
}
