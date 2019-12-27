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
	var chains []*model.Chain

	// 未验证(关联查询如何添加条件)
	if err = srv.DB.Model(&model.User{Model: model.Model{ID:userID},}).
		Where(&model.Chain{Model: model.Model{ID:chainID}}).
		Association("OwnerChains").Find(&chains).Error; nil != err {
		return false, err
	}

	return false, nil
}

// ChainDeployAdd 新增
func (srv *ChainDeployService) ChainDeployAdd(ctx *gin.Context) {
	// 获取请求参数
	var err error
	chainDeploy := &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID ,err)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证用户是否拥有链
	ok, err := srv.userHaveChain(user.ID, chainDeploy.ChainID)
	if nil != err || !ok {
		ginutil.Response(ctx, ADD_CHAIN_DEPLOY_FAIL, err)
		return
	}

	chainDeploy.UserID = user.ID

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
		ginutil.Response(ctx, REQUEST_PARAM_INVALID ,err)
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
	chainDeploy := &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID ,err)
		return
	}

	// 指定id获取实例数据
	if chainDeploy.ID == 0 {
		ginutil.Response(ctx, PARAMS_IS_NOT_ENOUGH, nil)
		return
	}
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
	chainDeploy := &model.ChainDeploy{}
	if err = ctx.ShouldBindJSON(chainDeploy); nil != err {
		ginutil.Response(ctx, REQUEST_PARAM_INVALID ,err)
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

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(ctx)

	// 验证当前用户是否有修改权(admin 不可以删除)
	if chainDeploy.UserID != user.ID {
		ginutil.Response(ctx, NOPERMISSION, nil)
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
