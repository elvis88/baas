package service

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChainService 区块链配置表
type ChainService struct {
	DB *gorm.DB
}

type ChainInfo struct {
	ID 		uint 	`json:"id"`
	Name 	string 	`json:"chainName"`
	UserID  uint	`json:"userID"`
	Desc    string  `json:"description"`
}

// Register
func (srv *ChainService) Register(router *gin.RouterGroup) {
	chain := router.Group("/chain")
	chain.POST("/add", srv.ChainAdd)
	chain.POST("/list", srv.ChainList)
	chain.POST("/delete", srv.ChainDelete)
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
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	if false == Verification(c, srv.DB, chain.UserID) {
		return
	}


	if err = srv.DB.Create(chain).Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}

	ginutil.Response(c, nil, chain)
}

// {"userID":1}
// 获取链列表
func (srv *ChainService) ChainList(c *gin.Context) {
	var err error
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	if false == Verification(c, srv.DB, chain.UserID) {
		return
	}


	var chains []*model.Chain
	if err = srv.DB.Where(chain).Find(&chains).Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}

	ginutil.Response(c, nil, chains)
}

// {"userID":1, "chainID": 1}
// 链删除
func (srv *ChainService)ChainDelete(c *gin.Context) {
	var err error
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	if chain.ID == 0 {
		ginutil.Response(c, PARAMS_IS_NOT_ENOUGH, nil)
		return
	}

	if err = srv.DB.First(chain).Error; nil != err {
		ginutil.Response(c, CHAINID_NOT_EXIST, nil)
		return
	}

	if false == Verification(c, srv.DB, chain.UserID) {
		return
	}

	// 删除chain的数据
	deleteDB := srv.DB.Unscoped().Where("id = ?", chain.ID).Delete(model.Chain{})
	if err := deleteDB.Error; nil != err {
		ginutil.Response(c, DELETE_FAIL, nil)
		return
	}

	ginutil.Response(c, nil, nil)
}


// {"chainID":4, "name":"ft", "userID":1, "description":"ft的私链1"}
// 链更新
func (srv *ChainService) ChainUpdate(c *gin.Context) {
	// 读取主体信息
	var err error
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 如果主体中不包含链ID，则直接返回
	if chain.ID == 0 {
		ginutil.Response(c, PARAMS_IS_NOT_ENOUGH, nil)
		return
	}

	// 获取链对应的账户ID
	chainVerify := &model.Chain{Model:model.Model{ID: chain.ID}}
	if err = srv.DB.First(chainVerify).Error; nil != err {
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
