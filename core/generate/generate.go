package generate

import "strings"

// AppSpec application spec
type AppSpec interface {
	Build() error
	Remove() error
	GetConfig() (string, error)
	SetConfig(config string) error
	GetConfigFile() string
	Join(user string) error
}

// NewAppSpec an application spec object
func NewAppSpec(user, name, org string) AppSpec {
	lorg := strings.ToLower(org)
	if lorg == "ft" {
		return NewFTSpec(user, name, org)
	}
	return nil
}

// AppDeploySpec application spec
type AppDeploySpec interface {
	Build() error
	Remove() error
	GetConfig() (string, error)
	SetConfig(config string) error
	GetConfigFile() string
	GetDeployFile() string
}

// NewAppDeploySpec an application spec object
func NewAppDeploySpec(user, name, org string) AppDeploySpec {
	lorg := strings.ToLower(org)
	if lorg == "ft" {
		return NewFTDeploySpec(user, name, org)
	}
	return nil
}
