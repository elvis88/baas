package model

import (
	"errors"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserService 用户表
type UserService struct {
	DB *gorm.DB
}

// LoginRequest 登陆请求参数
type LoginRequest struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Telephone string `json:"phone"`
	Email     string `json:"email"`
	Code      string `json:"code"`
}

var (
	LOGIN_NAME_NOT_EXIST = errors.New("用户不存在")
	LOGIN_PASSWORD_WRONG = errors.New("密码不正确")
	LOGIN_CODE_WRONG     = errors.New("验证码不正确")
	LOGIN_UNKOWN_TYPE    = errors.New("未知登陆方式")
)

// UserLogin 登陆
func (srv *UserService) UserLogin(ctx *gin.Context) {
	login := &LoginRequest{}
	if err := ctx.ShouldBindJSON(login); err != nil {
		ginutil.Response(ctx, err, nil)
		return
	}

	if len(login.UserName) != 0 {
		// 密码登陆

	} else if len(login.Email) != 0 {
		// 邮箱验证码登陆

	} else if len(login.Telephone) != 0 {
		// 手机验证码登陆
	}
	ginutil.Response(ctx, LOGIN_UNKOWN_TYPE, nil)
	return
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

// UserChangePWD 修改密码
func (srv *UserService) UserChangePWD(ctx *gin.Context) {

}

// Register ...
func (srv *UserService) Register(api *gin.RouterGroup) {
	api.POST("/user/register", srv.UserAdd)
	api.POST("/user/login", srv.UserLogin)
	api.POST("/user/logout", srv.UserLogout)
	//认证校验
	api.Use(srv.UserAuthorize)
	api.GET("/user/info", srv.UserInfo)
	api.GET("/user/list", srv.UserList)
	api.POST("/user/delete", srv.UserDelete)
	api.POST("/user/update", srv.UserUpdate)
	api.POST("/user/changepwd", srv.UserChangePWD)
}
