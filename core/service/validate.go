package service

import "github.com/thedevsaddam/govalidator"

func executeValidation(data interface{}, rules, messages govalidator.MapData) (bool, *string) {
	opts := govalidator.Options{
		Data:     data,
		Rules:    rules,
		Messages: messages,
	}
	e := govalidator.New(opts).ValidateStruct()
	if len(e) > 0 {
		for _, v := range e {
			return false, &v[0]
		}
	}

	return true, nil
}

// User api params validate.
// 登录字段参数验证
func (r *UserRequest) validateLogin() (bool, *string) {
	rules := govalidator.MapData{
		"name":  []string{"alpha_num"},
		"pwd":   []string{"between:6,20"},
		"phone": []string{"digits:11"},
		"email": []string{"email"},
		"code":  []string{"alpha_num"},
	}

	messages := govalidator.MapData{
		"name":  []string{"alpha_num:用户名必须是数字/字母，或者数字字母组合"},
		"pwd":   []string{"between:密码长度只能为6至20位"},
		"phone": []string{"digits:电话号码必须是11位"},
		"email": []string{"email:邮箱格式不正确"},
		"code":  []string{"alpha_num:用户名必须是数字/字母，或者数字字母组合"},
	}

	return executeValidation(r, rules, messages)
}

// 新增验证
func (r *UserRequest) validateAdd() (bool, *string) {
	rules := govalidator.MapData{
		"name":  []string{"alpha_num"},
		"pwd":   []string{"between:6,20"},
		"phone": []string{"digits:11"},
		"email": []string{"email"},
	}

	messages := govalidator.MapData{
		"name":  []string{"alpha_num:用户名必须是数字/字母，或者数字字母组合"},
		"pwd":   []string{"between:密码长度只能为6至20位"},
		"phone": []string{"digits:电话号码必须是11位"},
		"email": []string{"email:邮箱格式不正确"},
	}

	return executeValidation(r, rules, messages)
}

// 用户更新角色参数验证
func (r *UpdateRoleRequest) validateUpdateRole() (bool, *string) {
	rules := govalidator.MapData{
		"userid": []string{"required"},
		"roles":  []string{"required"},
	}

	messages := govalidator.MapData{
		"userid": []string{"required:用户ID不能为空"},
		"roles":  []string{"required:角色不能为空"},
	}

	return executeValidation(r, rules, messages)
}

// 修改密码参数验证
func (r *ChangePWDRequest) validateUpdatePwd() (bool, *string) {
	rules := govalidator.MapData{
		"code": []string{"required, alpha_num"},
		"pwd":  []string{"required,between:6,20"},
	}

	messages := govalidator.MapData{
		"code": []string{"required:验证码不能为空,alpha_num:用户名必须是数字/字母，或者数字字母组合"},
		"pwd":  []string{"required:新密码不能为空,between:密码长度只能为6至20位"},
	}

	return executeValidation(r, rules, messages)
}

// 修改手机号验证
func (r *ChangeTelRequest) validateUpdateTel() (bool, *string) {
	rules := govalidator.MapData{
		"code":  []string{"required, alpha_num"},
		"phone": []string{"required, digits:11"},
	}

	messages := govalidator.MapData{
		"code": []string{"required:验证码不能为空, alpha_num:用户名必须是数字/字母，或者数字字母组合"},
		"pwd":  []string{"required:新密码不能为空, digits:电话必须为11位数字"},
	}

	return executeValidation(r, rules, messages)
}

// 修改邮箱参数验证
func (r *ChangeEmailRequest) validateUpdateEmail() (bool, *string) {
	rules := govalidator.MapData{
		"code":  []string{"required, alpha_num"},
		"email": []string{"required, email"},
	}

	messages := govalidator.MapData{
		"code":  []string{"required:验证码不能为空,alpha_num:用户名必须是数字/字母，或者数字字母组合"},
		"email": []string{"required:新密码不能为空, email:邮箱不合法"},
	}

	return executeValidation(r, rules, messages)
}

// 获取验证码参数验证
func (r *CodeRequest) validateGetCode() (bool, *string) {
	rules := govalidator.MapData{
		"phone": []string{"digits:11"},
		"email": []string{"email"},
		"aim":   []string{"required"},
	}

	messages := govalidator.MapData{
		"phone": []string{"digits:电话格式不合法（需要11位数字）"},
		"email": []string{"email:邮箱格式不合法"},
		"aim":   []string{"required:验证码类型不能为空"},
	}

	return executeValidation(r, rules, messages)
}

// Chain api params validate.
func (r *requestChainParam) validateChainAdd() (bool, *string) {
	rules := govalidator.MapData{
		"name":     []string{"required"},
		"url":      []string{"url"},
		"originID": []string{"required"},
	}

	messages := govalidator.MapData{
		"name":     []string{"required:名字不能为空"},
		"url":      []string{"url:Url格式不正确"},
		"originID": []string{"required:必须指定链来源"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainParam) validateChainID() (bool, *string) {
	rules := govalidator.MapData{
		"id": []string{"required"},
	}

	messages := govalidator.MapData{
		"id": []string{"required:链ID不能为空"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainParam) validateChainUpdate() (bool, *string) {
	rules := govalidator.MapData{
		"id":   []string{"required"},
		"name": []string{"required"},
		"url":  []string{"url"},
	}

	messages := govalidator.MapData{
		"id":   []string{"required:链ID不能为空"},
		"name": []string{"required:名字不能为空"},
		"url":  []string{"url:Url格式不正确"},
	}
	return executeValidation(r, rules, messages)
}

func (r *requestChainConfig) validateGetConfig() (bool, *string) {
	rules := govalidator.MapData{
		"id": []string{"required"},
	}

	messages := govalidator.MapData{
		"id": []string{"required:链ID不能为空"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainConfig) validateSetConfig() (bool, *string) {
	rules := govalidator.MapData{
		"id":     []string{"required"},
		"config": []string{"required"},
	}

	messages := govalidator.MapData{
		"id":     []string{"required:链ID不能为空"},
		"config": []string{"required:config文件不能为空"},
	}

	return executeValidation(r, rules, messages)
}


// Chain deploy api params validate.
func (r *requestChainDeployParams) validateChainDeployAdd() (bool, *string) {
	rules := govalidator.MapData{
		"name":    []string{"required"},
		"chainID": []string{"required"},
	}

	messages := govalidator.MapData{
		"name":    []string{"required: 节点名不能为空"},
		"chainID": []string{"required: 链ID不能为空"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainDeployParams) validateChainDeployID() (bool, *string) {
	rules := govalidator.MapData{
		"id": []string{"required"},
	}

	messages := govalidator.MapData{
		"id": []string{"required:实例ID不能为空"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainDeployParams) validateChainDeployUpdate() (bool, *string) {
	rules := govalidator.MapData{
		"id":      []string{"required"},
		"name":    []string{"required"},
		"chainID": []string{"required"},
	}

	messages := govalidator.MapData{
		"id":      []string{"required:实例ID不能为空"},
		"name":    []string{"required: 节点名不能为空"},
		"chainID": []string{"required: 链ID不能为空"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainDeployConfig) validateGetFile() (bool, *string) {
	rules := govalidator.MapData{
		"id": []string{"required"},
	}

	messages := govalidator.MapData{
		"id": []string{"required:链实例ID不能为空"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainDeployConfig) validateSetConfig() (bool, *string) {
	rules := govalidator.MapData{
		"id":     []string{"required"},
		"config": []string{"required"},
	}

	messages := govalidator.MapData{
		"id":     []string{"required:链ID不能为空"},
		"config": []string{"required:config文件不能为空"},
	}

	return executeValidation(r, rules, messages)
}
