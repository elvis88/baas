package service

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/common/jwt"
	"github.com/elvis88/baas/core/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"time"
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

func Verification(c *gin.Context, userID uint) bool {
	token := c.GetHeader(headerTokenKey)
	session := ginutil.GetSession(c, token)
	if nil == session {
		ginutil.Response(c, TOKEN_NOT_EXIST, nil)
		c.Abort()
		return false
	}

	info, ok := jwt.ParseToken(token, TokenKey)
	if !ok {
		ginutil.Response(c, TOKEN_INVALID, nil)
		c.Abort()
		return false
	}

	var infoMap map[string]interface{}
	if infoMap = info.(map[string]interface{}); float64(time.Now().Unix()) >= infoMap["exp"].(float64) {
		ginutil.Response(c, TOKEN_EXPIRE, nil)
		c.Abort()
		return false
	}

	if tokenUserID,ok := infoMap["userId"].(float64); ok && userID == uint(tokenUserID) {
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

	if false == Verification(c, chain.UserID) {
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

	if false == Verification(c, chain.UserID) {
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

	if false == Verification(c, chain.UserID) {
		return
	}


	// 开启事务
	tx := srv.DB.Begin()

	if err := tx.Unscoped().Where("user_id = ? and chain_id = ?", chain.UserID, chain.ID).Delete(model.ChainDeploy{}).Error; nil != err {
		tx.Rollback()
		ginutil.Response(c, err, nil)
		return
	}


	// 删除chain的数据
	deleteDB := tx.Unscoped().Where("id = ?", chain.UserID).Delete(model.Chain{})
	if err := deleteDB.Error; nil != err {
		tx.Rollback()
		ginutil.Response(c, DELETE_FAIL, nil)
		return
	}

	if 0 == deleteDB.RowsAffected {
		ginutil.Response(c, CHAINID_NOT_EXIST, nil)
		return
	}

	// 结束事务
	tx.Commit()
	ginutil.Response(c, nil, nil)
}


// {"chainID":4, "name":"ft", "userID":1, "description":"ft的私链1"}
// 链更新
func (srv *ChainService) ChainUpdate(c *gin.Context) {
	var err error
	var chain = &model.Chain{}
	if err = c.ShouldBindJSON(chain); nil != err {
		ginutil.Response(c, REQUEST_PARAM_INVALID, nil)
		return
	}

	if false == Verification(c, chain.UserID) {
		return
	}

	result := &model.Chain{}
	result.ID = chain.ID

	updateDB := srv.DB.Model(result).Updates(chain)
	if err = updateDB.Error; nil != err {
		ginutil.Response(c, err, nil)
		return
	}
	if 0 == updateDB.RowsAffected {
		ginutil.Response(c, UPDATE_FAIL, nil)
		return
	}

	ginutil.Response(c, nil, result)
}
