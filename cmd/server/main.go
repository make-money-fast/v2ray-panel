package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	helpers2 "github.com/make-money-fast/v2ray/pkg/helpers"
	"github.com/make-money-fast/v2ray/pkg/server"
	"github.com/make-money-fast/v2ray/pkg/vars"
)

func init() {
	flag.BoolVar(&vars.IsServer, "server", true, "is server")
	flag.StringVar(&vars.HttpPort, "http-port", ":8091", "HTTP服务端口")
	flag.StringVar(&vars.Username, "u", "admin", "web端账号")
	flag.StringVar(&vars.Password, "p", "admin123", "web端密码")
}

func main() {
	flag.Parse()
	fmt.Println("运行模式：【服务端】")
	ip := helpers2.GetMyIP()
	fmt.Println(fmt.Sprintf("管理端地址: http://%s:%s@%s%s/server/index", vars.Username, vars.Password, ip, vars.HttpPort))
	gin.SetMode(gin.ReleaseMode)

	go func() {
		fmt.Println("自动启动服务中....")
		configs.Start(vars.GetConfigPath())
	}()
	server.StartServer()
}
