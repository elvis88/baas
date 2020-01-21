package service

import (
	"fmt"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainService 区块链配置表
type ChainService struct {
	DB *gorm.DB
}

// AncestorChain amin创建的链
func (srv *ChainService) getAncestorChain(ancestorID uint) (*model.Chain, error) {
	ancestorChain := &model.Chain{}
	if err := srv.DB.First(ancestorChain, ancestorID).Error; err != nil {
		return nil, err
	}

	user := &model.User{}
	if err := srv.DB.First(user, ancestorChain.UserID).Error; err != nil {
		return nil, err
	}

	if user.Name != "admin" {
		return nil, fmt.Errorf("ancestor id was wrong")
	}
	return ancestorChain, nil
}

// JoinChain public=true的链
func (srv *ChainService) getJoinChain(joinID uint) (*model.Chain, error) {
	JoinChain := &model.Chain{}
	if err := srv.DB.First(JoinChain, joinID).Error; err != nil {
		return nil, err
	}

	if JoinChain.Public != true {
		return nil, fmt.Errorf("join id was wrong")
	}
	return JoinChain, nil
}

// AncestorList 获取来源链列表
func (srv *ChainService) AncestorList(c *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}
	if err := c.ShouldBindJSON(req); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	user := &model.User{
		Name: "admin",
	}
	if err := srv.DB.First(user).Error; err != nil {
		panic(err.Error())
	}

	var chains []*model.Chain
	offset := req.Page * req.PageSize
	if err := srv.DB.Where(&model.Chain{
		UserID: user.ID,
	}).Offset(offset).Limit(req.PageSize).Find(&chains).Error; err != nil {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, chains)
}

// JoinList 获取加入链列表
func (srv *ChainService) JoinList(c *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}
	if err := c.ShouldBindJSON(req); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	var chains []*model.Chain
	offset := req.Page * req.PageSize
	if err := srv.DB.Not(&model.Chain{
		UserID: user.ID,
	}).Where(&model.Chain{
		Public: true,
	}).Offset(offset).Limit(req.PageSize).Find(&chains).Error; nil != err {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, chains)
}

// List 获取链列表
func (srv *ChainService) List(c *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}
	if err := c.ShouldBindJSON(req); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	var chains []*model.Chain
	offset := req.Page * req.PageSize
	// 获取与账户关联的链
	if err := srv.DB.Model(&user).Offset(offset).Limit(req.PageSize).Association("OwnerChains").Find(&chains).Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}
	ginutil.Response(c, nil, chains)
}

// Add 添加链
func (srv *ChainService) Add(c *gin.Context) {
	var chainInfo = &requestChainParam{}
	if err := c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateAdd(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 校验 UserID
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if user.Name == "admin" {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "admin")
		return
	}

	chain := &model.Chain{
		UserID:      user.ID,
		Name:        chainInfo.Name,
		URL:         chainInfo.URL,
		Description: chainInfo.Description,
		Public:      chainInfo.Public,
		AncestorID:  chainInfo.AncestorID,
	}

	// 创建链
	if err := chain.Add(srv.DB); nil != err {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, chain)
}

// Delete 删除链
func (srv *ChainService) Delete(c *gin.Context) {
	var chainInfo = &requestChainParam{}
	if err := c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateID(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链信息
	chain := &model.Chain{}
	if err := srv.DB.First(chain, chainInfo.ID).Error; nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if user.Name == "admin" {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "admin")
		return
	}
	// 验证当前用户是否有修改权(admin 不可以删除)
	if chain.UserID != user.ID {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "not self")
		return
	}

	// 删除链
	if err := chain.Remove(srv.DB); nil != err {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, nil)
}

// Join 添加已有链（admin不可以添加别人的链）
func (srv *ChainService) Join(c *gin.Context) {
	var chainInfo = &requestChainParam{}
	if err := c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateID(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链是否存在
	joinChain, err := srv.getJoinChain(chainInfo.ID)
	if nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if user.Name == "admin" {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "admin")
		return
	}
	if user.ID == joinChain.UserID {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "join self")
		return
	}

	// 已加入
	if srv.DB.Model(&model.User{
		Model: gorm.Model{
			ID: user.ID,
		},
	}).Where(&model.Chain{
		Model: gorm.Model{
			ID: chainInfo.ID,
		},
	}).Association("OwnerChains").Count() != 0 {
		ginutil.Response(c, EXEC_FAILED, "joined")
		return
	}

	// 加入链
	if err := joinChain.Join(srv.DB, user); nil != err {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, joinChain)
}

// Unjoin 退出加入链
func (srv *ChainService) Unjoin(c *gin.Context) {
	var chainInfo = &requestChainParam{}
	if err := c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验参数
	if ok, errMsg := chainInfo.validateID(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链信息
	joinChain, err := srv.getJoinChain(chainInfo.ID)
	if nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if user.Name == "admin" {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "admin")
		return
	}
	if user.ID == joinChain.UserID {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "unjoin self")
		return
	}

	if err := joinChain.Unjoin(srv.DB, user); nil != err {
		ginutil.Response(c, DELETE_FAIL, err.Error())
		return
	}
	ginutil.Response(c, nil, nil)
}

// GetConfig 获得链config内容
func (srv *ChainService) GetConfig(c *gin.Context) {
	var chainConfig = &requestChainConfig{}
	if err := c.ShouldBindJSON(chainConfig); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验参数
	if ok, errMsg := chainConfig.validateGetConfig(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链信息
	chain := &model.Chain{}
	if err := srv.DB.First(chain, chainConfig.ID).Error; nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if chain.UserID != user.ID {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "not self")
		return
	}

	spec, err := chain.Spec(srv.DB)
	if err != nil {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	config, err := spec.GetConfig()
	if err != nil {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, config)
}

// SetConfig 用户修改链config内容
func (srv *ChainService) SetConfig(c *gin.Context) {
	var chainConfig = &requestChainConfig{}
	if err := c.ShouldBindJSON(chainConfig); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验参数
	if ok, errMsg := chainConfig.validateSetConfig(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链信息
	chain := &model.Chain{}
	if err := srv.DB.First(chain, chainConfig.ID).Error; nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if chain.UserID != user.ID {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "not self")
		return
	}

	spec, err := chain.Spec(srv.DB)
	if err != nil {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}

	if err := spec.SetConfig(chainConfig.Config); err != nil {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, nil)
}

// Register ...
func (srv *ChainService) Register(router *gin.Engine, api *gin.RouterGroup) {
	chain := api.Group("/chain")
	chain.POST("/ancestorlist", srv.AncestorList)
	chain.POST("/joinlist", srv.JoinList)
	chain.POST("/list", srv.List)
	chain.POST("/add", srv.Add)
	chain.POST("/remove", srv.Delete)
	chain.POST("/join", srv.Join)
	chain.POST("/unjoin", srv.Unjoin)
	chain.POST("/getconfig", srv.GetConfig)
	chain.POST("/setconfig", srv.SetConfig)
}
