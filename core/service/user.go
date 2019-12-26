package service

import (
	"fmt"
	"strconv"
	"strings"
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
	UserName  string `json:"name"`
	Password  string `json:"pwd"`
	Telephone string `json:"phone"`
	Email     string `json:"email"`
	Code      string `json:"code"`
}

// UserLogin 登陆
func (srv *UserService) UserLogin(ctx *gin.Context) {
	login := &LoginRequest{}
	if err := ctx.ShouldBindJSON(login); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	usr := &model.User{}
	if len(login.UserName) != 0 { // 密码登陆
		// 用户是否存在
		if err := srv.DB.Where(&model.User{
			Name: login.UserName,
		}).First(usr).Error; err != nil {

			ginutil.Response(ctx, NAME_NOT_EXIST, err.Error())
			return
		}
		// 密码是否正确
		if val, _ := password.Validate(login.Password, usr.Password); !val {

			ginutil.Response(ctx, PASSWORD_WRONG, nil)
			return
		}
	} else if len(login.Code) != 0 { // 验证码登陆
		// 手机/邮箱是否存在
		codesessionkey := ""
		if len(login.Email) != 0 { // 邮箱验证码
			if err := srv.DB.Where(&model.User{
				Email: login.Email,
			}).First(usr).Error; err != nil {

				ginutil.Response(ctx, EMAIL_NOT_EXIST, err.Error())
				return
			}
			codesessionkey = CodeLoginKey + login.Email
		} else if len(login.Telephone) != 0 { // 手机验证码登陆
			if err := srv.DB.Where(&model.User{
				Telephone: login.Telephone,
			}).First(usr).Error; err != nil {

				ginutil.Response(ctx, TEL_NOT_EXIST, err.Error())
				return
			}
			codesessionkey = CodeLoginKey + login.Telephone
		} else if len(codesessionkey) == 0 {
			ginutil.Response(ctx, CODE_UNKOWN_TYPE, nil)
			return
		}

		// 验证码是否存在
		session := ginutil.GetSession(ctx, codesessionkey)
		if nil == session {
			ginutil.Response(ctx, CODE_NOT_EXIST, nil)
			return
		}
		// 验证码是否过期
		info := session.(map[string]interface{})
		if float64(time.Now().Unix()) >= info["exp"].(float64) {
			ginutil.Response(ctx, CODE_EXPIRE, nil)
			return
		}
		// 验证码是否匹配
		if strings.Compare(info["code"].(string), login.Code) != 0 {
			ginutil.Response(ctx, CODE_WRONG, nil)
			return
		}
		ginutil.RemoveSession(ctx, codesessionkey)
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

		ginutil.Response(ctx, LOGIN_FAIL, err.Error())
		return
	}

	ctx.Header(headerTokenKey, token)
	ginutil.SetSession(ctx, token, usr.ID)
	ginutil.Response(ctx, nil, usr)
	return
}

// UserAuthorize 用户验证
func (srv *UserService) UserAuthorize(ctx *gin.Context) {
	token := ctx.GetHeader(headerTokenKey)
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

		ginutil.Response(ctx, ID_NOT_EXIST, err.Error())
		ctx.Abort()
		return
	}

	ctx.Next()
}

// UserLogout 退出登陆
func (srv *UserService) UserLogout(ctx *gin.Context) {
	token := ctx.GetHeader(headerTokenKey)
	ginutil.RemoveSession(ctx, token)
	ginutil.Response(ctx, nil, nil)
}

// UserInfo 用户信息
func (srv *UserService) UserInfo(ctx *gin.Context) {
	token := ctx.GetHeader(headerTokenKey)
	session := ginutil.GetSession(ctx, token)
	usr := &model.User{}
	srv.DB.Preload("Roles").Where(&model.User{
		Model: model.Model{
			ID: session.(uint),
		},
	}).First(usr)
	ginutil.Response(ctx, nil, usr)
}

// PagerRequest ...
type PagerRequest struct {
	Page     int `json:"page"`     //当前页
	PageSize int `json:"pageSize"` //每页条数
}

// UserList 用户列表
func (srv *UserService) UserList(ctx *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}
	usrs := []*model.User{}
	offset := req.Page * req.PageSize
	if has, cusr := srv.hasAdminRole(ctx); !has {
		if err := srv.DB.Preload("Roles").Where(cusr).Offset(offset).Limit(req.PageSize).Find(&usrs).Error; err != nil {
			ginutil.Response(ctx, GET_FAIL, err.Error())
			return
		}
	} else {
		if err := srv.DB.Preload("Roles").Offset(offset).Limit(req.PageSize).Find(&usrs).Error; err != nil {
			ginutil.Response(ctx, GET_FAIL, err.Error())
			return
		}
	}

	ginutil.Response(ctx, nil, usrs)
	return
}

// UserAdd 新增用户
func (srv *UserService) UserAdd(ctx *gin.Context) {
	usr := &model.User{}
	if err := ctx.ShouldBindJSON(usr); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	password, err := password.CryTo(usr.Password, 12, "default")
	if err != nil {

		ginutil.Response(ctx, ADD_FAIL, err.Error())
		return
	}
	usr.Password = password
	userRole := &model.Role{}
	if err := srv.DB.Where(&model.Role{
		Key: "user",
	}).First(userRole).Error; err == nil {
		usr.Roles = append(usr.Roles, userRole)
	}

	if err := srv.DB.Create(&usr).Error; err != nil {

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
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	if has, cusr := srv.hasAdminRole(ctx); !has && cusr.ID != usr.ID {
		ginutil.Response(ctx, NOPERMISSION, nil)
		return
	}

	if res := srv.DB.Unscoped().Delete(&usr); res.RowsAffected == 0 {
		errstr := ""
		if res.Error != nil {
			errstr = res.Error.Error()
		}
		ginutil.Response(ctx, UPDATE_FAIL, errstr)
		return
	}

	ginutil.Response(ctx, nil, nil)
	return
}

// UserUpdate 修改
func (srv *UserService) UserUpdate(ctx *gin.Context) {
	usr := &model.User{}
	if err := ctx.ShouldBindJSON(usr); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	if has, cusr := srv.hasAdminRole(ctx); !has && cusr.ID != usr.ID {
		ginutil.Response(ctx, NOPERMISSION, nil)
		return
	}

	usr.CreatedAt = time.Time{}
	usr.Password = ""
	usr.Telephone = ""
	usr.Email = ""

	if res := srv.DB.Model(&model.User{
		Model: model.Model{
			ID: usr.ID,
		},
	}).Updates(usr); res.RowsAffected == 0 {
		errstr := ""
		if res.Error != nil {
			errstr = res.Error.Error()
		}
		ginutil.Response(ctx, UPDATE_FAIL, errstr)
		return
	}

	srv.DB.Preload("Roles").First(usr)
	ginutil.Response(ctx, nil, usr)
	return
}

// UpdateRoleRequest 修改权限
type UpdateRoleRequest struct {
	UserID uint   `json:"userid"`
	Roles  string `json:"roles"`
}

// UserUpdateRole 修改权限
func (srv *UserService) UserUpdateRole(ctx *gin.Context) {
	req := &UpdateRoleRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	if has, _ := srv.hasAdminRole(ctx); !has {
		ginutil.Response(ctx, NOPERMISSION, nil)
		return
	}

	roles := []*model.Role{}
	if err := srv.DB.Find(roles).Error; err != nil {
		ginutil.Response(ctx, ROLE_WRONG, err)
		return
	}

	usr := &model.User{
		Roles: roles,
	}
	if err := srv.DB.Model(&model.User{
		Model: model.Model{
			ID: req.UserID,
		},
	}).Updates(usr).Error; err != nil {
		ginutil.Response(ctx, UPDATE_FAIL, err)
		return
	}

	srv.DB.Preload("Roles").First(usr)
	ginutil.Response(ctx, nil, usr)
	return
}

// ChangePWDRequest 修改密码
type ChangePWDRequest struct {
	Code     string `json:"code"`
	Password string `json:"pwd"`
}

// UserChangePWD 修改密码
func (srv *UserService) UserChangePWD(ctx *gin.Context) {
	req := &ChangePWDRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	id := ginutil.GetSession(ctx, token).(uint)
	userID := strconv.FormatUint(uint64(id), 10)
	session := ginutil.GetSession(ctx, CodePWDKey+userID)
	if nil == session {
		ginutil.Response(ctx, CODE_NOT_EXIST, nil)
		return
	}

	infoMap := session.(map[string]interface{})
	if float64(time.Now().Unix()) >= infoMap["exp"].(float64) {
		ginutil.Response(ctx, CODE_EXPIRE, nil)
		return
	}

	if strings.Compare(infoMap["code"].(string), req.Code) != 0 {
		ginutil.Response(ctx, CODE_WRONG, nil)
		return
	}

	ginutil.RemoveSession(ctx, CodePWDKey+userID)

	password, err := password.CryTo(req.Password, 12, "default")
	if err != nil {

		ginutil.Response(ctx, CHANGE_PWD_FAIL, err.Error())
		return
	}

	usr := &model.User{
		Password: password,
	}
	if res := srv.DB.Model(&model.User{
		Model: model.Model{
			ID: id,
		},
	}).Updates(usr); res.RowsAffected == 0 {
		errstr := ""
		if res.Error != nil {
			errstr = res.Error.Error()
		}
		ginutil.Response(ctx, CHANGE_PWD_FAIL, errstr)
		return
	}

	srv.DB.Preload("Roles").First(usr)
	ginutil.Response(ctx, nil, usr)
	return
}

// ChangeTelRequest 修改手机号
type ChangeTelRequest struct {
	Code      string `json:"code"`
	Telephone string `json:"tel"`
}

// UserChangeTel 修改手机号
func (srv *UserService) UserChangeTel(ctx *gin.Context) {
	req := &ChangeTelRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {

		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	id := ginutil.GetSession(ctx, token).(uint)
	userID := strconv.FormatUint(uint64(id), 10)
	session := ginutil.GetSession(ctx, CodeTelKey+userID)
	if nil == session {
		ginutil.Response(ctx, CODE_NOT_EXIST, nil)
		return
	}

	infoMap := session.(map[string]interface{})
	if float64(time.Now().Unix()) >= infoMap["exp"].(float64) {
		ginutil.Response(ctx, CODE_EXPIRE, nil)
		return
	}

	if strings.Compare(infoMap["code"].(string), req.Code) != 0 {
		ginutil.Response(ctx, CODE_WRONG, nil)
		return
	}

	ginutil.RemoveSession(ctx, CodeTelKey+userID)

	usr := &model.User{
		Telephone: req.Telephone,
	}
	if res := srv.DB.Model(&model.User{
		Model: model.Model{
			ID: id,
		},
	}).Updates(usr); res.RowsAffected == 0 {
		errstr := ""
		if res.Error != nil {
			errstr = res.Error.Error()
		}
		ginutil.Response(ctx, CHANGE_TEL_FAIL, errstr)
		return
	}

	srv.DB.Preload("Roles").First(usr)
	ginutil.Response(ctx, nil, usr)
	return
}

// ChangeEmailRequest 修改邮箱
type ChangeEmailRequest struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

// UserChangeEmail 修改邮箱
func (srv *UserService) UserChangeEmail(ctx *gin.Context) {
	req := &ChangeEmailRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {

		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	id := ginutil.GetSession(ctx, token).(uint)
	userID := strconv.FormatUint(uint64(id), 10)
	session := ginutil.GetSession(ctx, CodeEmailKey+userID)
	if nil == session {
		ginutil.Response(ctx, CODE_NOT_EXIST, nil)
		return
	}

	infoMap := session.(map[string]interface{})
	if float64(time.Now().Unix()) >= infoMap["exp"].(float64) {
		ginutil.Response(ctx, CODE_EXPIRE, nil)
		return
	}

	if strings.Compare(infoMap["code"].(string), req.Code) != 0 {
		ginutil.Response(ctx, CODE_WRONG, nil)
		return
	}

	ginutil.RemoveSession(ctx, CodeEmailKey+userID)

	usr := &model.User{
		Email: req.Email,
	}
	if res := srv.DB.Model(&model.User{
		Model: model.Model{
			ID: id,
		},
	}).Updates(usr); res.RowsAffected == 0 {
		errstr := ""
		if res.Error != nil {
			errstr = res.Error.Error()
		}
		ginutil.Response(ctx, CHANGE_EMAIL_FAIL, errstr)
		return
	}

	srv.DB.Preload("Roles").First(usr)
	ginutil.Response(ctx, nil, usr)
	return
}

// CodeRequest 获取验证码
type CodeRequest struct {
	Telephone string `json:"tel"`
	Email     string `json:"email"`
	Aim       string `json:"aim"`
}

// UserLoginCode 获取验证码
func (srv *UserService) UserLoginCode(ctx *gin.Context) {
	req := &CodeRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	codesessionkey := ""
	if len(req.Email) != 0 {
		codesessionkey = CodeLoginKey + req.Telephone
	} else if len(req.Telephone) != 0 {
		codesessionkey = CodeLoginKey + req.Email
	} else if len(codesessionkey) == 0 {
		ginutil.Response(ctx, CODE_UNKOWN_TYPE, nil)
		return
	}

	session := ginutil.GetSession(ctx, CodeLoginKey)
	if nil != session {
		if info := session.(map[string]interface{}); float64(time.Now().Unix()) < info["exp"].(float64) {
			ginutil.Response(ctx, CODE_EXIST, nil)
			return
		}
	}

	info, err := srv.sendCode(req)
	if err != nil {
		ginutil.Response(ctx, err, nil)
		return
	}
	ginutil.SetSession(ctx, codesessionkey, info)
	ginutil.Response(ctx, nil, fmt.Sprintf("验证码已发送"))
	return
}

// UserChangeCode 获取验证码
func (srv *UserService) UserChangeCode(ctx *gin.Context) {
	req := &CodeRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	userID := strconv.FormatUint(uint64(ginutil.GetSession(ctx, token).(uint)), 10)

	codesessionkey := ""
	if strings.HasPrefix(CodePWDKey, req.Aim+"_") {
		codesessionkey = CodePWDKey + userID
	} else if strings.HasPrefix(CodeTelKey, req.Aim+"_") {
		codesessionkey = CodeTelKey + userID
	} else if strings.HasPrefix(CodeEmailKey, req.Aim+"_") {
		codesessionkey = CodeEmailKey + userID
	} else {
		ginutil.Response(ctx, CODE_CHANGE_AIM_INVALID, nil)
		return
	}

	session := ginutil.GetSession(ctx, codesessionkey)
	if nil != session {
		if info := session.(map[string]interface{}); float64(time.Now().Unix()) < info["exp"].(float64) {
			ginutil.Response(ctx, CODE_EXIST, nil)
			return
		}
	}

	info, err := srv.sendCode(req)
	if err != nil {
		ginutil.Response(ctx, err, nil)
		return
	}
	ginutil.SetSession(ctx, codesessionkey, info)
	ginutil.Response(ctx, nil, fmt.Sprintf("验证码已发送"))
	return
}

func (srv *UserService) sendCode(req *CodeRequest) (map[string]interface{}, error) {
	info := make(map[string]interface{})
	info["exp"] = time.Now().Add(time.Minute * 5).Unix() // 5分钟过期
	if len(req.Email) != 0 {
		info["code"] = "123456"
	} else if len(req.Telephone) != 0 {
		info["code"] = "123456"
	} else {
		return nil, CODE_UNKOWN_TYPE
	}
	return info, nil
}

// hasAdminRole 是否拥有admin
func (srv *UserService) hasAdminRole(ctx *gin.Context) (bool, *model.User) {
	token := ctx.GetHeader(headerTokenKey)
	session := ginutil.GetSession(ctx, token)
	usr := &model.User{
		Model: model.Model{
			ID: session.(uint),
		},
	}
	roles := []*model.Role{}
	srv.DB.Model(usr).Related(&roles, "Roles")
	for _, role := range roles {
		if strings.Compare(role.Key, "admin") == 0 {
			return true, usr
		}
	}
	return false, usr
}

// Register ...
func (srv *UserService) Register(api *gin.RouterGroup) {
	api.POST("/user/add", srv.UserAdd)
	api.POST("/user/login", srv.UserLogin)
	api.POST("/user/logincode", srv.UserLoginCode)
	//认证校验
	api.Use(srv.UserAuthorize)
	api.POST("/user/logout", srv.UserLogout)
	api.POST("/user/info", srv.UserInfo)
	api.POST("/user/list", srv.UserList)
	api.POST("/user/delete", srv.UserDelete)
	api.POST("/user/update", srv.UserUpdate)
	api.POST("/user/updaterole", srv.UserUpdateRole)
	api.POST("/user/changepwd", srv.UserChangePWD)
	api.POST("/user/changetel", srv.UserChangeTel)
	api.POST("/user/changeemail", srv.UserChangeEmail)
	api.POST("/user/changecode", srv.UserChangeCode)
}
