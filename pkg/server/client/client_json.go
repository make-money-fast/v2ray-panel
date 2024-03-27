package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

func ClientConfigJSON(ctx *gin.Context) {
	cfg, _ := configs.LoadClientConfig()
	ctx.IndentedJSON(200, cfg)
}
