package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

func ClientStop(ctx *gin.Context) {
	configs.Stop()
	ctx.JSON(200, gin.H{
		"code": 0,
	})
}
