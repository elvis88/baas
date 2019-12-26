package service

import (
	"errors"
)

const (
	headerTokenKey = "Authorization"
	TokenKey       = "baas user secret"
	CodeLoginKey   = "login_"
	CodePWDKey     = "pwd_"
	CodeTelKey     = "tel_"
	CodeEmailKey   = "email_"
)

var (
	REQUEST_PARAM_INVALID = errors.New("无效请求参数")
	NAME_NOT_EXIST        = errors.New("用户不存在")
	ID_NOT_EXIST          = errors.New("用户ID不存在")
	TEL_NOT_EXIST         = errors.New("手机号不存在")
	EMAIL_NOT_EXIST       = errors.New("邮箱不存在")
	PASSWORD_WRONG        = errors.New("密码不正确")
	CODE_EXIST            = errors.New("验证码已存在,请勿频繁发送")
	CODE_NOT_EXIST        = errors.New("验证码不存在")
	CODE_WRONG            = errors.New("验证码不正确")
	CODE_EXPIRE           = errors.New("验证码过期")
	UNKOWN_TYPE           = errors.New("未知登陆方式")
	LOGIN_FAIL            = errors.New("登陆失败,请稍后重试")
	TOKEN_NOT_EXIST       = errors.New("未登录,请登陆")
	TOKEN_INVALID         = errors.New("验证失败")
	TOKEN_EXPIRE          = errors.New("验证过期")
	GET_FAIL              = errors.New("获取失败")
	ADD_FAIL              = errors.New("添加失败")
	DELETE_FAIL           = errors.New("删除失败")
	UPDATE_FAIL           = errors.New("更新失败")
	CHANGE_PWD_FAIL       = errors.New("修改密码失败")
	CHANGE_TEL_FAIL       = errors.New("修改手机号失败")
	CHANGE_EMAIL_FAIL     = errors.New("修改邮箱失败")

	NAME_INVALID            = errors.New("无效用户名")
	EMAIL_INVALID           = errors.New("无效邮箱")
	TEL_INVALID             = errors.New("无效手机号")
	CODE_UNKOWN_TYPE        = errors.New("未知发送方式")
	CODE_CHANGE_AIM_INVALID = errors.New("无效修改验证码类型")
	NOPERMISSION            = errors.New("权限不够")
	ROLE_WRONG              = errors.New("权限不正确")

	ADD_CHAIN_FAIL			= errors.New("添加链失败")
	CHAINID_NOT_EXIST   	= errors.New("链不存在")
	CHAINID_DEPLOY_NOT_EXIST = errors.New("节点不存在")
	PARAMS_IS_NOT_ENOUGH 	= errors.New("参数不足")

	ADD_CHAIN_DEPLOY_FAIL   = errors.New("添加链实例失败")
)
