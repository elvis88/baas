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

// ChainDeployAdd 新增
func (srv *ChainDeployService) ChainDeployAdd(ctx *gin.Context) {
	var err error
	chainDeploy := &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	if err = srv.DB.Create(chainDeploy).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, chainDeploy)
}

func (srv *ChainDeployService) ChainDeployList(ctx *gin.Context) {
	var err error
	chainDeploy := &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
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
	var err error
	chainDeploy := &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
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
	var err error
	var chainDeploy = &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}
	chainDeployResult := &model.ChainDeploy{}
	chainDeployResult.ID = chainDeploy.ID

	updateDB := srv.DB.Model(&chainDeployResult).Updates(chainDeploy)
	if err = updateDB.Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}
	if 0 == updateDB.RowsAffected {
		ginutil.Response(ctx, UPDATE_FAIL, nil)
		return
	}

	ginutil.Response(ctx, nil, chainDeployResult)
}

// Register ...
func (srv *ChainDeployService) Register(api *gin.RouterGroup) {
	chainDeployGroup := api.Group("/chaindeploy")
	chainDeployGroup.POST("/add", srv.ChainDeployAdd)
	chainDeployGroup.POST("/list", srv.ChainDeployList)
	chainDeployGroup.POST("/delete", srv.ChainDeployDelete)
	chainDeployGroup.POST("/update", srv.ChainDeployUpdate)
}
