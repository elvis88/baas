package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/elvis88/baas/core"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("baas") // name of config file
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the feconfig.yaml file
	if err != nil {             // Handle errors reading the config file
		fmt.Println("read config file error: ", err)
		os.Exit(-1)
	}
}

func main() {
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
		panic(fmt.Sprintf("not support db engine %s\n", engine))
	}
	db, err := gorm.Open(engine, connStr)
	if err != nil {
		panic(fmt.Sprintf("db connect err %s\n", err))
	}
	defer db.Close()
	db.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}

	router := gin.Default()

	if err := core.Register(router, db); err != nil {
		panic(err)
	}

	// 设置服务端口
	servicePort := viper.GetString("baas.config.port")
	if err := router.Run(fmt.Sprintf(":%s", servicePort)); err != nil {
		panic(err)
	}
}
