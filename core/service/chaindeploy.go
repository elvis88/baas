package service

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainDeployService 区块链配置表
type ChainDeployService struct {
	DB *gorm.DB
}

func (srv *ChainDeployService) userHaveChain(userID, chainID uint) (b bool, err error) {
	var chain model.Chain

	// 查看用户是否和该链有关联
	if err = srv.DB.Model(&model.User{Model: gorm.Model{ID: userID}}).
		Where(&model.Chain{Model: gorm.Model{ID: chainID}}).
		Association("OwnerChains").Find(&chain).Error; nil != err || chain.ID == 0 {
		return false, err
	}

	return true, nil
}

// List 列表
func (srv *ChainDeployService) List(ctx *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}
	if err := ctx.ShouldBindJSON(req); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	var chainDeploys []*model.ChainDeploy
	offset := req.Page * req.PageSize
	if err := srv.DB.Where(&model.ChainDeploy{
		UserID: user.ID,
	}).Offset(offset).Limit(req.PageSize).Find(&chainDeploys).Error; nil != err {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, chainDeploys)
}

// Add 新增
func (srv *ChainDeployService) Add(ctx *gin.Context) {
	chainDeployParams := &requestChainDeployParams{}
	if err := ctx.ShouldBindJSON(chainDeployParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := chainDeployParams.validateAdd(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证用户是否拥有链
	ok, err := srv.userHaveChain(user.ID, chainDeployParams.ChainID)
	if nil != err || !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 添加链实例
	chainDeploy := &model.ChainDeploy{
		Name:        chainDeployParams.Name,
		UserID:      user.ID,
		ChainID:     chainDeployParams.ChainID,
		Description: chainDeployParams.Description,
	}

	if err = chainDeploy.Add(srv.DB); nil != err {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, chainDeploy)
}

// Delete 删除
func (srv *ChainDeployService) Delete(ctx *gin.Context) {
	chainDeployParams := &requestChainDeployParams{}
	if err := ctx.ShouldBindJSON(chainDeployParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := chainDeployParams.validateChainDeployID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 指定id获取实例数据
	chainDeploy := &model.ChainDeploy{}
	if err := srv.DB.First(chainDeploy, chainDeployParams.ID).Error; nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if chainDeploy.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	if err := chainDeploy.Remove(srv.DB); err != nil {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, nil)
}

// GetConfig 用户获得config内容
func (srv *ChainDeployService) GetConfig(c *gin.Context) {
	var chainDeployConfig = &requestChainDeployConfig{}
	if err := c.ShouldBindJSON(chainDeployConfig); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验参数
	if ok, errMsg := chainDeployConfig.validateGetFile(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链实例信息
	chainDeploy := &model.ChainDeploy{}
	if err := srv.DB.First(chainDeploy, chainDeployConfig.ID).Error; nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if chainDeploy.UserID != user.ID {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "not self")
		return
	}

	spec, err := chainDeploy.Spec(srv.DB)
	if nil != err {
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

// SetConfig 用户修改config内容
func (srv *ChainDeployService) SetConfig(c *gin.Context) {
	var chainDeployConfig = &requestChainDeployConfig{}
	if err := c.ShouldBindJSON(chainDeployConfig); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 校验参数
	if ok, errMsg := chainDeployConfig.validateSetConfig(); !ok {
		ginutil.Response(c, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 查询链实例信息
	chainDeploy := &model.ChainDeploy{}
	if err := srv.DB.First(chainDeploy, chainDeployConfig.ID).Error; nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if chainDeploy.UserID != user.ID {
		ginutil.Response(c, REQUEST_PARAM_INVALID, "not self")
		return
	}

	spec, err := chainDeploy.Spec(srv.DB)
	if nil != err {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	err = spec.SetConfig(chainDeployConfig.Config)
	if err != nil {
		ginutil.Response(c, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(c, nil, nil)
}

// Register ...
func (srv *ChainDeployService) Register(router *gin.Engine, api *gin.RouterGroup) {
	chainDeployGroup := api.Group("/chaindeploy")
	chainDeployGroup.POST("/list", srv.List)
	chainDeployGroup.POST("/add", srv.Add)
	chainDeployGroup.POST("/remove", srv.Delete)
	chainDeployGroup.POST("/getconfig", srv.GetConfig)
	chainDeployGroup.POST("/setconfig", srv.SetConfig)
}
