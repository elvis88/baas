package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/elvis88/baas/common/util"
	"github.com/spf13/viper"
)

// ApplicationDeploySpec 定义
type ApplicationDeploySpec struct {
	Org             string
	Name            string
	CoinfigFileName string
	Account         string
}

// 数据目录路径
func (app *ApplicationDeploySpec) datadir() string {
	dirPath := filepath.Join(app.Name, Deployment)
	root := filepath.Join(os.Getenv("GOPATH"), viper.GetString("baas.shared"))
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	return filepath.Join(root, "data", dirPath)
}

// 模本目录路径
func (app *ApplicationDeploySpec) templatedir() string {
	root := filepath.Join(os.Getenv("GOPATH"), viper.GetString("baas.shared"))
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	return filepath.Join(root, "template", app.Org, Deployment)
}

// Build 创建数据目录
func (app *ApplicationDeploySpec) Build(copyTo func(tfilename, filename string) error) error {
	shellFilePath := app.datadir()
	err := util.CreatedDir(shellFilePath)
	if err != nil {
		return err
	}
	// 模本目录 ===> 数据目录
	if err = util.CopyDir(app.templatedir(), shellFilePath, copyTo); nil != err {
		return err
	}

	// 修改参数
	var oldText []byte
	shellFileName := fmt.Sprintf("%s/%s", shellFilePath, DeploymentFile)
	if oldText, err = ioutil.ReadFile(shellFileName); nil != err {
		return err
	}

	modifiedText := fmt.Sprintf(string(oldText), app.Account, app.Name)
	if err = ioutil.WriteFile(shellFileName, []byte(modifiedText), 0644); nil != err {
		return err
	}

	return nil
}

// Remove 删除数据目录
func (app *ApplicationDeploySpec) Remove() error {
	return util.RemoveDir(app.datadir())
}

// GetConfig 获取配置内容
func (app *ApplicationDeploySpec) GetConfig() (string, error) {
	cfilename := filepath.Join(app.datadir(), app.CoinfigFileName)
	config, err := ioutil.ReadFile(cfilename)
	return string(config), err
}

// SetConfig 设置配置内容
func (app *ApplicationDeploySpec) SetConfig(config string) error {
	cfilename := filepath.Join(app.datadir(), app.CoinfigFileName)
	return ioutil.WriteFile(cfilename, []byte(config), os.ModePerm)
}

// GetConfigFile 获取配置文件
func (app *ApplicationDeploySpec) GetConfigFile() string {
	return filepath.Join(app.datadir(), app.CoinfigFileName)
}

// GetDeployFile 获取配置文件
func (app *ApplicationDeploySpec) GetDeployFile() string {
	return fmt.Sprintf("/file/%s/%s/%s", app.Name, Deployment, DeploymentFile)
	//return filepath.Join(app.datadir(), DeploymentFile)
}
