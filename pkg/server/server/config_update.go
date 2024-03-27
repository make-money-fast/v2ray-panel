package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
)

type ConfigChangeRequest struct {
	Config string `json:"config" form:"config" binding:"required"`
}

func ServerConfig(ctx *gin.Context) {
	var req ConfigChangeRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	cfg, err := configs.ServerConfigFromJSON(req.Config)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "配置文件错误:" + err.Error(),
		})
		return
	}

	if err := configs.SaveConfig(cfg); err != nil {
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
