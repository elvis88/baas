package core

import (
	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/common/password"
	"github.com/elvis88/baas/common/sms"
	"github.com/elvis88/baas/core/model"
	srv "github.com/elvis88/baas/core/service"
	"github.com/elvis88/baas/core/ws"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Register 提供服务
func Register(router *gin.Engine, db *gorm.DB) error {
	// 创建数据库表
	if err := db.AutoMigrate(
		model.User{},
		model.Chain{}, model.ChainDeploy{}, model.ChainDeployNode{},
		model.ChainStatus{}, model.Agent{}, model.ChainDeployNodeStatus{}).Error; err != nil {
		return err
	}

	// 初始化用户Admin
	pwd, _ := password.CryTo(sms.GetRandomString(6), 12, "default")
	adminUsr := &model.User{
		Name:     "admin",
		Password: pwd,
	}
	usrs := []*model.User{
		adminUsr,
	}
	for _, usr := range usrs {
		if err := db.FirstOrCreate(usr, &model.User{
			Name: usr.Name,
		}).Error; err != nil {
			return err
		}
	}

	// 初始chain
	ftChain := &model.Chain{
		Name:        "ft",
		URL:         "https://github.com/fractalplatform/fractal",
		Description: "fractalplatform ",
		UserID:      adminUsr.ID,
	}
	fttestChain := &model.Chain{
		Name:        "fttest",
		URL:         "https://github.com/fractalplatform/fractal",
		Description: "fractalplatform",
		UserID:      adminUsr.ID,
	}
	chains := []*model.Chain{
		ftChain,
		fttestChain,
	}
	for _, chain := range chains {
		if err := db.FirstOrCreate(chain, &model.Chain{
			Name: chain.Name,
		}).Error; err != nil {
			return err
		}
	}

	ws.Run(router, db, gin.Mode() != gin.ReleaseMode)

	services := []service{
		&srv.UserService{
			DB: db,
		},
		&srv.AgentService{
			DB: db,
		},
		&srv.ChainService{
			DB: db,
		},
		&srv.ChainDeployService{
			DB: db,
		},
		&srv.ChainDeployNodeService{
			DB: db,
		},
		&srv.ChainDeployNodeStatusService{
			DB: db,
		},
	}

	ginutil.UseSession(router)
	apiv1 := router.Group("api/v1")
	for _, service := range services {
		service.Register(router, apiv1)
	}
	return nil
}

type service interface {
	Register(router *gin.Engine, api *gin.RouterGroup)
}
