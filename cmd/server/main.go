package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/helpers"
	"github.com/make-money-fast/v2ray/server"
	"github.com/make-money-fast/v2ray/system"
)

func init() {
	flag.BoolVar(&helpers.IsServer, "server", true, "is server")
	flag.StringVar(&helpers.HttpPort, "http-port", ":8091", "HTTP服务端口")
	flag.StringVar(&helpers.Username, "u", "admin", "web端账号")
	flag.StringVar(&helpers.Password, "p", "admin123", "web端密码")
}

func main() {
	flag.Parse()
	fmt.Println("运行模式：【服务端】")
	ip := helpers.GetMyIP()
	fmt.Println(fmt.Sprintf("管理端地址: http://%s:%s@%s%s/server/index", helpers.Username, helpers.Password, ip, helpers.HttpPort))
	gin.SetMode(gin.ReleaseMode)

	go func() {
		fmt.Println("自动启动服务中....")
		system.Start(helpers.GetConfigPath())
	}()
	server.StartServer()
}
