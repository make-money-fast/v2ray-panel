package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/helpers"
	"strings"
)

func StartServer() {
	g := gin.Default()
	g.Delims("${{", "}}")
	g.Static("/static", "web/static")
	g.LoadHTMLGlob("web/templates/*.gohtml")
	g.Use(serverClientCheck())

	// 1. 启动服务端.
	server := g.Group("/server")
	{
		server.GET("/index", ServerIndex)
		server.GET("/api/start", ServerStart)
		server.GET("/api/stop", ServerStop)
		server.GET("/api/reload", ServerReload)
		server.POST("/api/config", ServerConfig)
		server.GET("/api/initDefaultConfig", ServerInitDefaultConfig)
		server.GET("/config.json", ServerConfigJSON)
		server.GET("/client.json", ServerConfigToClientJSON)
		server.GET("/api/state", GetServerProxyState)

		server.GET("/vmess", ServeQRCode)
	}

	// 2. 启动客户端
	client := g.Group("/client")
	{
		client.GET("/index", ClientIndex)
	}

	g.Run(helpers.HttpPort)
}

func serverClientCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if helpers.IsServer && strings.HasPrefix(ctx.Request.RequestURI, "/server") {
			ctx.Next()
			return
		}
		if helpers.IsServer && strings.HasPrefix(ctx.Request.RequestURI, "/client") {
			ctx.Next()
			return
		}

		if strings.HasPrefix(ctx.Request.RequestURI, "/static") ||
			strings.HasPrefix(ctx.Request.RequestURI, "/favico.ico") {
			ctx.Next()
			return
		}

		ctx.String(200, "invalid url")
		ctx.Abort()
	}
}
