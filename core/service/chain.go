package service

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainService 区块链配置表
type ChainService struct {
	DB *gorm.DB
}

// chainValidate
type requestChainParam struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Public      bool   `json:"public"`
	OriginID    uint   `json:"originID"`
	Description string `json:"description"`
}

// 添加链
func (srv *ChainService) ChainAdd(c *gin.Context) {
	var err error

	// 获取主体信息
	var chainInfo = &requestChainParam{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateChainAdd(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 设置用户ID
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	chain := &model.Chain{
		UserID:      user.ID,
		Name:        chainInfo.Name,
		Url:         chainInfo.Url,
		Description: chainInfo.Description,
		Public:      chainInfo.Public,
		OriginID:    chainInfo.OriginID,
	}

	// 创建链
	if err = srv.DB.Create(chain).Error; nil != err {
		ginutil.Response(c, ADD_CHAIN_FAIL, err)
		return
	}

	// 建立用户与链的关联
	if err = srv.DB.Model(user).Association("OwnerChains").Append(chain).Error; nil != err {
		ginutil.Response(c, ADD_CHAIN_FAIL, err)
		return
	}

	ginutil.Response(c, nil, chain)
}

// 添加已有链（admin不可以添加别人的链）
func (srv *ChainService) ChainJoin(c *gin.Context) {
	var err error

	// 获取主体信息
	var chainInfo = &requestChainParam{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateChainID(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链是否存在
	chain := &model.Chain{Model: model.Model{ID: chainInfo.ID}}
	if err = srv.DB.First(chain).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, err)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	// 建立联系
	if err = srv.DB.Model(user).Association("OwnerChains").Append(chain).Error; nil != err {
		ginutil.Response(c, ADD_CHAIN_FAIL, err)
	}

	ginutil.Response(c, nil, chain)
}

// 获取链列表（admin可以查看任何人的链）
func (srv *ChainService) ChainList(c *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}
	var err error

	// 获取主体信息
	if err = c.ShouldBindJSON(req); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}
	offset := req.Page * req.PageSize

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	ok, user := userService.hasAdminRole(c)

	var chains []*model.Chain

	// admin查看所有链
	if ok {
		if err = srv.DB.Offset(offset).Limit(req.PageSize).Find(&chains).Error; nil != err {
			ginutil.Response(c, GET_CHAINS_FAIL, err)
			return
		}
	} else {
		// 获取与账户关联的链
		if err = srv.DB.Model(&user).Offset(offset).Limit(req.PageSize).Association("OwnerChains").Find(&chains).Error; nil != err {
			ginutil.Response(c, err, nil)
			return
		}
	}

	ginutil.Response(c, nil, chains)
}

// 删除链
func (srv *ChainService) ChainDelete(c *gin.Context) {
	var err error

	// 获取主体信息
	var chainInfo = &requestChainParam{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateChainID(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链信息
	chain := &model.Chain{Model: model.Model{ID: chainInfo.ID}}
	if err = srv.DB.First(chain).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, err)
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if chain.UserID != user.ID {
		ginutil.Response(c, NOPERMISSION, nil)
		return
	}

	// 删除chain的数据
	deleteDB := srv.DB.Unscoped().Where(chain).Delete(model.Chain{})
	if err := deleteDB.Error; nil != err {
		ginutil.Response(c, DELETE_FAIL, nil)
		return
	}

	ginutil.Response(c, nil, nil)
}

// 退出链
func (srv *ChainService) ChainExit(c *gin.Context) {
	var err error

	// 获取主体信息
	var chainInfo = &requestChainParam{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateChainID(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链信息
	chain := &model.Chain{Model: model.Model{ID: chainInfo.ID}}
	if err = srv.DB.First(chain).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, err)
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	// 如果是该用户是创建者（不可退出）
	if user.ID == chain.UserID {
		ginutil.Response(c, DELETE_FAIL, err)
		return
	}

	// 删除联系
	if err = srv.DB.Model(user).Association("OwnerChains").Delete(chain).Error; nil != err {
		ginutil.Response(c, DELETE_FAIL, err)
		return
	}

	ginutil.Response(c, nil, chain)
}

// 链更新
func (srv *ChainService) ChainUpdate(c *gin.Context) {
	// 读取主体信息
	var err error
	var chainInfo = &requestChainParam{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 校验ID是否存在
	if ok, errMsg := chainInfo.validateChainUpdate(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 获取链对应的账户ID
	chainVerify := &model.Chain{Model: model.Model{ID: chainInfo.ID}}
	if err = srv.DB.First(chainVerify).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, nil)
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	// 验证当前用户是否有修改权(admin 不可以更新)
	if chainVerify.UserID != user.ID {
		ginutil.Response(c, NOPERMISSION, nil)
		return
	}

	// 更新链
	chain := &model.Chain{Model: model.Model{ID: chainInfo.ID}}
	updateDB := srv.DB.Model(chain).Updates(
		&model.Chain{
			Name:        chainInfo.Name,
			Url:         chainInfo.Url,
			Public:      chainInfo.Public,
			Description: chainInfo.Description,
		})
	if err = updateDB.Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}

	// 获取最新链数据
	if err = srv.DB.First(chain).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, nil)
		return
	}

	// 返回更新链结果
	ginutil.Response(c, nil, chain)
}

// Register
func (srv *ChainService) Register(router *gin.Engine, api *gin.RouterGroup) {
	chain := api.Group("/chain")
	chain.POST("/add", srv.ChainAdd)
	chain.POST("/join", srv.ChainJoin)
	chain.POST("/list", srv.ChainList)
	chain.POST("/delete", srv.ChainDelete)
	chain.POST("/exit", srv.ChainExit)
	chain.POST("/update", srv.ChainUpdate)
}
