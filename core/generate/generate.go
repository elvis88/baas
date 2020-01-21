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
	Unjoin(user string) error
}

// NewAppSpec an application spec object
func NewAppSpec(user, name, org string) AppSpec {
	lorg := strings.ToLower(org)
	switch lorg {
	case "ft":
		return NewFTSpec(user, name, lorg)
	case "fttest":
		return NewFTSpec(user, name, lorg)
	default:
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
	GetScriptPath() string
}

// NewAppDeploySpec an application spec object
func NewAppDeploySpec(user, name, org, chain string) AppDeploySpec {
	lorg := strings.ToLower(org)
	switch lorg {
	case "ft":
		return NewFTDeploySpec(user, name, lorg, chain)
	case "fttest":
		return NewFTDeploySpec(user, name, lorg, chain)
	default:
	}
	return nil
}
