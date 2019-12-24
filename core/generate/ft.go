package generate

import (
	"io/ioutil"
	"os"
	"strings"
)

// FTSpec ft block chain
type FTSpec struct {
	ApplicationSpec
}

// NewFTSpec ft block chain
func NewFTSpec(account string, name string) *FTSpec {
	return &FTSpec{
		ApplicationSpec{
			Name:            name,
			Account:         account,
			CoinfigFileName: FTConfigFileName,
		},
	}
}

// SetConfig 设置配置内容
func (app *FTSpec) SetConfig(config string) error {
	copyto := func(fname, tname string) error {
		fi, err := os.Stat(fname)
		if err != nil {
			return err
		}
		if strings.Compare(fi.Name(), app.CoinfigFileName) == 0 {
			return nil
		}

		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(tname, bts, os.ModePerm)
	}
	return app.ApplicationSpec.SetConfig(config, copyto)
}

// FTDeploySpec ft block chain deployment
type FTDeploySpec struct {
	ApplicationDeploySpec
}

// NewFTDeploySpec ft block chain deployment
func NewFTDeploySpec(account string, name string) *FTDeploySpec {
	return &FTDeploySpec{
		ApplicationDeploySpec{
			Name:            name,
			Account:         account,
			CoinfigFileName: FTDeployConfigFileName,
		},
	}
}

// SetConfig 设置配置内容
func (app *FTDeploySpec) SetConfig(config string) error {
	copyto := func(fname, tname string) error {
		fi, err := os.Stat(fname)
		if err != nil {
			return err
		}
		if strings.Compare(fi.Name(), app.CoinfigFileName) == 0 {
			return nil
		}

		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(tname, bts, os.ModePerm)
	}
	return app.ApplicationDeploySpec.SetConfig(config, copyto)
}
