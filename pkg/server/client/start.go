package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

func ClientStart(ctx *gin.Context) {
	if err := configs.Start("./config.json"); err != nil {
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
