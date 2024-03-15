package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	helpers2 "github.com/make-money-fast/v2ray/helpers"
	system2 "github.com/make-money-fast/v2ray/system"
	"github.com/skip2/go-qrcode"
)

func ServerIndex(ctx *gin.Context) {
	cfg := system2.LoadServerConfigStruct()
	ip := helpers2.GetMyIP()
	uuid := cfg.GetUUID()
	port := cfg.GetPort()
	net := cfg.GetNetwork()
	ctx.HTML(200, "server_index.gohtml", gin.H{
		"configPath": helpers2.GetConfigPath(),
		"config":     cfg,
		"isRunning":  system2.IsRunning(),
		"ip":         helpers2.GetMyIP(),
		"vmess":      helpers2.VMessLink(port, net, uuid, ip),
		"port":       port,
		"uuid":       uuid,
		"configJSON": cfg.GetIntentJSON(),
		"version":    system2.Version,
		"token":      helpers2.Username + helpers2.Password,
	})
}

func ServeQRCode(ctx *gin.Context) {
	cfg := system2.LoadServerConfigStruct()
	ip := helpers2.GetMyIP()
	uuid := cfg.GetUUID()
	port := cfg.GetPort()
	net := cfg.GetNetwork()
	vmess := helpers2.VMessLink(port, net, uuid, ip)
	data, _ := qrcode.Encode(vmess, qrcode.High, 300)
	ctx.Data(200, "image/png", data)
}

func ServerStart(ctx *gin.Context) {
	uri := fmt.Sprintf("http://localhost%s/server/config.json", helpers2.HttpPort)
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

func ServerStop(ctx *gin.Context) {
	system2.Stop()
	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

func ServerReload(ctx *gin.Context) {
	uri := fmt.Sprintf("http://localhost%s/server/config.json", helpers2.HttpPort)
	if err := system2.Reload(uri); err != nil {
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

type ConfigChangeRequest struct {
	Config string `json:"config" form:"config" binding:"required"`
}

func ServerConfig(ctx *gin.Context) {
	var req ConfigChangeRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	cfg, err := system2.ServerConfigFromJSON(req.Config)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "配置文件错误:" + err.Error(),
		})
		return
	}

	if err := system2.SaveConfig(cfg); err != nil {
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

func ServerConfigToClientJSON(ctx *gin.Context) {
	cfg := system2.LoadServerConfigStruct()
	clientConfig := cfg.ToTestClientConfig()
	ctx.IndentedJSON(200, clientConfig)
}

func ServerConfigJSON(ctx *gin.Context) {
	cfg := system2.LoadServerConfigStruct()
	ctx.IndentedJSON(200, cfg)
}

func ServerInitDefaultConfig(ctx *gin.Context) {
	if err := system2.SaveConfig(system2.DefaultServerConfig()); err != nil {
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

func GetServerProxyState(ctx *gin.Context) {
	state := system2.GetServerProxyState()
	ctx.JSON(200, gin.H{
		"state": state,
		"code":  0,
	})
}
