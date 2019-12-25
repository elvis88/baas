package service

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// BrowserService 浏览器配置服务
type BrowserService struct {
	DB *gorm.DB
}

// BrowserAdd 新增
func (srv *BrowserService) BrowserAdd(ctx *gin.Context) {
	var err error
	var browser = &model.Browser{}
	if err = ctx.ShouldBindJSON(browser); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	if err = srv.DB.Create(browser).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, browser)
}

func (srv *BrowserService) BrowserList(ctx *gin.Context) {
	var err error
	var browser = &model.Browser{}
	if err = ctx.ShouldBindJSON(browser); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	var browsers []*model.Browser
	if err = srv.DB.Where(browser).Find(&browsers).Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}

	ginutil.Response(ctx, nil, browsers)
}

// BrowserDelete 删除
func (srv *BrowserService) BrowserDelete(ctx *gin.Context) {
	var err error
	var browser = &model.Browser{}
	if err = ctx.ShouldBindJSON(browser); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 开启事务
	tx := srv.DB.Begin()

	if err := tx.Unscoped().Where("user_id = ? and browser_id = ?", browser.UserID, browser.ID).Delete(model.BrowserDeploy{}).Error; nil != err {
		tx.Rollback()
		ginutil.Response(ctx, err, nil)
		return
	}


	// 删除浏览器的数据
	deleteDB := tx.Unscoped().Where("user_id = ?", browser.UserID).Delete(model.Browser{})
	if err := deleteDB.Error; nil != err {
		tx.Rollback()
		ginutil.Response(ctx, DELETE_FAIL, nil)
		return
	}

	if 0 == deleteDB.RowsAffected {
		ginutil.Response(ctx, CHAINID_NOT_EXIST, nil)
	}

	// 结束事务
	tx.Commit()
	ginutil.Response(ctx, nil, nil)
}

// BrowserUpdate 修改
func (srv *BrowserService) BrowserUpdate(ctx *gin.Context) {
	var err error
	var browser = &model.Browser{}
	if err = ctx.ShouldBindJSON(browser); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	result := &model.Chain{}
	result.ID = browser.ID

	updateDB := srv.DB.Model(result).Updates(browser)
	if err = updateDB.Error; nil != err {
		ginutil.Response(ctx, err, nil)
		return
	}
	if 0 == updateDB.RowsAffected {
		ginutil.Response(ctx, UPDATE_FAIL, nil)
		return
	}

	ginutil.Response(ctx, nil, result)
}

// Register ...
func (srv *BrowserService) Register(api *gin.RouterGroup) {
	browserGroup := api.Group("/browser")
	browserGroup.POST("/add", srv.BrowserAdd)
	browserGroup.POST("/list", srv.BrowserList)
	browserGroup.POST("/delete", srv.BrowserDelete)
	browserGroup.POST("/update", srv.BrowserUpdate)
}
