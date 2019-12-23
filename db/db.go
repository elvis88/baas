package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDb(connStr string) {
	var err error
	DB, err = gorm.Open("mysql", connStr)
	if err != nil {
		_ = fmt.Errorf("connect", err)
		os.Exit(-1)
	}
}
