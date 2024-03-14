package server

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	helpers2 "github.com/make-money-fast/v2ray/helpers"
	system2 "github.com/make-money-fast/v2ray/system"
	"github.com/pkg/browser"
	"github.com/skip2/go-qrcode"
)

//go:embed gittips.md
var gitTips string

//go:embed helper.md
var helper string

func ClientIndex(ctx *gin.Context) {
	cfg, _ := system2.LoadClientConfig()
	httpPxy, socksPxy := cfg.GetProxy()
	ctx.HTML(200, "client_index.gohtml", gin.H{
		"configPath": helpers2.GetConfigPath(),
		"config":     cfg,
		"isRunning":  system2.IsRunning(),
		"configJSON": cfg.GetIntentJSON(),
		"version":    system2.Version,
		"vmess":      cfg.GetVmess(),
		"httpProxy":  httpPxy,
		"socksPorxy": socksPxy,
		"gitTips":    fmt.Sprintf(gitTips, httpPxy, httpPxy),
	})
}

type VmessImportRequest struct {
	Vmess string `json:"vmess" binding:"required"`
}

func ClientImportVmess(ctx *gin.Context) {
	var req VmessImportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	vmess, err := helpers2.FromVmess(req.Vmess)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "解析vmess失败:" + err.Error(),
		})
		return
	}
	cfg := system2.ClientConfigFromVmess(vmess)

	if err := system2.SaveClientConfig(cfg); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "保存配置失败:" + err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func ClientConfigJSON(ctx *gin.Context) {
	cfg, _ := system2.LoadClientConfig()
	ctx.IndentedJSON(200, cfg)
}

func ClientStart(ctx *gin.Context) {
	//uri := fmt.Sprintf("http://localhost%s/client/config.json", helpers2.HttpPort)
	if err := system2.Start("./config.json"); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "操作失败:" + err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func ClientStop(ctx *gin.Context) {
	system2.Stop()
	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func ClientState(ctx *gin.Context) {
	state := system2.GetClientProxyState()
	ctx.JSON(200, gin.H{
		"code":  0,
		"state": state,
	})
}

func ClientQRCode(ctx *gin.Context) {
	cfg, _ := system2.LoadClientConfig()
	if cfg != nil {
		data, _ := qrcode.Encode(cfg.GetVmess(), qrcode.High, 300)
		ctx.Data(200, "image/png", data)
		return
	}
	ctx.String(404, "")
}

type CheckRequest struct {
	Url string `json:"url" binding:"required"`
}

func ClientCheck(ctx *gin.Context) {
	var req CheckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, gin.H{
			"code":  -1,
			"state": "参数错误:" + err.Error(),
		})
		return
	}

	cfg, err := system2.LoadClientConfig()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":  -1,
			"state": "系统错误:" + err.Error(),
		})
		return
	}
	var httpPort int
	for _, inbound := range cfg.Inbounds {
		if inbound.Tag == "http" {
			httpPort = inbound.Port
		}
	}

	if system2.CheckPorxy(fmt.Sprintf("http://localhost:%d", httpPort), req.Url) {
		ctx.JSON(200, gin.H{
			"code": 0,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": -1,
	})
}

func OpenBrowser() {
	browser.OpenURL(fmt.Sprintf("http://localhost%s/client/index", helpers2.HttpPort))
}

func ClientConfig(ctx *gin.Context) {
	var req ConfigChangeRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	cfg, err := system2.ClientConfigFromJSON(req.Config)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "配置文件错误:" + err.Error(),
		})
		return
	}

	if err := system2.SaveClientConfig(cfg); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "保存文件失败:" + err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func ClientInitDefaultConfig(ctx *gin.Context) {
	if err := system2.SaveClientConfig(system2.DefaultClientConfig()); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "保存文件失败:" + err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func SetSysProxy(ctx *gin.Context) {
	state := ctx.Query("state")
	cfg, _ := system2.LoadClientConfig()
	httpPxy, _ := cfg.GetProxy()

	switch state {
	case "on":
		if httpPxy == "" {
			ctx.JSON(200, gin.H{
				"code": -1,
				"msg":  "服务器配置异常",
			})
		}
		if err := helpers2.SetProxy(httpPxy); err != nil {
			ctx.JSON(200, gin.H{
				"code": -1,
				"msg":  "操作失败: " + err.Error(),
			})
		}
	case "off":
		if err := helpers2.UnSetProxy(); err != nil {
			ctx.JSON(200, gin.H{
				"code": -1,
				"msg":  "操作失败: " + err.Error(),
			})
		}
	}

	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func ClientHelper(ctx *gin.Context) {
	ctx.HTML(200, "client_helper.gohtml", gin.H{
		"helper": helper,
	})
}
