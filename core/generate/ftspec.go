package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/elvis88/baas/common/util"
	"github.com/spf13/viper"
)

// FTSpec 定义
type FTSpec struct {
	Account string
	Name    string
}

func (c *FTSpec) workdir() string {
	dirPath := filepath.Join(c.Account, c.Name, "chain")
	root := viper.GetString("baas.shared")
	if root == "" {
		root = filepath.Join(os.Getenv("GOPATH"), "src/github.com/elvis88/baas/shared")
	}
	dirPath = filepath.Join(root, "data", dirPath)
	return dirPath
}

func (c *FTSpec) configfilename() string {
	return filepath.Join(c.workdir(), fmt.Sprintf("%s_%s", c.Name, "genesis.json"))
}

func (c *FTSpec) genscript() error {
	return nil
}

// Build 创建目录
func (c *FTSpec) Build() error {
	return util.CreatedDir(c.workdir())
}

// Remove 删除目录
func (c *FTSpec) Remove() error {
	return util.RemoveDir(c.workdir())
}

// SetConfig 获取
func (c *FTSpec) SetConfig(config string) error {
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
func (c *FTSpec) GetConfig() (string, error) {
	config, err := ioutil.ReadFile(c.configfilename())
	return string(config), err
}
