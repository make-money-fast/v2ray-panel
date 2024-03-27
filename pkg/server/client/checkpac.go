package client

import "github.com/gin-gonic/gin"

func CheckPac(ctx *gin.Context) {
	ctx.HTML(200, "check_pac.gohtml", nil)
}
