package service

import (
	"strconv"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainDeployNodeStatusService 区块链运行节点
type ChainDeployNodeStatusService struct {
	DB *gorm.DB
}

// List 列表
func (srv *ChainDeployNodeStatusService) List(ctx *gin.Context) {
	nodeid := ctx.Param("chaindeploynode")
	id, err := strconv.ParseUint(nodeid, 10, 64)
	if err != nil || id == 0 {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

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

	var node model.ChainDeployNode
	if err := srv.DB.First(&node, uint(id)).Error; err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}
	if node.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	var nodestatuss []*model.ChainDeployNodeStatus
	offset := req.Page * req.PageSize
	if err := srv.DB.Where(&model.ChainDeployNodeStatus{
		ChainDeployNodeID: uint(id),
	}).Offset(offset).Limit(req.PageSize).Find(&nodestatuss).Error; nil != err {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, nodestatuss)
}

// Register ...
func (srv *ChainDeployNodeStatusService) Register(router *gin.Engine, api *gin.RouterGroup) {
	chainDeployGroup := api.Group("/chaindeploynodestatus")
	chainDeployGroup.POST("/list", srv.List)
}
