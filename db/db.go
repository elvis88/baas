package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	//"github.com/elvis88/baas/core/model"
	"os"
)

var DB *gorm.DB

func InitDb(connStr string) {
	var err error
	DB, err = gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Println("connect fail", err)
		os.Exit(-1)
	}
	//model.ModelInit()
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}
	DB.SingularTable(true)
}
