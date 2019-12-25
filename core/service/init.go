package service

import (
	"errors"

	"github.com/elvis88/baas/common/log"
)

var logger = log.GetLogger("baas", log.DEBUG)

const headerTokenKey = "Authorization"
const TokenKey = "baas user secret"

var (
	REQUEST_PARAM_INVALID = errors.New("无效请求参数")
	NAME_NOT_EXIST        = errors.New("用户不存在")
	ID_NOT_EXIST          = errors.New("用户ID不存在")
	TEL_NOT_EXIST         = errors.New("手机号不存在")
	EMAIL_NOT_EXIST       = errors.New("邮箱不存在")
	PASSWORD_WRONG        = errors.New("密码不正确")
	CODE_WRONG            = errors.New("验证码不正确")
	UNKOWN_TYPE           = errors.New("未知登陆方式")
	LOGIN_FAIL            = errors.New("登陆失败,请稍后重试")
	TOKEN_NOT_EXIST       = errors.New("Token不存在")
	TOKEN_INVALID         = errors.New("Token无效")
	TOKEN_EXPIRE          = errors.New("Token过期")
	GET_FAIL              = errors.New("获取失败")
	ADD_FAIL              = errors.New("添加失败")
	DELETE_FAIL           = errors.New("删除失败")
	UPDATE_FAIL           = errors.New("更新失败")
)
