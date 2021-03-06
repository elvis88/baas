package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elvis88/baas/common/sms"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/common/password"
	"github.com/elvis88/baas/core/generate"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserService 用户表
type UserService struct {
	DB *gorm.DB
}

// LoginRequest 登陆请求参数
type UserRequest struct {
	ID        string `json:"id"`
	UserName  string `json:"name"`
	Password  string `json:"pwd"`
	Nick      string `json:"nick"`
	Telephone string `json:"phone"`
	Email     string `json:"email"`
	Code      string `json:"code"`
}

// UserLogin 登陆
func (srv *UserService) UserLogin(ctx *gin.Context) {
	login := &UserRequest{}
	if err := ctx.ShouldBindJSON(login); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := login.validateLogin(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
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
			codesessionkey = CodeLoginKey + login.Email
		} else if len(login.Telephone) != 0 { // 手机验证码登陆
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
		info := map[string]interface{}{}
		json.Unmarshal(session.([]byte), &info)
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

		if len(login.Email) != 0 { // 邮箱验证码
			usr.Name = sms.GetRandomString(6)
			usr.Email = login.Email
			if err := srv.DB.FirstOrCreate(usr, &model.User{
				Email: login.Email,
			}).Error; err != nil {
				ginutil.Response(ctx, LOGIN_FAIL, err.Error())
				return
			}
		} else if len(login.Telephone) != 0 { // 手机验证码登陆
			usr.Name = sms.GetRandomString(6)
			usr.Telephone = login.Telephone
			if err := srv.DB.FirstOrCreate(usr, &model.User{
				Telephone: login.Telephone,
			}).Error; err != nil {
				ginutil.Response(ctx, LOGIN_FAIL, err.Error())
				return
			}
		}
	} else {
		ginutil.Response(ctx, UNKOWN_TYPE, nil)
		return
	}

	token := newAuthorizeToken(usr.ID)
	jwttoken, err := token.toJWT()
	if err != nil {
		ginutil.Response(ctx, LOGIN_FAIL, err.Error())
		return
	}

	ctx.Header(headerTokenKey, jwttoken)
	ginutil.Response(ctx, nil, usr)
	return
}

// UserAuthorize 用户验证
func (srv *UserService) UserAuthorize(ctx *gin.Context) {
	jwttoken := ctx.GetHeader(headerTokenKey)

	token := newFromJWT(jwttoken)
	if token == nil {
		ginutil.Response(ctx, TOKEN_INVALID, nil)
		ctx.Abort()
		return
	}

	if time.Now().Unix() >= token.Exp {
		ginutil.Response(ctx, TOKEN_EXPIRE, nil)
		ctx.Abort()
		return
	}

	usr := &model.User{}
	if err := srv.DB.Where(&model.User{
		Model: gorm.Model{
			ID: token.UserID,
		},
	}).First(usr).Error; err != nil {
		ginutil.Response(ctx, ID_NOT_EXIST, err.Error())
		ctx.Abort()
		return
	}
	ginutil.SetSession(ctx, jwttoken, usr.ID)
	ctx.Next()
}

// UserLogout 退出登陆
func (srv *UserService) UserLogout(ctx *gin.Context) {
	ctx.Header(headerTokenKey, "")
	ginutil.Response(ctx, nil, nil)
}

// UserInfo 用户信息
func (srv *UserService) UserInfo(ctx *gin.Context) {
	token := ctx.GetHeader(headerTokenKey)
	session := ginutil.GetSession(ctx, token)
	usr := &model.User{}
	srv.DB.Where(&model.User{
		Model: gorm.Model{
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
		if err := srv.DB.Where(cusr).Offset(offset).Limit(req.PageSize).Find(&usrs).Error; err != nil {
			ginutil.Response(ctx, GET_FAIL, err.Error())
			return
		}
	} else {
		if err := srv.DB.Offset(offset).Limit(req.PageSize).Find(&usrs).Error; err != nil {
			ginutil.Response(ctx, GET_FAIL, err.Error())
			return
		}
	}

	ginutil.Response(ctx, nil, usrs)
	return
}

// UserAdd 新增用户
func (srv *UserService) UserAdd(ctx *gin.Context) {
	userRequest := &UserRequest{}
	if err := ctx.ShouldBindJSON(userRequest); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验用户新增参数
	if ok, errMsg := userRequest.validateAdd(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 用户密码加密
	password, err := password.CryTo(userRequest.Password, 12, "default")
	if err != nil {
		ginutil.Response(ctx, ADD_FAIL, err.Error())
		return
	}

	// 构建新增结构体
	usr := &model.User{
		Name:      userRequest.UserName,
		Password:  password,
		Nick:      userRequest.Nick,
		Telephone: userRequest.Telephone,
		Email:     userRequest.Email,
	}

	// 新建用户
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
		Model: gorm.Model{
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

	if ok, errMsg := req.validateUpdateRole(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	if has, _ := srv.hasAdminRole(ctx); !has {
		ginutil.Response(ctx, NOPERMISSION, nil)
		return
	}

	usr := &model.User{}
	if err := srv.DB.Model(&model.User{
		Model: gorm.Model{
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

	// 验证参数合法性
	if ok, errMsg := req.validateUpdatePwd(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
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

	info := map[string]interface{}{}
	json.Unmarshal(session.([]byte), &info)
	if float64(time.Now().Unix()) >= info["exp"].(float64) {
		ginutil.Response(ctx, CODE_EXPIRE, nil)
		return
	}

	if strings.Compare(info["code"].(string), req.Code) != 0 {
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
		Model: gorm.Model{
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
	Telephone string `json:"phone"`
}

// UserChangeTel 修改手机号
func (srv *UserService) UserChangeTel(ctx *gin.Context) {
	req := &ChangeTelRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数合法性
	if ok, errMsg := req.validateUpdateTel(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
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

	info := map[string]interface{}{}
	json.Unmarshal(session.([]byte), &info)
	if float64(time.Now().Unix()) >= info["exp"].(float64) {
		ginutil.Response(ctx, CODE_EXPIRE, nil)
		return
	}

	if strings.Compare(info["code"].(string), req.Code) != 0 {
		ginutil.Response(ctx, CODE_WRONG, nil)
		return
	}

	ginutil.RemoveSession(ctx, CodeTelKey+userID)

	usr := &model.User{
		Telephone: req.Telephone,
	}
	if res := srv.DB.Model(&model.User{
		Model: gorm.Model{
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

	// 验证参数合法性
	if ok, errMsg := req.validateUpdateEmail(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
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

	info := map[string]interface{}{}
	json.Unmarshal(session.([]byte), &info)
	if float64(time.Now().Unix()) >= info["exp"].(float64) {
		ginutil.Response(ctx, CODE_EXPIRE, nil)
		return
	}

	if strings.Compare(info["code"].(string), req.Code) != 0 {
		ginutil.Response(ctx, CODE_WRONG, nil)
		return
	}

	ginutil.RemoveSession(ctx, CodeEmailKey+userID)

	usr := &model.User{
		Email: req.Email,
	}
	if res := srv.DB.Model(&model.User{
		Model: gorm.Model{
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
	Telephone string `json:"phone"`
	Email     string `json:"email"`
	Aim       string `json:"aim" binding:"requried"`
}

// UserLoginCode 获取验证码
func (srv *UserService) UserLoginCode(ctx *gin.Context) {
	req := &CodeRequest{
		Aim: CodeLoginKey,
	}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数合法性
	if ok, errMsg := req.validateGetCode(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	if strings.Compare(req.Aim, CodeLoginKey) != 0 {
		ginutil.Response(ctx, CODE_AIM_INVALID, nil)
		return
	}

	codesessionkey := ""
	if len(req.Email) != 0 {
		codesessionkey = CodeLoginKey + req.Email
	} else if len(req.Telephone) != 0 {
		codesessionkey = CodeLoginKey + req.Telephone
	} else if len(codesessionkey) == 0 {
		ginutil.Response(ctx, CODE_UNKOWN_TYPE, nil)
		return
	}

	session := ginutil.GetSession(ctx, CodeLoginKey)
	if nil != session {
		info := map[string]interface{}{}
		json.Unmarshal(session.([]byte), &info)
		if float64(time.Now().Unix()) < info["exp"].(float64) {
			ginutil.Response(ctx, CODE_EXIST, nil)
			return
		}
	}

	info, err := srv.sendCode(req)
	if err != nil {
		ginutil.Response(ctx, err, nil)
		return
	}
	bts, _ := json.Marshal(info)
	ginutil.SetSession(ctx, codesessionkey, bts)
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

	// 验证参数合法性
	if ok, errMsg := req.validateGetCode(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	userID := strconv.FormatUint(uint64(ginutil.GetSession(ctx, token).(uint)), 10)

	codesessionkey := ""
	if strings.Compare(CodePWDKey, req.Aim) == 0 {
		codesessionkey = CodePWDKey + userID
	} else if strings.Compare(CodeTelKey, req.Aim) == 0 {
		codesessionkey = CodeTelKey + userID
	} else if strings.Compare(CodeEmailKey, req.Aim) == 0 {
		codesessionkey = CodeEmailKey + userID
	} else {
		ginutil.Response(ctx, CODE_AIM_INVALID, nil)
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
	code := sms.GetRandomString(6)
	if len(req.Email) != 0 {
		if err := emailClient.Send(&sms.Message{
			Subject: req.Aim,
			To:      []string{req.Email},
			Content: bytes.NewBufferString(fmt.Sprintf(mailBody, req.Aim, code)),
		}, false); err != nil {
			return nil, CODE_SEND_FAIL
		}
	} else if len(req.Telephone) != 0 {
		code = "123456"
	} else {
		return nil, CODE_UNKOWN_TYPE
	}
	info := make(map[string]interface{})
	info["exp"] = time.Now().Add(time.Minute * 5).Unix() // 5分钟过期
	info["code"] = code
	return info, nil
}

// hasAdminRole 是否拥有admin
func (srv *UserService) hasAdminRole(ctx *gin.Context) (bool, *model.User) {
	jwttoken := ctx.GetHeader(headerTokenKey)
	session := ginutil.GetSession(ctx, jwttoken)

	usr := &model.User{}
	if err := srv.DB.First(usr, session.(uint)).Error; nil != err {
		return false, nil
	}

	return false, usr
}

// UserGetFile 获取文件
func (srv *UserService) UserGetFile(ctx *gin.Context) (sysErr, err error) {
	nodename := ctx.Param("nodename")
	action := ctx.Param("action")

	_, cusr := srv.hasAdminRole(ctx)
	cnt := 0
	// 判断 nodename 是否是自己创建
	if strings.Compare(action, generate.Application) == 0 {
		// 获取链信息
		chain := &model.Chain{}
		if err := srv.DB.Model(&model.Chain{}).Where(&model.Chain{
			Name: nodename,
		}).Find(chain).Count(&cnt).Error; err != nil {
			return GET_FAIL, err
		}
		// 查询用户是否与链有关联
		if ok, err := (&ChainDeployService{srv.DB}).userHaveChain(cusr.ID, chain.ID); !ok {
			return NOPERMISSION, err
		}
	} else if strings.Compare(action, generate.Deployment) == 0 {
		if err := srv.DB.Model(&model.ChainDeploy{}).Where(&model.ChainDeploy{
			Name:   nodename,
			UserID: cusr.ID,
		}).Count(&cnt).Error; err != nil {
			return GET_FAIL, err
		}
	} else {
		return ACTION_UNKOWN_TYPE, nil
	}

	if cnt == 0 {
		ginutil.Response(ctx, NOPERMISSION, nil)
		return NOPERMISSION, nil
	}
	usrName := cusr.Name
	ctx.Request.URL.Path = strings.Replace(ctx.Request.URL.Path, fmt.Sprintf("file/%s/%s", action, nodename), fmt.Sprintf("data/%s/%s/%s", usrName, action, nodename), 1)
	return nil, nil
}

// Register ...
func (srv *UserService) Register(router *gin.Engine, api *gin.RouterGroup) {
	api.POST("/user/add", srv.UserAdd)
	api.POST("/user/login", srv.UserLogin)
	//api.POST("/user/logincode", srv.UserLoginCode)
	//认证校验
	api.Use(srv.UserAuthorize)
	// api.POST("/user/logout", srv.UserLogout)
	api.POST("/user/info", srv.UserInfo)
	api.POST("/user/list", srv.UserList)
	// api.POST("/user/delete", srv.UserDelete)
	// api.POST("/user/update", srv.UserUpdate)
	// api.POST("/user/updaterole", srv.UserUpdateRole)
	// api.POST("/user/changepwd", srv.UserChangePWD)
	// api.POST("/user/changetel", srv.UserChangeTel)
	// api.POST("/user/changeemail", srv.UserChangeEmail)
	// api.POST("/user/changecode", srv.UserChangeCode)

	// 脚本下载
	api.Static("/data", "./shared/data")
	api.Static("/agent", "./shared/agent")
	api.GET("/file/:action/:nodename/:fname", func(ctx *gin.Context) {
		if sysErr, err := srv.UserGetFile(ctx); nil != sysErr {
			ginutil.Response(ctx, sysErr, err.Error())
			return
		}
		router.HandleContext(ctx)
	})
}
