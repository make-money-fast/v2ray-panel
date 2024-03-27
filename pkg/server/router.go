package server

import (
	"github.com/gin-gonic/gin"
	client2 "github.com/make-money-fast/v2ray/pkg/server/client"
	server2 "github.com/make-money-fast/v2ray/pkg/server/server"
	"github.com/make-money-fast/v2ray/pkg/vars"
	"github.com/make-money-fast/v2ray/web"
	"html/template"
	"net/http"
	"strings"
)

func StartServer() {
	g := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	g.Delims("${{", "}}")
	//g.Static("/static", "web/static")
	//g.LoadHTMLGlob("web/templates/*.gohtml")
	//
	g.Any("/static/*action", func(ctx *gin.Context) {
		http.FileServer(http.FS(web.Static)).ServeHTTP(ctx.Writer, ctx.Request)
	})

	g.SetHTMLTemplate(template.Must(template.New("").Delims("${{", "}}").ParseFS(web.Templates, "templates/*")))

	if vars.IsServer {
		g.Use(serverClientCheck())
		// 1. 启动服务端.
		server := g.Group("/server")
		{
			basic := server.Group("/")
			basic.Use(gin.BasicAuth(gin.Accounts{
				vars.Username: vars.Password,
			}))
			{
				basic.GET("/index", server2.ServerIndex)
			}

			server.GET("/config.json", server2.ServerConfigJSON)
			server.GET("/client.json", server2.ServerConfigToClientJSON)
			server.GET("/vmess", server2.ServeQRCode)
			server.GET("/api/start", server2.ServerStart)
			server.GET("/api/stop", server2.ServerStop)
			server.GET("/api/reload", server2.ServerReload)
			server.POST("/api/config", server2.ServerConfig)
			server.GET("/api/initDefaultConfig", server2.ServerInitDefaultConfig)
			server.GET("/api/state", server2.GetServerProxyState)
		}
	}

	if vars.IsClient {
		// 2. 启动客户端
		client := g.Group("/client")
		{
			client.GET("/index", client2.ClientIndex)
			client.GET("/helper", client2.ClientHelper)
			client.GET("/config.json", client2.ClientConfigJSON)
			client.GET("/vmess", client2.ClientQRCode)

			client.POST("/api/vmess", client2.ClientImportVmess)
			client.POST("/api/config", client2.ClientConfig)
			client.GET("/api/initDefaultConfig", client2.ClientInitDefaultConfig)
			client.GET("/api/start", client2.ClientStart)
			client.GET("/api/stop", client2.ClientStop)
			client.GET("/api/state", client2.ClientState)
			client.POST("/api/check", client2.ClientCheck)
			client.GET("/api/set_proxy", client2.SetSysProxy)
			client.GET("/pac.js", client2.Pacjs)
			client.GET("/checkpac", client2.CheckPac)
		}
	}

	g.Run(vars.HttpPort)
}

func serverClientCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if vars.IsServer && strings.HasPrefix(ctx.Request.RequestURI, "/server") {
			ctx.Next()
			return
		}
		if vars.IsClient && strings.HasPrefix(ctx.Request.RequestURI, "/client") {
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
