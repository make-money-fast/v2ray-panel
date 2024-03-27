package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

func ServerConfigToClientJSON(ctx *gin.Context) {
	cfg := configs.LoadServerConfigStruct()
	clientConfig := cfg.ToTestClientConfig()
	ctx.IndentedJSON(200, clientConfig)
}

func ServerConfigJSON(ctx *gin.Context) {
	cfg := configs.LoadServerConfigStruct()
	ctx.IndentedJSON(200, cfg)
}

func ServerInitDefaultConfig(ctx *gin.Context) {
	if err := configs.SaveConfig(configs.DefaultServerConfig()); err != nil {
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
