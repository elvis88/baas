package service

import (
	"errors"
	"time"

	"github.com/elvis88/baas/common/json"
	"github.com/elvis88/baas/common/jwt"
)

const (
	headerTokenKey = "Authorization"
	headerUserID   = "UserID"
	TokenKey       = "baas user secret"
	CodeLoginKey   = "login"
	CodePWDKey     = "changepwd"
	CodeTelKey     = "changetel"
	CodeEmailKey   = "changeemail"
)

var (
	REQUEST_PARAM_INVALID = errors.New("无效请求参数")
	EXEC_FAILED           = errors.New("执行失败")

	NAME_NOT_EXIST    = errors.New("用户不存在")
	ID_NOT_EXIST      = errors.New("用户ID不存在")
	TEL_NOT_EXIST     = errors.New("手机号不存在")
	EMAIL_NOT_EXIST   = errors.New("邮箱不存在")
	PASSWORD_WRONG    = errors.New("密码不正确")
	CODE_EXIST        = errors.New("验证码已存在,请勿频繁发送")
	CODE_NOT_EXIST    = errors.New("验证码不存在")
	CODE_WRONG        = errors.New("验证码不正确")
	CODE_EXPIRE       = errors.New("验证码过期")
	UNKOWN_TYPE       = errors.New("未知登陆方式")
	LOGIN_FAIL        = errors.New("登陆失败,请稍后重试")
	TOKEN_NOT_EXIST   = errors.New("未登录,请登陆")
	TOKEN_INVALID     = errors.New("验证失败")
	TOKEN_EXPIRE      = errors.New("验证过期")
	GET_FAIL          = errors.New("获取失败")
	ADD_FAIL          = errors.New("添加失败")
	DELETE_FAIL       = errors.New("删除失败")
	UPDATE_FAIL       = errors.New("更新失败")
	CHANGE_PWD_FAIL   = errors.New("修改密码失败")
	CHANGE_TEL_FAIL   = errors.New("修改手机号失败")
	CHANGE_EMAIL_FAIL = errors.New("修改邮箱失败")

	NAME_INVALID       = errors.New("无效用户名")
	EMAIL_INVALID      = errors.New("无效邮箱")
	TEL_INVALID        = errors.New("无效手机号")
	CODE_UNKOWN_TYPE   = errors.New("未知发送方式")
	CODE_SEND_FAIL     = errors.New("验证码发送失败")
	CODE_AIM_INVALID   = errors.New("无效修改验证码类型")
	NOPERMISSION       = errors.New("权限不够")
	ROLE_WRONG         = errors.New("权限不正确")
	ACTION_UNKOWN_TYPE = errors.New("未知程序文件下载")

	ADD_CHAIN_FAIL           = errors.New("添加链失败")
	GET_CHAINS_FAIL          = errors.New("获取链列表失败")
	CHAINID_NOT_EXIST        = errors.New("链不存在")
	NOT_SUPPORT_ORIGIN_CHAIN = errors.New("不支持该源链")
	CHAIN_DEPLOY_NOT_EXIST   = errors.New("节点不存在")
	CHAIN_DEPLOY_ADD_FAIL    = errors.New("链实例添加失败")
	PARAMS_IS_NOT_ENOUGH     = errors.New("参数不足")
)

// AuthorizeToken login jwt
type AuthorizeToken struct {
	UserID uint  `json:"userID"`
	Exp    int64 `json:"exp"`
	Iat    int64 `json:"iat"`
}

func newAuthorizeToken(UserID uint) *AuthorizeToken {
	token := &AuthorizeToken{}
	token.UserID = UserID
	now := time.Now()
	token.Exp = now.Add(time.Hour * 1).Unix()
	token.Iat = now.Unix()
	return token
}

func newFromJWT(jwttoken string) *AuthorizeToken {
	info, ok := jwt.ParseToken(jwttoken, TokenKey)
	if !ok {
		return nil
	}
	bts, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	token := &AuthorizeToken{}
	err = json.Unmarshal(bts, token)
	if err != nil {
		panic(err)
	}
	return token
}

func (a *AuthorizeToken) toJWT() (string, error) {
	bts, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	info := make(map[string]interface{})
	err = json.Unmarshal(bts, &info)
	if err != nil {
		panic(err)
	}
	return jwt.CreateToken(TokenKey, info)
}
