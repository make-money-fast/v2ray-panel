package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

func ClientState(ctx *gin.Context) {
	state := configs.GetClientProxyState()
	ctx.JSON(200, gin.H{
		"code":  0,
		"state": state,
	})
}
