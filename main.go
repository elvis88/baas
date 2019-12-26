package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/elvis88/baas/common/log"
	"github.com/elvis88/baas/core"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

func main() {
	// 初始化log
	logger := log.GetLogger("baas", log.DEBUG)

	// service signal handler
	signalHandler()

	// db 初始化
	connStr := ""
	engine := viper.GetString("baas.dbengine")
	if strings.Compare(engine, "mysql") == 0 {
		username := viper.GetString("baas.mysql.user")
		password := viper.GetString("baas.mysql.password")
		ip := viper.GetString("baas.mysql.ip")
		port := viper.GetString("baas.mysql.port")
		database := viper.GetString("baas.mysql.database")
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=%s&parseTime=true",
			username, password, ip, port, database, url.QueryEscape("Asia/Shanghai"))
	} else if strings.Compare(engine, "sqlite3") == 0 {
		connStr = viper.GetString("baas.sqlite3.root")
	}

	if len(connStr) == 0 {
		logger.Errorf("not support db engine %s\n", engine)
		return
	}

	db, err := gorm.Open(engine, connStr)
	if err != nil {
		logger.Errorf("db connect err %s\n", err)
		return
	}

	defer db.Close()
	db.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}

	// 创建服务
	router := gin.Default()

	// 注册服务
	if err := core.Server(router, db); err != nil {
		logger.Error(err)
		return
	}

	// 设置服务端口
	servicePort := viper.GetString("baas.config.port")
	if err := router.Run(fmt.Sprintf(":%s", servicePort)); err != nil {
		logger.Error(err)
		return
	}
}

func init() {
	viper.SetConfigName("baas") // name of config file
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the feconfig.yaml file
	if err != nil {             // Handle errors reading the config file
		fmt.Println("read config file error: ", err)
		os.Exit(-1)
	}

	//全局配置
	// fmt.Println("load config: ", viper.AllSettings())
}

func signalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("system exit")
		os.Exit(0)
	}()
}
