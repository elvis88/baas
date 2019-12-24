package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/elvis88/baas/core/model"
	"github.com/elvis88/baas/db"
	"github.com/gin-gonic/gin"
)

// ChainService 区块链配置表
type ChainService struct {
}

// Register
func (srv *ChainService) Register(router *gin.Engine) {
	chain := router.Group("/chain")
	chain.POST("/add", AddChain)
	chain.POST("/gets", GetChains)
	chain.POST("/delete", DeleteChain)
	chain.POST("/update", UpdateChain)
}

func getParams(c *gin.Context) (map[string]interface{}, error) {
	// 获取请求主题
	body, err := ioutil.ReadAll(c.Request.Body)
	if nil != err {
		return nil, err
	}

	var params map[string]interface{}
	err = json.Unmarshal(body, &params)
	if nil != err {
		return nil, err
	}

	return params, nil
}

// {"name":"ft", "userID":1, "description":"ft的私链"}
func AddChain(c *gin.Context) {
	// 验证身份

	// 获取请求主体
	params, err := getParams(c)
	if nil != err {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  err.Error(),
		})
		return
	}

	// 获取主体参数
	var chainName string
	var userID    float64
	var desc      string
	var ok bool
	if chainName, ok = params["name"].(string); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	if userID, ok = params["userID"].(float64); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	if desc, ok = params["description"].(string); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	// install 链数据
	var chain = &model.Chain{
		Name: chainName,
		UserID: uint(userID),
		Description: desc,
	}

	if err = db.DB.Create(chain).Error; nil != err {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":true,
	})
}

// {"userID":1}
func GetChains(c *gin.Context) {
	// 验证身份

	// 获取请求主题
	params, err := getParams(c)
	if nil != err {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  err.Error(),
		})
		return
	}

	// 获取账户id
	var userID float64
	var ok bool
	if userID, ok = params["userID"].(float64); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	var chains []*model.Chain
	if err = db.DB.Where("user_id = ?", uint(userID)).Find(&chains).Error; nil != err {
		c.JSON(http.StatusOK, gin.H{
			"code":-32000,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": chains,
	})	
}

// {"userID":1, "chainID": 1}
func DeleteChain(c *gin.Context) {
	// 验证身份
	
	// 获取请求主题
	params, err := getParams(c)
	if nil != err {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  err.Error(),
		})
		return
	}

	// 获取账户id 链id
	var userID float64
	var chainID float64
	var ok bool
	if userID, ok = params["userID"].(float64); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	if chainID, ok = params["chainID"].(float64); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	// 开启事务
	tx := db.DB.Begin()

	if err := tx.Where("user_id = ? and chain_id = ?", uint(userID), uint(chainID)).Delete(model.ChainDeploy{}).Error; nil != err {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"code":-32000,
			"msg": err.Error(),
		})
		return
	}


	// 删除chain的数据
	if err := tx.Where("user_id = ?", uint(userID)).Delete(model.Chain{}).Error; nil != err {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"code":-32000,
			"msg": err.Error(),
		})
		return
	}

	// 结束事务
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"result":true,
	})
}


// {"chainID":4, "name":"ft", "userID":1, "description":"ft的私链1"}
func UpdateChain(c *gin.Context) {
	// 验证身份

	// 获取请求主题
	params, err := getParams(c)
	if nil != err {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  err.Error(),
		})
		return
	}

	var chainID   float64
	var chainName string
	var userID    float64
	var desc      string
	var ok bool

	if chainID, ok = params["chainID"].(float64); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}
	if chainName, ok = params["name"].(string); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	if userID, ok = params["userID"].(float64); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	if desc, ok = params["description"].(string); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -32000,
			"msg":  "body params format error",
		})
		return
	}

	var chain = &model.Chain{}
	chain.ID = uint(chainID)

	if err = db.DB.Model(&chain).Updates(&model.Chain{Name: chainName, UserID: uint(userID),Description: desc,}).Error; nil != err {
		c.JSON(http.StatusOK, gin.H{
			"code":-32000,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
