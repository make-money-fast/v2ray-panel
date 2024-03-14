package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	helpers2 "github.com/make-money-fast/v2ray/helpers"
	system2 "github.com/make-money-fast/v2ray/system"
)

func ClientIndex(ctx *gin.Context) {
	cfg, _ := system2.LoadClientConfig()
	ctx.HTML(200, "client_index.gohtml", gin.H{
		"configPath": helpers2.GetConfigPath(),
		"config":     cfg,
		"isRunning":  system2.IsRunning(),
		"configJSON": cfg.GetIntentJSON(),
		"version":    system2.Version,
		"vmess":      cfg.GetVmess(),
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

func ClientConfig(ctx *gin.Context) {
	cfg, _ := system2.LoadClientConfig()
	ctx.IndentedJSON(200, cfg)
}

func ClientStart(ctx *gin.Context) {
	uri := fmt.Sprintf("http://localhost%s/client/config.json", helpers2.HttpPort)
	if err := system2.Start(uri); err != nil {
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
