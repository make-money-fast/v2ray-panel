package helpers

import (
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	IsServer   bool
	IsClient   bool
	ConfigPath string
	TestingUrl = "https://www.google.com"

	Username string = "admin"
	Password string = "admin123"
)

func init() {
	go func() {
		conn, err := net.DialTimeout("tcp", "www.google.com", 1*time.Second)
		if err != nil {
			TestingUrl = "https://www.baidu.com"
			return
		}
		conn.Close()
	}()
}

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
