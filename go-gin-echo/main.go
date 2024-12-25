package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	echo := r.Group("/echo")
	{
		echo.GET("/*path", func(ctx *gin.Context) {
			p := ctx.Param("path")
			ctx.String(200, "[GET] Path: %s\n", p)
		})
		echo.POST("/*path", func(ctx *gin.Context) {
			p := ctx.Param("path")
			raw, _ := ctx.GetRawData()
			ctx.String(200, "[POST] Path: %s, Data: %s\n", p, string(raw))
		})
		echo.PUT("/*path", func(ctx *gin.Context) {
			p := ctx.Param("path")
			raw, _ := ctx.GetRawData()
			ctx.String(200, "[PUT] Path: %s, Data: %s\n", p, string(raw))
		})
		echo.DELETE("/*path", func(ctx *gin.Context) {
			p := ctx.Param("path")
			raw, _ := ctx.GetRawData()
			ctx.String(200, "[DELETE] Path: %s, Data: %s\n", p, string(raw))
		})
	}

	redirect := r.Group("/redirect")
	{
		redirect.GET("/hello1", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/redirect/hello2")
		})
		redirect.GET("/hello2", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})
	}
	r.Run(":9520")
}
