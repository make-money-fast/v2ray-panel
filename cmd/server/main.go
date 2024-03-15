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
}

func main() {
	flag.Parse()
	fmt.Println("运行模式：【服务端】")
	gin.SetMode(gin.ReleaseMode)

	go func() {
		fmt.Println("自动启动服务中....")
		system.Start(helpers.GetConfigPath())
	}()
	server.StartServer()
}
