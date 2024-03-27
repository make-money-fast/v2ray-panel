package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/skip2/go-qrcode"
)

func ClientQRCode(ctx *gin.Context) {
	cfg, _ := configs.LoadClientConfig()
	if cfg != nil {
		data, _ := qrcode.Encode(cfg.GetVmess(), qrcode.High, 300)
		ctx.Data(200, "image/png", data)
		return
	}
	ctx.String(404, "")
}
