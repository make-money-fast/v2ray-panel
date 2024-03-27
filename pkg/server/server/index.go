package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/make-money-fast/v2ray/pkg/helpers"
	"github.com/make-money-fast/v2ray/pkg/vars"
)

func ServerIndex(ctx *gin.Context) {
	cfg := configs.LoadServerConfigStruct()
	ip := helpers.GetMyIP()
	uuid := cfg.GetUUID()
	port := cfg.GetPort()
	net := cfg.GetNetwork()
	ctx.HTML(200, "server_index.gohtml", gin.H{
		"configPath": vars.GetConfigPath(),
		"config":     cfg,
		"isRunning":  configs.IsRunning(),
		"ip":         helpers.GetMyIP(),
		"vmess":      helpers.VMessLink(port, net, uuid, ip),
		"port":       port,
		"uuid":       uuid,
		"configJSON": cfg.GetIntentJSON(),
		"version":    configs.Version,
		"token":      vars.Username + vars.Password,
	})
}
