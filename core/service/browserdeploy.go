package service

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// BrowserDeployService 浏览器配置表
type BrowserDeployService struct {
	DB *gorm.DB
}

// BrowserDeployAdd 新增
func (srv *BrowserDeployService) BrowserDeployAdd(ctx *gin.Context) {
	var err error
	browserDeploy := &model.BrowserDeploy{}
	if err = ctx.ShouldBindJSON(browserDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	if err = srv.DB.Create(browserDeploy).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, browserDeploy)
}

func (srv *BrowserDeployService) BrowserDeployList(ctx *gin.Context) {
	var err error
	browserDeploy := &model.BrowserDeploy{}
	if err = ctx.ShouldBindJSON(browserDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	var browserDeploys []*model.BrowserDeploy
	if err = srv.DB.Where(browserDeploy).Find(&browserDeploys).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, browserDeploys)
}

// BrowserDeployDelete 删除
func (srv *BrowserDeployService) BrowserDeployDelete(ctx *gin.Context) {
	var err error
	browserDeploy := &model.BrowserDeploy{}
	if err = ctx.ShouldBindJSON(browserDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	deleteDB := srv.DB.Unscoped().Delete(browserDeploy)
	if err = deleteDB.Error; nil != err {
		ginutil.Response(ctx, DELETE_FAIL, nil)
		return
	}

	ginutil.Response(ctx, nil, nil)
}

// BrowserDeployUpdate 修改
func (srv *BrowserDeployService) BrowserDeployUpdate(ctx *gin.Context) {
	var err error
	var browserDeploy = &model.BrowserDeploy{}
	if err = ctx.ShouldBindJSON(browserDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}
	chainDeployResult := &model.ChainDeploy{}
	chainDeployResult.ID = browserDeploy.ID

	updateDB := srv.DB.Model(&chainDeployResult).Updates(browserDeploy)
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
func (srv *BrowserDeployService) Register(router *gin.Engine, api *gin.RouterGroup) {
	browserDeployGroup := api.Group("/browserdeploy")
	browserDeployGroup.POST("/add", srv.BrowserDeployAdd)
	browserDeployGroup.POST("/list", srv.BrowserDeployList)
	browserDeployGroup.POST("/delete", srv.BrowserDeployDelete)
	browserDeployGroup.POST("/update", srv.BrowserDeployUpdate)
}
