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

// {"chainName":"ft", "userID":1, "description":"ft的私链"}
// 添加链
func (srv *ChainService) ChainAdd(c *gin.Context) {
	var err error
	var chainInfo = &ChainInfo{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	chain := &model.Chain{
		Name: chainInfo.Name,
		UserID: chainInfo.UserID,
		Description: chainInfo.Desc,
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
	var chainInfo = &ChainInfo{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	chain := &model.Chain{
		UserID: chainInfo.UserID,
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
	var chainInfo = &ChainInfo{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	// 开启事务
	tx := srv.DB.Begin()

	if err := tx.Unscoped().Where("user_id = ? and chain_id = ?", chainInfo.UserID, chainInfo.ID).Delete(model.ChainDeploy{}).Error; nil != err {
		tx.Rollback()
		ginutil.Response(c, err, nil)
		return
	}


	// 删除chain的数据
	deleteDB := tx.Unscoped().Where("user_id = ?", chainInfo.UserID).Delete(model.Chain{})
	if err := deleteDB.Error; nil != err {
		tx.Rollback()
		ginutil.Response(c, DELETE_FAIL, nil)
		return
	}

	if 0 == deleteDB.RowsAffected {
		ginutil.Response(c, CHAINID_NOT_EXIST, nil)
	}

	// 结束事务
	tx.Commit()
	ginutil.Response(c, nil, nil)
}


// {"chainID":4, "name":"ft", "userID":1, "description":"ft的私链1"}
// 链更新
func (srv *ChainService) ChainUpdate(c *gin.Context) {
	var err error
	var chainInfo = &ChainInfo{}
	if err = c.ShouldBindJSON(chainInfo); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}
	chain := &model.Chain{}
	chain.ID = chainInfo.ID

	result := &model.Chain{
		Name: chainInfo.Name,
		UserID: chainInfo.UserID,
		Description: chainInfo.Desc,
	}

	updateDB := srv.DB.Model(&chain).Updates(result)
	if err = updateDB.Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}
	if 0 == updateDB.RowsAffected {
		ginutil.Response(c, UPDATE_FAIL, nil)
		return
	}

	ginutil.Response(c, nil, chain)
}
