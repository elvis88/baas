package service

import (
	"fmt"
	"strings"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AgentService 监控代理服务
type AgentService struct {
	DB *gorm.DB
}

// List 列表
func (srv *AgentService) List(ctx *gin.Context) {
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

	var agents []*model.Agent
	offset := req.Page * req.PageSize
	if err := srv.DB.Where(&model.Agent{
		UserID: user.ID,
	}).Offset(offset).Limit(req.PageSize).Find(&agents).Error; nil != err {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, agents)
}

// Add 新增
func (srv *AgentService) Add(ctx *gin.Context) {
	agentParams := &requestAgentParams{}
	if err := ctx.ShouldBindJSON(agentParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := agentParams.validateAdd(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	agent := &model.Agent{
		Name:        agentParams.Name,
		Description: agentParams.Description,
		UserID:      user.ID,
	}
	if err := agent.Add(srv.DB); nil != err {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, agent)
}

// Delete 删除
func (srv *AgentService) Delete(ctx *gin.Context) {
	agentParams := &requestAgentParams{}
	if err := ctx.ShouldBindJSON(agentParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := agentParams.validateID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	agent := &model.Agent{}
	if err := srv.DB.First(agent, agentParams.ID).Error; nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if agent.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	if err := agent.Remove(srv.DB); err != nil {
		ginutil.Response(ctx, EXEC_FAILED, err.Error())
		return
	}
	ginutil.Response(ctx, nil, nil)
}

// Start 安装
func (srv *AgentService) Start(ctx *gin.Context) {
	agentParams := &requestAgentParams{}
	if err := ctx.ShouldBindJSON(agentParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := agentParams.validateID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	agent := &model.Agent{}
	if err := srv.DB.First(agent, agentParams.ID).Error; nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if agent.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	host := ctx.Request.Host
	path := strings.Replace(ctx.Request.URL.Path, "/start", "/agent.sh", 1)
	res := fmt.Sprintf(`export AgentID=%d; export BASS_Authorization=%s; curl http://%s%s --header Authorization:${BASS_Authorization} -sSf | sh -s start`, agent.ID, token, host, path)
	ginutil.Response(ctx, nil, res)
}

// Stop 安装
func (srv *AgentService) Stop(ctx *gin.Context) {
	agentParams := &requestAgentParams{}
	if err := ctx.ShouldBindJSON(agentParams); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 验证参数
	if ok, errMsg := agentParams.validateID(); !ok {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, errMsg)
		return
	}

	agent := &model.Agent{}
	if err := srv.DB.First(agent, agentParams.ID).Error; nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, err.Error())
		return
	}

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if agent.UserID != user.ID {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID, "not self")
		return
	}

	token := ctx.GetHeader(headerTokenKey)
	host := ctx.Request.Host
	path := strings.Replace(ctx.Request.URL.Path, "/stop", "/agent.sh", 1)
	res := fmt.Sprintf(`export AgentID=%d; export BASS_Authorization=%s; curl http://%s%s --header Authorization:${BASS_Authorization} -sSf | sh -s stop`, agent.ID, token, host, path)
	ginutil.Response(ctx, nil, res)
}

// Register ...
func (srv *AgentService) Register(router *gin.Engine, api *gin.RouterGroup) {
	chainDeployGroup := api.Group("/agent")
	chainDeployGroup.POST("/list", srv.List)
	chainDeployGroup.POST("/add", srv.Add)
	chainDeployGroup.POST("/remove", srv.Delete)
	chainDeployGroup.POST("/start", srv.Start)
	chainDeployGroup.POST("/stop", srv.Stop)
}
