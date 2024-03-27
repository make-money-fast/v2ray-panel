package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/make-money-fast/v2ray/pkg/helpers"
	"github.com/skip2/go-qrcode"
)

func ServeQRCode(ctx *gin.Context) {
	cfg := configs.LoadServerConfigStruct()
	ip := helpers.GetMyIP()
	uuid := cfg.GetUUID()
	port := cfg.GetPort()
	net := cfg.GetNetwork()
	vmess := helpers.VMessLink(port, net, uuid, ip)
	data, _ := qrcode.Encode(vmess, qrcode.High, 300)
	ctx.Data(200, "image/png", data)
}
