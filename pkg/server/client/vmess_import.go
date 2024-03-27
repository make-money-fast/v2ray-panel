package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/make-money-fast/v2ray/pkg/helpers"
)

type VmessImportRequest struct {
	Vmess string `json:"vmess" binding:"required"`
}

func ClientImportVmess(ctx *gin.Context) {
	var req VmessImportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	vmess, err := helpers.FromVmess(req.Vmess)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "解析vmess失败:" + err.Error(),
		})
		return
	}
	cfg := configs.ClientConfigFromVmess(vmess)

	if err := configs.SaveClientConfig(cfg); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "保存配置失败:" + err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
	})
}
