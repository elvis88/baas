package service

import (
	"errors"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/generate"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainDeployService 区块链配置表
type ChainDeployService struct {
	DB *gorm.DB
}


// chainDeployValidate
type requestChainDeployParams struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	UserID      string `json:"userID"`
	ChainID     uint   `json:"chainID"`
	Description string `json:"desc"`
}


func (srv *ChainDeployService) userHaveChain(userID, chainID uint) (b bool, err error) {
	var chain model.Chain

	// 查看用户是否和该链有关联
	if err = srv.DB.Model(&model.User{Model: model.Model{ID: userID}}).
		Where(&model.Chain{Model: model.Model{ID: chainID}}).
		Association("OwnerChains").Find(&chain).Error; nil != err || chain.ID == 0 {
		return false, err
	}

	return true, nil
}

// ChainDeployAdd 新增
func (srv *ChainDeployService) ChainDeployAdd(ctx *gin.Context) {
	// 获取请求参数
	var err error
	chainDeployParams := &requestChainDeployParams{}
	if err = ctx.ShouldBindJSON(chainDeployParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err)
		return
	}

	// 验证参数
	if ok, errMsg := chainDeployParams.validateChainDeployAdd(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证用户是否拥有链
	ok, err := srv.userHaveChain(user.ID, chainDeployParams.ChainID)
	if nil != err || !ok {
		ginutil.Response(ctx, CHAINID_DEPLOY_ADD_FAIL, err)
		return
	}

	// 添加链实例
	chainDeploy := &model.ChainDeploy{
		Name: chainDeployParams.Name,
		UserID: user.ID,
		ChainID: chainDeployParams.ChainID,
		Description: chainDeployParams.Description,
	}

	chain := &model.Chain{}
	if err := srv.DB.First(chain, chainDeploy.ChainID).Error; err != nil {
		ginutil.Response(ctx, CHAINID_NOT_EXIST, err)
		return
	}
	orginChain := &model.Chain{}
	if err := srv.DB.First(orginChain, chain.OriginID).Error; err != nil {
		ginutil.Response(ctx, CHAINID_NOT_EXIST, err)
		return
	}

	spec := generate.NewAppDeploySpec(user.Name, chainDeploy.Name, orginChain.Name)
	if spec == nil {
		ginutil.Response(ctx, ADD_CHAIN_FAIL, errors.New("not support"))
		return
	}
	err = spec.Build()
	if err != nil {
		ginutil.Response(ctx, ADD_CHAIN_FAIL, err)
		return
	}
	defer func() {
		if err != nil {
			spec.Remove()
		}
	}()

	if err = srv.DB.Create(chainDeploy).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, chainDeploy)
}

func (srv *ChainDeployService) ChainDeployList(ctx *gin.Context) {
	req := &PagerRequest{
		Page:     0,
		PageSize: 10,
	}

	// 获取请求参数
	var err error
	if err = ctx.ShouldBindJSON(req); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err)
		return
	}

	// 获取账户
	offset := req.Page * req.PageSize
	userService := &UserService{DB: srv.DB}
	ok, user := userService.hasAdminRole(ctx)

	// 获取实例列表
	var chainDeploys []*model.ChainDeploy

	// admin查看所有节点
	if ok {
		if err = srv.DB.Offset(offset).Limit(req.PageSize).Where(&model.ChainDeploy{}).Find(&chainDeploys).Error; nil != err {
			ginutil.Response(ctx, err, nil)
			return
		}
	} else {
		if err = srv.DB.Offset(offset).Limit(req.PageSize).Where(&model.ChainDeploy{
			UserID: user.ID,
		}).Find(&chainDeploys).Error; nil != err {
			ginutil.Response(ctx, err, nil)
			return
		}
	}

	// 返回链实例列表
	ginutil.Response(ctx, nil, chainDeploys)
}

// ChainDeployDelete 删除
func (srv *ChainDeployService) ChainDeployDelete(ctx *gin.Context) {
	// 获取请求参数
	var err error
	chainDeployParams := &requestChainDeployParams{}
	if err = ctx.ShouldBindJSON(chainDeployParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err)
		return
	}

	// 验证参数
	if ok, errMsg := chainDeployParams.validateChainDeployID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 指定id获取实例数据
	chainDeploy := &model.ChainDeploy{Model: model.Model{ID: chainDeployParams.ID}}
	if err = srv.DB.First(chainDeploy).Error; nil != err {
		ginutil.Response(ctx, CHAINID_DEPLOY_NOT_EXIST, nil)
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if chainDeploy.UserID != user.ID {
		ginutil.Response(ctx, NOPERMISSION, nil)
		return
	}

	deleteDB := srv.DB.Unscoped().Delete(chainDeploy)
	if err = deleteDB.Error; nil != err {
		ginutil.Response(ctx, DELETE_FAIL, nil)
		return
	}

	ginutil.Response(ctx, nil, nil)
}

// ChainDeployUpdate 修改
func (srv *ChainDeployService) ChainDeployUpdate(ctx *gin.Context) {
	// 获取请求参数
	var err error
	chainDeployParams := &requestChainDeployParams{}
	if err = ctx.ShouldBindJSON(chainDeployParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err)
		return
	}

	// 验证参数
	if ok, errMsg := chainDeployParams.validateChainDeployUpdate(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 获取该ID对应的数据
	chainDeployVerify := &model.ChainDeploy{Model: model.Model{ID: chainDeployParams.ID}}
	if err = srv.DB.First(chainDeployVerify).Error; nil != err {
		ginutil.Response(ctx, CHAINID_DEPLOY_NOT_EXIST, nil)
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if chainDeployVerify.UserID != user.ID {
		ginutil.Response(ctx, NOPERMISSION, nil)
		return
	}

	// 存储更新的结果
	updateDB := srv.DB.Model(chainDeployVerify).Updates(&model.Chain{
		Name: chainDeployParams.Name,
		Description: chainDeployParams.Description,
	})
	if err = updateDB.Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	// 获取最新链实例数据
	if err = srv.DB.First(chainDeployVerify).Error; nil != err {
		ginutil.Response(ctx, CHAINID_DEPLOY_NOT_EXIST, nil)
		return
	}

	// 返回更新后的数据
	ginutil.Response(ctx, nil, chainDeployVerify)
}

func (srv *ChainDeployService) ChainGetConfig(c *gin.Context) {

	chainDeploy := &model.ChainDeploy{}
	user := &model.User{}
	orgChain := &model.Chain{}
	spec := generate.NewAppDeploySpec(user.Name, chainDeploy.Name, orgChain.Name)
	if spec == nil {
		ginutil.Response(c, nil, nil)
		return
	}
	config, err := spec.GetConfig()
	if err != nil {
		ginutil.Response(c, nil, err)
		return
	}
	ginutil.Response(c, nil, config)
}

func (srv *ChainDeployService) ChainSetConfig(c *gin.Context) {

}

func (srv *ChainDeployService) ChainGetDeploy(c *gin.Context) {

}

// Register ...
func (srv *ChainDeployService) Register(router *gin.Engine, api *gin.RouterGroup) {
	chainDeployGroup := api.Group("/chaindeploy")
	chainDeployGroup.POST("/add", srv.ChainDeployAdd)
	chainDeployGroup.POST("/list", srv.ChainDeployList)
	chainDeployGroup.POST("/delete", srv.ChainDeployDelete)
	chainDeployGroup.POST("/update", srv.ChainDeployUpdate)
	chainDeployGroup.POST("/getcoinfig", srv.ChainSetConfig)
	chainDeployGroup.POST("/setcoinfig", srv.ChainGetConfig)
	chainDeployGroup.POST("/deploy", srv.ChainGetDeploy)
}
