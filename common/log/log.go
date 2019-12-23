package log

import (
	"os"

	"github.com/op/go-logging"
)

const (
	CRITICAL logging.Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var Log *logging.Logger

// GetLogger 获取日志
func GetLogger(module string, level logging.Level) *logging.Logger {
	Log = logging.MustGetLogger(module)
	var format = logging.MustStringFormatter(
		`[%{level:.4s}] %{module} %{time:2006-01-02 15:04:05} [%{longfunc}] : %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")
	logging.SetBackend(backendLeveled)
	return Log
}

// Secret 隐藏字符
func Secret(s string) interface{} {
	return logging.Redact(s)
}
