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

func (srv *ChainDeployService) getAndCheckParams(ctx *gin.Context) (cDeploy *model.ChainDeploy, err error) {
	chainDeploy := &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		return nil, err
	}

	return chainDeploy, nil
}

// ChainDeployAdd 新增
func (srv *ChainDeployService) ChainDeployAdd(ctx *gin.Context) {
	chainDeploy, err := srv.getAndCheckParams(ctx)
	if nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	if chainDeploy.UserID  == 0  || chainDeploy.ChainID == 0 {
		ginutil.Response(ctx, PARAMS_IS_NOT_ENOUGH, nil)
		return
	}

	if false == Verification(ctx, srv.DB, chainDeploy.UserID) {
		return
	}

	if err = srv.DB.Create(chainDeploy).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, chainDeploy)
}

func (srv *ChainDeployService) ChainDeployList(ctx *gin.Context) {
	chainDeploy, err := srv.getAndCheckParams(ctx)
	if nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	if chainDeploy.UserID  == 0  || chainDeploy.ChainID == 0 {
		ginutil.Response(ctx, PARAMS_IS_NOT_ENOUGH, nil)
		return
	}

	if false == Verification(ctx, srv.DB, chainDeploy.UserID) {
		return
	}

	var chainDeploys []*model.ChainDeploy
	if err = srv.DB.Where(chainDeploy).Find(&chainDeploys).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, chainDeploys)
}

// ChainDeployDelete 删除
func (srv *ChainDeployService) ChainDeployDelete(ctx *gin.Context) {
	chainDeploy, err := srv.getAndCheckParams(ctx)
	if nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	if chainDeploy.ID == 0 {
		ginutil.Response(ctx, PARAMS_IS_NOT_ENOUGH, nil)
		return
	}

	if err = srv.DB.First(chainDeploy).Error; nil != err {
		ginutil.Response(ctx, CHAINID_DEPLOY_NOT_EXIST, nil)
		return
	}

	if false == Verification(ctx, srv.DB, chainDeploy.UserID) {
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
	// 获取主体信息
	chainDeploy, err := srv.getAndCheckParams(ctx)
	if nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	// 如果实例id不存在则返回
	if chainDeploy.ID == 0 {
		ginutil.Response(ctx, PARAMS_IS_NOT_ENOUGH, nil)
		return
	}

	// 获取该ID对应的数据
	chainDeployVerify := &model.ChainDeploy{Model: model.Model{ID:chainDeploy.ID}}
	if err = srv.DB.First(chainDeployVerify).Error; nil != err {
		ginutil.Response(ctx, CHAINID_DEPLOY_NOT_EXIST, nil)
		return
	}

	// 验证用户是否有修改权限
	if false == Verification(ctx, srv.DB, chainDeployVerify.UserID) {
		return
	}

	// 存储更新的结果
	updateDB := srv.DB.Model(chainDeploy).Updates(&model.Chain{Name:chainDeploy.Name, Description:chainDeploy.Description})
	if err = updateDB.Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	// 获取最新链实例数据
	if err = srv.DB.First(chainDeploy).Error; nil != err {
		ginutil.Response(ctx, CHAINID_DEPLOY_NOT_EXIST, nil)
		return
	}

	// 返回更新后的数据
	ginutil.Response(ctx, nil, chainDeploy)
}

// Register ...
func (srv *ChainDeployService) Register(api *gin.RouterGroup) {
	chainDeployGroup := api.Group("/chaindeploy")
	chainDeployGroup.POST("/add", srv.ChainDeployAdd)
	chainDeployGroup.POST("/list", srv.ChainDeployList)
	chainDeployGroup.POST("/delete", srv.ChainDeployDelete)
	chainDeployGroup.POST("/update", srv.ChainDeployUpdate)
}
