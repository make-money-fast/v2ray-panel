package main

import (
	"flag"
	"github.com/getlantern/systray"
	"github.com/make-money-fast/v2ray/helpers"
	"github.com/make-money-fast/v2ray/server"
	"github.com/make-money-fast/v2ray/system"
	"github.com/make-money-fast/v2ray/web"
	"net/http"
	"os"
	"time"
)

func init() {
	flag.BoolVar(&helpers.IsServer, "server", false, "is server")
	flag.BoolVar(&helpers.IsClient, "client", false, "is client")

	flag.StringVar(&helpers.HttpPort, "http-port", ":8091", "HTTP服务端口")
}

func main() {
	flag.Parse()

	helpers.IsClient = true
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
		go func() {
			time.Sleep(3 * time.Second)
			http.Get("http://localhost" + helpers.HttpPort + "/client/api/start")
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
