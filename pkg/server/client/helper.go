package client

import "github.com/gin-gonic/gin"

func ClientHelper(ctx *gin.Context) {
	ctx.HTML(200, "client_helper.gohtml", gin.H{
		"helper": helper,
	})
}
