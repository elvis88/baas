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

type requestChainParam struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Public      bool   `json:"public"`
	OriginID    uint   `json:"originID"`
	Description string `json:"description"`
}

func (r *requestChainParam) validateChainAdd() (bool, *string) {
	rules := govalidator.MapData{
		"name":     []string{"required"},
		"url":      []string{"url"},
		"originID": []string{"required, numeric_between:1,"},
	}

	messages := govalidator.MapData{
		"name":     []string{"required:名字不能为空"},
		"url":      []string{"url:Url格式不正确"},
		"originID": []string{"required:必须指定链来源, numeric_between:来源链ID无效"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainParam) validateChainID() (bool, *string) {
	rules := govalidator.MapData{
		"id": []string{"required, numeric_between:1,"},
	}

	messages := govalidator.MapData{
		"id": []string{"required:链ID不能为空, numeric_between:链ID无效"},
	}

	return executeValidation(r, rules, messages)
}

func (r *requestChainParam) validateChainUpdate() (bool, *string) {
	rules := govalidator.MapData{
		"id":   []string{"required, numeric_between:1,"},
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