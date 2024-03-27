package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/make-money-fast/v2ray/pkg/vars"
)

func ServerReload(ctx *gin.Context) {
	uri := fmt.Sprintf("http://localhost%s/server/config.json", vars.HttpPort)
	if err := configs.Reload(uri); err != nil {
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
