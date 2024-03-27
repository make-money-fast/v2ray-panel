package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

func GetServerProxyState(ctx *gin.Context) {
	state := configs.GetServerProxyState()
	ctx.JSON(200, gin.H{
		"state": state,
		"code":  0,
	})
}
