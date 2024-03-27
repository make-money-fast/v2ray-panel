package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

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

	cfg, err := configs.LoadClientConfig()
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

	if configs.CheckPorxy(fmt.Sprintf("http://localhost:%d", httpPort), req.Url) {
		ctx.JSON(200, gin.H{
			"code": 0,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": -1,
	})
}
