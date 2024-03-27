package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/make-money-fast/v2ray/pkg/helpers"
)

func SetSysProxy(ctx *gin.Context) {
	state := ctx.Query("state")
	cfg, _ := configs.LoadClientConfig()
	httpPxy, _ := cfg.GetProxy()

	switch state {
	case "on":
		if httpPxy == "" {
			ctx.JSON(200, gin.H{
				"code": -1,
				"msg":  "服务器配置异常",
			})
		}
		if err := helpers.SetProxy(httpPxy); err != nil {
			ctx.JSON(200, gin.H{
				"code": -1,
				"msg":  "操作失败: " + err.Error(),
			})
		}
	case "off":
		if err := helpers.UnSetProxy(); err != nil {
			ctx.JSON(200, gin.H{
				"code": -1,
				"msg":  "操作失败: " + err.Error(),
			})
		}
	}

	ctx.JSON(200, gin.H{
		"code": 0,
	})
}
