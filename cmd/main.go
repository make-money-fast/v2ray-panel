package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/helpers"
	"github.com/make-money-fast/v2ray/server"
	"github.com/make-money-fast/v2ray/system"
	"github.com/make-money-fast/v2ray/web"
	"os"
)

func init() {
	flag.BoolVar(&helpers.IsServer, "server", false, "is server")
	flag.BoolVar(&helpers.IsClient, "client", true, "is client")

	flag.StringVar(&helpers.HttpPort, "http-port", ":8091", "HTTP服务端口")
}

func main() {
	flag.Parse()

	if helpers.IsServer {
		fmt.Println("运行模式：【服务端】")
	} else {
		fmt.Println("运行模式：【客户端】")
	}

	gin.SetMode(gin.ReleaseMode)

	fmt.Println("自动启动服务中....")
	system.Start(helpers.GetConfigPath())

	if helpers.IsClient {
		mainGUI()
	} else {
		server.StartServer()
	}
}

func mainGUI() {
	systray.Run(func() {
		helpers.UnSetProxy()

		go func() {
			server.StartServer()
		}()

		ico, _ := web.Static.ReadFile("static/favicon.ico")

		systray.SetTemplateIcon(ico, ico)
		systray.SetTitle("v2")
		systray.SetTooltip("v2-client")

		webBar := systray.AddMenuItem("web", "打开管理面板")
		exitBar := systray.AddMenuItem("exit", "退出")

		for {
			select {
			case <-webBar.ClickedCh:
				server.OpenBrowser()
			case <-exitBar.ClickedCh:
				system.Stop()
				os.Exit(0)
			}
		}
	}, nil)
}
