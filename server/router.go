package server

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/helpers"
	"github.com/make-money-fast/v2ray/web"
	"html/template"
	"net/http"
	"strings"
)

func StartServer() {
	g := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	//g.Delims("${{", "}}")
	//g.Static("/static", "web/static")
	//g.LoadHTMLGlob("web/templates/*.gohtml")

	g.Any("/static/*action", func(ctx *gin.Context) {
		http.FileServer(http.FS(web.Static)).ServeHTTP(ctx.Writer, ctx.Request)
	})

	g.SetHTMLTemplate(template.Must(template.New("").Delims("${{", "}}").ParseFS(web.Templates, "templates/*")))

	g.Use(serverClientCheck())

	// 1. 启动服务端.
	server := g.Group("/server")
	server.Use(gin.BasicAuth(gin.Accounts{
		helpers.Username: helpers.Password,
	}))
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
		client.GET("/helper", ClientHelper)
		client.GET("/config.json", ClientConfigJSON)
		client.GET("/vmess", ClientQRCode)

		client.POST("/api/vmess", ClientImportVmess)
		client.POST("/api/config", ClientConfig)
		client.GET("/api/initDefaultConfig", ClientInitDefaultConfig)
		client.GET("/api/start", ClientStart)
		client.GET("/api/stop", ClientStop)
		client.GET("/api/state", ClientState)
		client.POST("/api/check", ClientCheck)
		client.GET("/api/set_proxy", SetSysProxy)
	}

	g.Run(helpers.HttpPort)
}

func serverClientCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if helpers.IsServer && strings.HasPrefix(ctx.Request.RequestURI, "/server") {
			ctx.Next()
			return
		}
		if helpers.IsClient && strings.HasPrefix(ctx.Request.RequestURI, "/client") {
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
