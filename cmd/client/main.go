package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/make-money-fast/v2ray/pkg/configs"
	helpers2 "github.com/make-money-fast/v2ray/pkg/helpers"
	"github.com/make-money-fast/v2ray/pkg/menus"
	server2 "github.com/make-money-fast/v2ray/pkg/server"
	"github.com/make-money-fast/v2ray/pkg/vars"
	"github.com/make-money-fast/v2ray/web"
)

func init() {
	flag.StringVar(&vars.HttpPort, "http-port", ":8091", "HTTP服务端口")
}

func main() {
	vars.IsClient = true
	flag.Parse()

	fmt.Println("运行模式：【客户端】")
	fmt.Println("自动启动服务中....")

	go func() {
		configs.Start(vars.GetConfigPath())
	}()

	mainGUI()
}

func mainGUI() {
	systray.Run(func() {
		helpers2.UnSetProxy()

		go func() {
			server2.StartServer()
		}()

		ico, _ := web.Static.ReadFile("static/favicon.png")

		systray.SetTemplateIcon(ico, ico)
		systray.SetTooltip("v2-client")
		menus.SetMenus()
	}, nil)
}
