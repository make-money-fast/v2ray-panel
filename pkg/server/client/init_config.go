package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

func ClientInitDefaultConfig(ctx *gin.Context) {
	if err := configs.SaveClientConfig(configs.DefaultClientConfig()); err != nil {
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
