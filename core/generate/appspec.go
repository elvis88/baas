package generate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/elvis88/baas/common/util"
	"github.com/spf13/viper"
)

// ApplicationSpec 定义
type ApplicationSpec struct {
	Org             string
	Name            string
	CoinfigFileName string
	Account         string // user
}

// 数据目录路径
func (app *ApplicationSpec) datadir() string {
	dirPath := filepath.Join(app.Name, Application)
	root := filepath.Join(os.Getenv("GOPATH"), viper.GetString("baas.shared"))
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	return filepath.Join(root, "data", dirPath)
}

// 模本目录路径
func (app *ApplicationSpec) templatedir() string {
	root := filepath.Join(os.Getenv("GOPATH"), viper.GetString("baas.shared"))
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	return filepath.Join(root, "template", app.Org, Application)
}

// Build 创建数据目录
func (app *ApplicationSpec) Build(copyTo func(tfilename, filename string) error) error {
	err := util.CreatedDir(app.datadir())
	if err != nil {
		return err
	}
	// 模本目录 ===> 数据目录
	return util.CopyDir(app.templatedir(), app.datadir(), copyTo)
}

// Remove 删除数据目录
func (app *ApplicationSpec) Remove() error {
	return util.RemoveDir(app.datadir())
}

// GetConfig 获取配置内容
func (app *ApplicationSpec) GetConfig() (string, error) {
	config, err := ioutil.ReadFile(app.GetConfigFile())
	return string(config), err
}

// SetConfig 设置配置内容
func (app *ApplicationSpec) SetConfig(config string) error {
	return ioutil.WriteFile(app.GetConfigFile(), []byte(config), os.ModePerm)
}

// GetConfigFile 获取配置文件
func (app *ApplicationSpec) GetConfigFile() string {
	return filepath.Join(app.datadir(), app.CoinfigFileName)
}

// Join 加入该应用
func (app *ApplicationSpec) Join(user string) error {
	cdatadir := app.datadir()
	ndatadir := strings.Replace(cdatadir, filepath.Join(app.Account, app.Name, Application), filepath.Join(user, app.Name, Application), 1)
	if err := util.CreatedDir(ndatadir); err != nil {
		return err
	}

	cfilename := filepath.Join(cdatadir, app.CoinfigFileName)
	nfilename := filepath.Join(ndatadir, app.CoinfigFileName)
	return os.Symlink(cfilename, nfilename)
}
