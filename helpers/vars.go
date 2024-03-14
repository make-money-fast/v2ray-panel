package helpers

import (
	"os"
	"path/filepath"
)

var (
	IsServer   bool
	IsClient   bool
	ConfigPath string
)

var (
	HttpPort       = ":8091"
	ServerTestPort = 30221
)

func GetConfigPath() string {
	path := ConfigPath
	if ConfigPath == "" {
		wd, _ := os.Getwd()
		path = filepath.Join(wd, "config.json")
	}
	return path
}
