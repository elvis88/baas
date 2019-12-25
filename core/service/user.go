package service

import (
	"fmt"
	"time"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/common/jwt"
	"github.com/elvis88/baas/common/password"
	"github.com/elvis88/baas/core/model"
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

// LoginResponse 登陆响应参数
type LoginResponse struct {
	Token string `json:"token"`
}

// UserLogin 登陆
func (srv *UserService) UserLogin(ctx *gin.Context) {
	login := &LoginRequest{}
	if err := ctx.ShouldBindJSON(login); err != nil {
		logger.Error(err)
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	usr := &model.User{}
	if len(login.UserName) != 0 { // 密码登陆
		if err := srv.DB.Where(&model.User{
			Name: login.UserName,
		}).First(usr).Error; err != nil {
			logger.Error(err)
			ginutil.Response(ctx, NAME_NOT_EXIST, err.Error())
			return
		}
		if val, err := password.Validate(login.Password, usr.Password); !val {
			logger.Error(err)
			ginutil.Response(ctx, PASSWORD_WRONG, err.Error())
			return
		}
	} else if len(login.Email) != 0 { // 邮箱验证码登陆
		if err := srv.DB.Where(&model.User{
			Email: login.Email,
		}).First(usr).Error; err != nil {
			logger.Error(err)
			ginutil.Response(ctx, EMAIL_NOT_EXIST, err.Error())
			return
		}
	} else if len(login.Telephone) != 0 { // 手机验证码登陆
		if err := srv.DB.Where(&model.User{
			Telephone: login.Telephone,
		}).First(usr).Error; err != nil {
			logger.Error(err)
			ginutil.Response(ctx, TEL_NOT_EXIST, err.Error())
			return
		}
	} else {
		ginutil.Response(ctx, UNKOWN_TYPE, nil)
		return
	}

	now := time.Now()
	info := make(map[string]interface{})
	info["userId"] = usr.ID
	info["exp"] = now.Add(time.Hour * 1).Unix() // 1 小时过期
	info["iat"] = now.Unix()
	token, err := jwt.CreateToken(TokenKey, info)
	if err != nil {
		logger.Error(err)
		ginutil.Response(ctx, LOGIN_FAIL, err.Error())
		return
	}

	ginutil.SetSession(ctx, token, usr.ID)
	ginutil.Response(ctx, nil, &LoginResponse{
		Token: token,
	})
	return
}

// UserLogout 退出登陆
func (srv *UserService) UserLogout(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	ginutil.RemoveSession(ctx, token)
	ginutil.Response(ctx, nil, nil)
}

// UserAuthorize 用户验证
func (srv *UserService) UserAuthorize(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	session := ginutil.GetSession(ctx, token)
	if nil == session {
		ginutil.Response(ctx, TOKEN_NOT_EXIST, nil)
		ctx.Abort()
		return
	}

	info, ok := jwt.ParseToken(token, TokenKey)
	if !ok {
		ginutil.Response(ctx, TOKEN_INVALID, nil)
		ctx.Abort()
		return
	}

	if infoMap := info.(map[string]interface{}); float64(time.Now().Unix()) >= infoMap["exp"].(float64) {
		ginutil.Response(ctx, TOKEN_EXPIRE, nil)
		ctx.Abort()
		return
	}

	usr := &model.User{}
	if err := srv.DB.Where(&model.User{
		Model: model.Model{
			ID: session.(uint),
		},
	}).First(usr).Error; err != nil {
		logger.Error(err)
		ginutil.Response(ctx, ID_NOT_EXIST, err.Error())
		ctx.Abort()
		return
	}

	ctx.Next()
}

// UserInfo 用户信息
func (srv *UserService) UserInfo(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	session := ginutil.GetSession(ctx, token)
	usr := &model.User{}
	srv.DB.Where(&model.User{
		Model: model.Model{
			ID: session.(uint),
		},
	}).First(usr)
	ginutil.Response(ctx, nil, usr)
}

type PagerRequest struct {
	Page     int `json:"page"`     //当前页
	PageSize int `json:"pageSize"` //每页条数
	Total    int `json:"total"`    //总条数
}

// UserList 用户列表
func (srv *UserService) UserList(ctx *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Error(err)
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}
	usrs := &model.User{}
	offset := req.Page * req.PageSize
	if err := srv.DB.Offset(offset).Limit(req.PageSize).Find(usrs).Error; err != nil {
		logger.Error(err)
		ginutil.Response(ctx, GET_FAIL, err.Error())
		return
	}
	ginutil.Response(ctx, nil, usrs)
	return
}

// UserAdd 新增用户
func (srv *UserService) UserAdd(ctx *gin.Context) {
	usr := &model.User{}
	if err := ctx.ShouldBindJSON(usr); err != nil {
		logger.Error(err)
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}
	if err := srv.DB.Create(&usr).Error; err != nil {
		logger.Error(err)
		ginutil.Response(ctx, ADD_FAIL, err.Error())
		return
	}
	ginutil.Response(ctx, nil, usr)
	return
}

// UserDelete 删除用户
func (srv *UserService) UserDelete(ctx *gin.Context) {
	usr := &model.User{}
	if err := ctx.ShouldBindJSON(usr); err != nil {
		logger.Error(err)
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}
	fmt.Println("delete", usr.ID)
	if err := srv.DB.Unscoped().Delete(&usr).Error; err != nil {
		logger.Error(err)
		ginutil.Response(ctx, DELETE_FAIL, err.Error())
		return
	}
	ginutil.Response(ctx, nil, nil)
	return
}

// UserUpdate 修改
func (srv *UserService) UserUpdate(ctx *gin.Context) {
	usr := &model.User{}
	if err := ctx.ShouldBindJSON(usr); err != nil {
		logger.Error(err)
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	if err := srv.DB.Model(&model.User{
		Model: model.Model{
			ID: usr.ID,
		},
	}).Updates(usr).Error; err != nil {
		logger.Error(err)
		ginutil.Response(ctx, UPDATE_FAIL, err.Error())
		return
	}

	nusr := &model.User{}
	if err := srv.DB.Where(&model.User{
		Model: model.Model{
			ID: usr.ID,
		},
	}).First(nusr).Error; err != nil {
		ginutil.Response(ctx, UPDATE_FAIL, err.Error())
	} else {
		ginutil.Response(ctx, nil, nusr)
	}

	return
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
	//api.Use(srv.UserAuthorize)
	api.POST("/user/info", srv.UserInfo)
	api.POST("/user/list", srv.UserList)
	api.POST("/user/delete", srv.UserDelete)
	api.POST("/user/update", srv.UserUpdate)
	api.POST("/user/changepwd", srv.UserChangePWD)
}
