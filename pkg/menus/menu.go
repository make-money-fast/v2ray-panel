package menus

import (
	"github.com/getlantern/systray"
	"github.com/make-money-fast/v2ray/pkg/configs"
	helpers2 "github.com/make-money-fast/v2ray/pkg/helpers"
	"github.com/make-money-fast/v2ray/pkg/server/client"
	"github.com/make-money-fast/v2ray/pkg/vars"
	"os"
	"os/exec"
	"time"
)

func SetMenus() {
	panelCh := systray.AddMenuItem("管理后台", "")
	v2Ch := systray.AddMenuItem("启动v2", "")
	pxyCh := systray.AddMenuItem("设置代理", "")
	cfgCh := systray.AddMenuItem("查看配置", "")
	exitCh := systray.AddMenuItem("退出", "")

	go func() {
		helpers2.UnSetProxy()

		time.Sleep(2 * time.Second)
		// 查看系统有没有自动启动代理
		if configs.IsRunning() {
			v2Ch.Check()
		}
	}()

	for {
		select {
		case <-panelCh.ClickedCh:
			client.OpenBrowser()
		case <-v2Ch.ClickedCh:
			if v2Ch.Checked() {
				configs.Stop()
				v2Ch.Uncheck()
			} else {
				configs.Start("./config.json")
				v2Ch.Check()
			}
		case <-pxyCh.ClickedCh:
			if pxyCh.Checked() {
				helpers2.UnSetProxy()
				pxyCh.Uncheck()
			} else {
				helpers2.SetProxy(vars.GetPacAddress())
				pxyCh.Check()
			}
		case <-cfgCh.ClickedCh:
			cmd := exec.Command("open", "config.json")
			go cmd.Run()
		case <-exitCh.ClickedCh:
			helpers2.UnSetProxy()
			configs.Stop()
			os.Exit(0)
		}
	}
}
