package log

import (
	"os"

	logging "github.com/op/go-logging"
)

const (
	CRITICAL logging.Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

// GetLogger 获取日志
func GetLogger(module string, level logging.Level) *logging.Logger {
	var log = logging.MustGetLogger(module)
	var format = logging.MustStringFormatter(
		`[%{level:.4s}] %{module} %{time:2006-01-02 15:04:05} [%{longfunc}] : %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")
	logging.SetBackend(backendLeveled)
	return log
}

// Secret 隐藏字符
func Secret(s string) interface{} {
	return logging.Redact(s)
}
