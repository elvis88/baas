package generate

// FTBrowserSpec ft block chain
type FTBrowserSpec struct {
	ApplicationSpec
}

// NewFTBrowserSpec ft block chain
func NewFTBrowserSpec(account string, name string) *FTBrowserSpec {
	return &FTBrowserSpec{
		ApplicationSpec{
			Name:            name,
			Account:         account,
			CoinfigFileName: FTConfigFileName,
		},
	}
}

// SetConfig 设置配置内容
func (app *FTBrowserSpec) SetConfig(config string) error {
	return nil
}

// FTBrowserDeploySpec ft block chain deployment
type FTBrowserDeploySpec struct {
	ApplicationDeploySpec
}

// NewFTBrowserDeploySpec ft block chain deployment
func NewFTBrowserDeploySpec(account string, name string) *FTBrowserDeploySpec {
	return &FTBrowserDeploySpec{
		ApplicationDeploySpec{
			Name:            name,
			Account:         account,
			CoinfigFileName: FTDeployConfigFileName,
		},
	}
}

// SetConfig 设置配置内容
func (app *FTBrowserDeploySpec) SetConfig(config string) error {
	return nil
}
