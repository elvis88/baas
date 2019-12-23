package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/elvis88/baas/common/util"
	"github.com/spf13/viper"
)

// FTDeploySpec 定义
type FTDeploySpec struct {
	Account string
	Name    string
}

func (c *FTDeploySpec) workdir() string {
	dirPath := filepath.Join(c.Account, c.Name, "deploy")
	root := viper.GetString("baas.shared")
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	dirPath = filepath.Join(root, "data", dirPath)
	return dirPath
}

func (c *FTDeploySpec) configfilename() string {
	return filepath.Join(c.workdir(), fmt.Sprintf("%s_%s", c.Name, "config.yaml"))
}

func (c *FTDeploySpec) genscript() error {
	startcontent := ""
	stopcontent := ""
	restartcontent := ""

	start := filepath.Join(c.workdir(), fmt.Sprintf("%s_%s", c.Name, "start.sh"))
	if err := ioutil.WriteFile(start, []byte(startcontent), os.ModePerm); err != nil {
		return err
	}

	stop := filepath.Join(c.workdir(), fmt.Sprintf("%s_%s", c.Name, "stop.sh"))
	if err := ioutil.WriteFile(stop, []byte(stopcontent), os.ModePerm); err != nil {
		return err
	}

	restart := filepath.Join(c.workdir(), fmt.Sprintf("%s_%s", c.Name, "restart.sh"))
	if err := ioutil.WriteFile(restart, []byte(restartcontent), os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Build 创建目录
func (c *FTDeploySpec) Build() error {
	return util.CreatedDir(c.workdir())
}

// Remove 删除目录
func (c *FTDeploySpec) Remove() error {
	return util.RemoveDir(c.workdir())
}

// SetConfig 获取
func (c *FTDeploySpec) SetConfig(config string) error {
	cfilename := c.configfilename()
	if err := ioutil.WriteFile(cfilename, []byte(config), os.ModePerm); err != nil {
		return err
	}

	if err := c.genscript(); err != nil {
		return err
	}

	return nil
}

// GetConfig 获取
func (c *FTDeploySpec) GetConfig() (string, error) {
	config, err := ioutil.ReadFile(c.configfilename())
	return string(config), err
}
