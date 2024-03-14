package server

import "github.com/gin-gonic/gin"

func ClientIndex(ctx *gin.Context) {
	ctx.HTML(200, "client-index", gin.H{})
}
