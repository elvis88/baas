package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/elvis88/baas/core/ws"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainDeployNodeService 区块链运行节点
type ChainDeployNodeService struct {
	DB *gorm.DB
}

// List 列表
func (srv *ChainDeployNodeService) List(ctx *gin.Context) {
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

	node := &model.ChainDeployNode{
		UserID: user.ID,
	}
	agent, ok := ctx.GetQuery("agent")
	if ok {
		id, err := strconv.ParseUint(agent, 10, 64)
		if err != nil {
			ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
			return
		}
		node.AgentID = uint(id)
	}
	chaindeploy, ok := ctx.GetQuery("chaindeploy")
	if ok {
		id, err := strconv.ParseUint(chaindeploy, 10, 64)
		if err != nil {
			ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
			return
		}
		node.ChainDeployID = uint(id)
	}

	var nodes []*model.ChainDeployNode
	offset := req.Page * req.PageSize
	if err := srv.DB.Where(node).Offset(offset).Limit(req.PageSize).Find(&nodes).Error; nil != err {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, nodes)
}

// Add 新增
func (srv *ChainDeployNodeService) Add(ctx *gin.Context) {
	nodeParams := &requestChainDeloyNodeParams{}
	if err := ctx.ShouldBindJSON(nodeParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := nodeParams.validateAdd(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	agent := &model.Agent{}
	if err := srv.DB.First(agent, nodeParams.AgentID).Error; err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	deploy := &model.ChainDeploy{}
	if err := srv.DB.First(deploy, nodeParams.ChainDeployID).Error; err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)
	_ = user
	if agent.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "agent not self")
		return
	}

	if deploy.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "deploy not self")
		return
	}

	// 添加链实例
	node := &model.ChainDeployNode{
		Name:          nodeParams.Name,
		Description:   nodeParams.Description,
		AgentID:       nodeParams.AgentID,
		ChainDeployID: nodeParams.ChainDeployID,
		UserID:        user.ID,
	}

	if err := node.Add(srv.DB); nil != err {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}

	var chain model.Chain
	srv.DB.First(&chain, deploy.ChainID)
	ws.HandleAddNode(agent.ID, deploy.ID, deploy.Name, chain.Name)
	ginutil.Response(ctx, nil, node)
}

// Delete 删除
func (srv *ChainDeployNodeService) Delete(ctx *gin.Context) {
	nodeParams := &requestChainDeloyNodeParams{}
	if err := ctx.ShouldBindJSON(nodeParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := nodeParams.validateID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 指定id获取实例数据
	node := &model.ChainDeployNode{}
	if err := srv.DB.First(node, nodeParams.ID).Error; nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)
	_ = user
	if node.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	if err := node.Remove(srv.DB); err != nil {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}

	ws.HandleRemoveNode(node.AgentID, node.ID)
	ginutil.Response(ctx, nil, nil)
}

// Start 启动
func (srv *ChainDeployNodeService) Start(ctx *gin.Context) {
	nodeParams := &requestChainDeloyNodeParams{}
	if err := ctx.ShouldBindJSON(nodeParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := nodeParams.validateID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 指定id获取实例数据
	node := &model.ChainDeployNode{}
	if err := srv.DB.First(node, nodeParams.ID).Error; nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)
	_ = user
	if node.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	var chaindeploy model.ChainDeploy
	srv.DB.First(&chaindeploy, node.ChainDeployID)
	path, err := chaindeploy.GetScriptPath(srv.DB)
	if err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	host := ctx.Request.Host
	path = strings.Replace(ctx.Request.URL.Path, "/chaindeploynode/start", path, 1)
	res := fmt.Sprintf(`export BASS_Authorization=%s; curl http://%s%s --header Authorization:${BASS_Authorization} -sSf | sh -s start`, token, host, path)

	var agent model.Agent
	srv.DB.First(&agent, node.AgentID)

	res, err = ws.HandleCommand(agent.ID, res)
	if err != nil {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}

	ginutil.Response(ctx, nil, res)
}

// Stop 停止
func (srv *ChainDeployNodeService) Stop(ctx *gin.Context) {
	nodeParams := &requestChainDeloyNodeParams{}
	if err := ctx.ShouldBindJSON(nodeParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := nodeParams.validateID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 指定id获取实例数据
	node := &model.ChainDeployNode{}
	if err := srv.DB.First(node, nodeParams.ID).Error; nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)
	_ = user
	if node.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	var chaindeploy model.ChainDeploy
	srv.DB.First(&chaindeploy, node.ChainDeployID)
	path, err := chaindeploy.GetScriptPath(srv.DB)
	if err != nil {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	host := ctx.Request.Host
	path = strings.Replace(ctx.Request.URL.Path, "/chaindeploynode/stop", path, 1)
	res := fmt.Sprintf(`export BASS_Authorization=%s; curl http://%s%s --header Authorization:${BASS_Authorization} -sSf | sh -s stop`, token, host, path)

	var agent model.Agent
	srv.DB.First(&agent, node.AgentID)

	res, err = ws.HandleCommand(agent.ID, res)
	if err != nil {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}

	ginutil.Response(ctx, nil, res)
}

// Register ...
func (srv *ChainDeployNodeService) Register(router *gin.Engine, api *gin.RouterGroup) {
	chainDeployGroup := api.Group("/chaindeploynode")
	chainDeployGroup.POST("/list", srv.List)
	chainDeployGroup.POST("/add", srv.Add)
	chainDeployGroup.POST("/remove", srv.Delete)
	chainDeployGroup.POST("/start", srv.Start)
	chainDeployGroup.POST("/stop", srv.Stop)
}
