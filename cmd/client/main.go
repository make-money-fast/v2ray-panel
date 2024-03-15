package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/make-money-fast/v2ray/helpers"
	"github.com/make-money-fast/v2ray/server"
	"github.com/make-money-fast/v2ray/system"
	"github.com/make-money-fast/v2ray/web"
	"os"
	"time"
)

func init() {
	flag.BoolVar(&helpers.IsClient, "client", true, "is client")
	flag.StringVar(&helpers.HttpPort, "http-port", ":8091", "HTTP服务端口")
}

func main() {
	flag.Parse()

	fmt.Println("运行模式：【客户端】")
	fmt.Println("自动启动服务中....")

	go func() {
		time.Sleep(1 * time.Second)
		server.OpenBrowser()
	}()
	go func() {
		system.Start(helpers.GetConfigPath())
	}()

	mainGUI()
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
