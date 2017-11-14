package logger

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/config"
)

// ParseLogLevel turns a log string to int type
func ParseLogLevel(level string) (int, error) {
	switch level {
	case "emergency":
		return logs.LevelEmergency, nil
	case "alert":
		return logs.LevelAlert, nil
	case "critical":
		return logs.LevelCritical, nil
	case "error":
		return logs.LevelError, nil
	case "warn":
		return logs.LevelWarning, nil
	case "notice":
		return logs.LevelNotice, nil
	case "info":
		return logs.LevelInformational, nil
	case "debug":
		return logs.LevelDebug, nil
	}

	return logs.LevelInformational, fmt.Errorf("fail to parse log level: %s is not supported", level)
}

// ParseAdapter turns an adapter name from log to beego type
func ParseAdapter(name string) (string, error) {
	switch name {
	case "":
		return "", errors.New("log Adapter should not be empty")
	case "console":
		return logs.AdapterConsole, nil
	case "file":
		return logs.AdapterFile, nil
	case "multifile":
		return logs.AdapterMultiFile, fmt.Errorf("adapter '%s' is not supported", name)
	case "smtp":
		return logs.AdapterMail, fmt.Errorf("adapter '%s' is not supported", name)
	case "conn":
		return logs.AdapterConn, fmt.Errorf("adapter '%s' is not supported", name)
	case "es":
		return logs.AdapterEs, fmt.Errorf("adapter '%s' is not supported", name)
	case "jianliao":
		return logs.AdapterJianLiao, fmt.Errorf("adapter '%s' is not supported", name)
	case "slack":
		return logs.AdapterSlack, fmt.Errorf("adapter '%s' is not supported", name)
	case "alils":
		return logs.AdapterAliLS, fmt.Errorf("adapter '%s' is not supported", name)

	}
	return "", fmt.Errorf("adapter '%s' is invalid", name)
}

// ParseAdapterArgs turns 'level' to human readable string
func ParseAdapterArgs(name string, values map[string]interface{}) string {
	v, ok := values["level"]
	if ok {
		if l, err := ParseLogLevel(v.(string)); err == nil {
			values["level"] = l
		} else {
			values["level"] = logs.LevelInformational
		}
	}

	str, _ := json.Marshal(values)
	return string(str)
}

// ParseLog reads user's config
func ParseLog(cfg config.LogConfig) (string, string, error) {
	for n, v := range cfg {
		name, err := ParseAdapter(n)
		if err != nil {
			return "", "", err
		}
		args := ParseAdapterArgs(name, v)
		return name, args, nil
		// If multiple log adapters were detected, use the first one
	}

	return "", "", errors.New("log is not set")
}

// InitLogger setup the logger
func InitLogger(cfg config.LogConfig) error {
	n, args, err := ParseLog(cfg)
	if err != nil {
		logs.SetLogger(logs.AdapterConsole, fmt.Sprintf("{\"level\": %d}", logs.LevelInformational))
		return errors.New("fail to parse logger, fallback to 'console', the debug level 'info'")
	}

	logs.SetLogger(n, args)
	return nil
}
