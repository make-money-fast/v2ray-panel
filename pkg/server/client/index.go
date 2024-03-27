package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/make-money-fast/v2ray/pkg/vars"
)

type ClientData struct {
	ConfigPath string                `json:"configPath"`
	Config     *configs.ClientConfig `json:"config"`
	SrvAddr    string                `json:"srvAddr"`
	IsRunning  bool                  `json:"isRunning"`
	ConfigJSON string                `json:"configJSON"`
	Version    string                `json:"version"`
	Vmess      string                `json:"vmess"`
	HttpProxy  string                `json:"httpProxy"`
	SocksProxy string                `json:"socksPorxy"`
	GitTips    string                `json:"gitTips"`
}

func ClientIndex(ctx *gin.Context) {
	cfg, _ := configs.LoadClientConfig()
	httpPxy, socksPxy := cfg.GetProxy()
	addr, port := cfg.GetServerAddress()

	var data = ClientData{
		ConfigPath: vars.GetConfigPath(),
		Config:     cfg,
		SrvAddr:    fmt.Sprintf("%s:%d", addr, port),
		IsRunning:  configs.IsRunning(),
		ConfigJSON: cfg.GetIntentJSON(),
		Version:    configs.Version,
		Vmess:      cfg.GetVmess(),
		HttpProxy:  httpPxy,
		SocksProxy: socksPxy,
		GitTips:    fmt.Sprintf(gitTips, httpPxy, httpPxy),
	}

	ctx.HTML(200, "client_index.gohtml", data)
}
