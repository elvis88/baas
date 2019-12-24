package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/elvis88/baas/common/util"
	"github.com/spf13/viper"
)

// ApplicationDeploySpec 定义
type ApplicationDeploySpec struct {
	Name            string
	CoinfigFileName string
	Account         string
}

// 程序名
func (app *ApplicationDeploySpec) basename() string {
	return strings.Split(app.Name, "_")[0]
}

// 数据目录路径
func (app *ApplicationDeploySpec) datadir() string {
	dirPath := filepath.Join(app.Account, app.Name, Deployment)
	root := viper.GetString("baas.shared")
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	return filepath.Join(root, "data", dirPath)
}

// 模本目录路径
func (app *ApplicationDeploySpec) templatedir() string {
	root := viper.GetString("baas.shared")
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	return filepath.Join(root, "template", app.basename(), Deployment)
}

// 配置文件名路径
func (app *ApplicationDeploySpec) fullconfigfilename() string {
	return filepath.Join(app.datadir(), fmt.Sprintf("%s_%s", app.Name, app.CoinfigFileName))
}

// Build 创建数据目录
func (app *ApplicationDeploySpec) Build() error {
	return util.CreatedDir(app.datadir())
}

// Remove 删除数据目录
func (app *ApplicationDeploySpec) Remove() error {
	return util.RemoveDir(app.datadir())
}

// GetConfig 获取配置内容
func (app *ApplicationDeploySpec) GetConfig() (string, error) {
	cfilename := app.fullconfigfilename()
	if !util.Exists(cfilename) {
		cfilename = filepath.Join(app.templatedir(), app.CoinfigFileName)
	}
	config, err := ioutil.ReadFile(cfilename)
	return string(config), err
}

// SetConfig 设置配置内容
func (app *ApplicationDeploySpec) SetConfig(config string, copyTo func(tfilename, filename string) error) error {
	cfilename := app.fullconfigfilename()
	if err := ioutil.WriteFile(cfilename, []byte(config), os.ModePerm); err != nil {
		return err
	}
	if copyTo == nil {
		return nil
	}
	return util.CopyDir(app.templatedir(), app.datadir(), copyTo)
}
