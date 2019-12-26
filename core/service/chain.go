package service

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainService 区块链配置表
type ChainService struct {
	DB        *gorm.DB
}

// Register
func (srv *ChainService) Register(router *gin.RouterGroup) {
	chain := router.Group("/chain")
	chain.POST("/add", srv.ChainAdd)
	chain.POST("/join",srv.ChainJoin)
	chain.POST("/list", srv.ChainList)
	chain.POST("/delete", srv.ChainDelete)
	chain.POST("/exit", srv.ChainExit)
	chain.POST("/update", srv.ChainUpdate)
}

func Verification(c *gin.Context, db *gorm.DB, userID uint) bool {
	userService := &UserService{DB: db}
	ok, user := userService.hasAdminRole(c)
	if ok {
		return true
	}

	if userID == user.ID {
		return true
	} else {
		ginutil.Response(c, NOPERMISSION, nil)
		c.Abort()
		return false
	}
}

// {"chainName":"ft", "userID":1, "description":"ft的私链"}
// 添加链
func (srv *ChainService) ChainAdd(c *gin.Context) {
	var err error

	// 获取主体信息
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 设置用户ID
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	chain.UserID = user.ID

	// 创建链
	if err = srv.DB.Create(chain).Error; nil != err {
		ginutil.Response(c, ADD_CHAIN_FAIL, err)
		return
	}

	// 建立用户与链的关联
	if err = srv.DB.Model(user).Association("OwnerChains").Append(chain).Error; nil != err {
		ginutil.Response(c, ADD_CHAIN_FAIL, err)
	}

	ginutil.Response(c, nil, chain)
}

// 添加已有链
func (srv *ChainService) ChainJoin(c *gin.Context) {
	var err error

	// 获取主体信息
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 判断链ID是否存在
	if chain.ID == 0 {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 查询链是否存在
	if err = srv.DB.First(chain).Error; nil != err{
		ginutil.Response(c, CHAINID_NOT_EXIST, err)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	// 建立联系
	if err = srv.DB.Model(user).Association("OwnerChains").Append(chain).Error; nil != err {
		ginutil.Response(c, ADD_CHAIN_FAIL, err)
	}

	ginutil.Response(c, nil, chain)
}

// 获取链列表
func (srv *ChainService) ChainList(c *gin.Context) {
	var err error

	// 获取账户信息
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)

	// 获取与账户关联的链
	var chains []*model.Chain
	if err = srv.DB.Model(&user).Association("OwnerChains").Find(&chains).Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}

	ginutil.Response(c, nil, chains)
}

// 删除链
func (srv *ChainService)ChainDelete(c *gin.Context) {
	var err error

	// 获取主体信息
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 判断链ID是否存在
	if chain.ID == 0 {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 查询链信息
	if err = srv.DB.First(chain).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, err)
		return
	}

	// 验证当前用户是否有修改权
	if false == Verification(c, srv.DB, chain.UserID) {
		return
	}

	// 删除chain的数据
	deleteDB := srv.DB.Unscoped().Where(chain).Delete(model.Chain{})
	if err := deleteDB.Error; nil != err {
		ginutil.Response(c, DELETE_FAIL, nil)
		return
	}

	ginutil.Response(c, nil, nil)
}

// 退出链
func (srv *ChainService) ChainExit(c *gin.Context) {
	var err error

	// 获取主体信息
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 判断链ID是否存在
	if chain.ID == 0 {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 查询链是否存在
	if err = srv.DB.First(chain).Error; nil != err{
		ginutil.Response(c, CHAINID_NOT_EXIST, err)
		return
	}

	// 获取账户
	userService := &UserService{DB: srv.DB}
	_, user := userService.hasAdminRole(c)
	if user.ID != chain.UserID {
		ginutil.Response(c, NOPERMISSION, nil)
		return
	}

	// 删除联系
	if err = srv.DB.Model(user).Association("OwnerChains").Delete(chain).Error; nil != err {
		ginutil.Response(c, ADD_CHAIN_FAIL, err)
	}

	ginutil.Response(c, nil, chain)
}

// 链更新
func (srv *ChainService) ChainUpdate(c *gin.Context) {
	// 读取主体信息
	var err error
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 判断链ID是否存在
	if chain.ID == 0 {
		ginutil.Response(c, REQUEST_PARAM_INVALID, err)
		return
	}

	// 获取链对应的账户ID
	chainVerify := &model.Chain{Model:model.Model{ID: chain.ID}}
	if err = srv.DB.First(chainVerify).Error; nil != err{
		ginutil.Response(c, CHAINID_NOT_EXIST, nil)
		return
	}

	// 验证对该链有没有操作权限
	if false == Verification(c, srv.DB, chainVerify.UserID) {
		return
	}

	// 更新链
	updateDB := srv.DB.Model(chain).Updates(&model.Chain{Name: chain.Name, Description:chain.Description})
	if err = updateDB.Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}

	// 获取最新链数据
	if err = srv.DB.First(chain).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, nil)
		return
	}

	// 返回更新链结果
	ginutil.Response(c, nil, chain)
}
