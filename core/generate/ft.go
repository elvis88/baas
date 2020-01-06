package generate

import (
	"bytes"
	"fmt"
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

		bname := fi.Name()
		if bname == "deployment.sh" {
			bts = bytes.Replace(bts, []byte(`DEPLOY_USER="admin"`), []byte(fmt.Sprintf(`DEPLOY_USER="%s"`, app.Account)), 1)
			bts = bytes.Replace(bts, []byte(`DEPLOY_NAME="ft"`), []byte(fmt.Sprintf(`DEPLOY_NAME="%s"`, app.Name)), 1)
			bts = bytes.Replace(bts, []byte(`APP_NAME="ft"`), []byte(fmt.Sprintf(`APP_NAME="%s"`, app.Org)), 1)
		}

		return ioutil.WriteFile(tname, bts, os.ModePerm)
	}
	return app.ApplicationDeploySpec.Build(copyto)
}
