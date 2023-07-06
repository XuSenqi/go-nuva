package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var _config Config

type (
	Config struct {
		ServiceInfo      ServiceInfo `yaml:"service"`
	}

	ServiceInfo struct {
		Name      string  `yaml:"name"`
		Ip        string  `yaml:"ip"`
		Port      int32   `yaml:"port"`
		Env       string  `yaml:"env"`
		Log       Log     `yaml:"log"`
		StartTime int64
	}

	Log struct {
		LogPath  string `yaml:"log_path"`
		LogLevel string `yaml:"log_level"`
	}
)

func InitConfig(cfgFile string) error {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("read config file: %s", err.Error())
	}
	if err = yaml.Unmarshal(data, &_config); err != nil {
		return err
	}


	return nil
}

func GetListenIp() string {
	return _config.ServiceInfo.Ip
}

func SetListenIp(ip string) {
	_config.ServiceInfo.Ip = ip
}

func SetListenPort(port int32) {
	_config.ServiceInfo.Port = port
}

func GetListenPort() int32 {
	return _config.ServiceInfo.Port
}

func SetEnv(env string) {
	_config.ServiceInfo.Env = env
}

func GetEnv() string {
	return _config.ServiceInfo.Env
}

func GetServiceName() string {
	return _config.ServiceInfo.Name
}

func GetLogPath() string {
	return _config.ServiceInfo.Log.LogPath
}

func GetLogLevel() string {
	return _config.ServiceInfo.Log.LogLevel
}

func SetLogPath(logPath string) {
	_config.ServiceInfo.Log.LogPath = logPath
}

func GetAPIAddr() string {
	return fmt.Sprintf("%s:%d", _config.ServiceInfo.Ip, _config.ServiceInfo.Port)
}