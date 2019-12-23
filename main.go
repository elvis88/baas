package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/elvis88/baas/common/ginutil"
	"github.com/elvis88/baas/common/log"
	"github.com/elvis88/baas/core"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func main() {
	// db
	username := viper.GetString("baas.mysql.user")
	password := viper.GetString("baas.mysql.password")
	ip := viper.GetString("baas.mysql.ip")
	port := viper.GetString("baas.mysql.port")
	database := viper.GetString("baas.mysql.database")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=%s&parseTime=true",
		username, password, ip, port, database, url.QueryEscape("Asia/Shanghai"))
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		_ = fmt.Errorf("connect", err)
		os.Exit(-1)
	}

	logger := log.GetLogger("baas", log.DEBUG)

	router := gin.New()
	router.Use(ginutil.UseLogger(router, logger.Debugf))
	router.Use(gin.Recovery())
	core.Server(router, db)
}

func init() {
	viper.SetConfigName("bass") // name of config file
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the feconfig.yaml file
	if err != nil {             // Handle errors reading the config file
		fmt.Println("read config file error: %s \n", err)
		os.Exit(-1)
	}

	//全局配置
	fmt.Println("load config %v", viper.AllSettings())
}
