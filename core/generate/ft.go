package generate

import (
	"io/ioutil"
	"os"
)

// FTSpec ft block chain
type FTSpec struct {
	ApplicationSpec
}

// NewFTSpec ft block chain
func NewFTSpec(account string, name string, org string) *FTSpec {
	return &FTSpec{
		ApplicationSpec{
			Org:             org,
			Name:            name,
			Account:         account,
			CoinfigFileName: FTConfigFileName,
		},
	}
}

// Build 创建数据目录
func (app *FTSpec) Build() error {
	copyto := func(fname, tname string) error {
		fi, err := os.Stat(fname)
		if err != nil {
			return err
		}
		_ = fi.Name()

		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(tname, bts, os.ModePerm)
	}
	return app.ApplicationSpec.Build(copyto)
}

// FTDeploySpec ft block chain deployment
type FTDeploySpec struct {
	ApplicationDeploySpec
}

// NewFTDeploySpec ft block chain deployment
func NewFTDeploySpec(account string, name string, org string) *FTDeploySpec {
	return &FTDeploySpec{
		ApplicationDeploySpec{
			Org:             org,
			Name:            name,
			Account:         account,
			CoinfigFileName: FTDeployConfigFileName,
		},
	}
}

// Build 创建数据目录
func (app *FTDeploySpec) Build() error {
	copyto := func(fname, tname string) error {
		fi, err := os.Stat(fname)
		if err != nil {
			return err
		}
		_ = fi.Name()

		bts, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(tname, bts, os.ModePerm)
	}
	return app.ApplicationDeploySpec.Build(copyto)
}
