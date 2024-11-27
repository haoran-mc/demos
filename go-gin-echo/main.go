package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/*path", func(ctx *gin.Context) {
		p := ctx.Param("path")
		ctx.String(200, "[GET] Path: %s\n", p)
	})
	r.POST("/*path", func(ctx *gin.Context) {
		p := ctx.Param("path")
		raw, _ := ctx.GetRawData()
		ctx.String(200, "[POST] Path: %s, Data: %s\n", p, string(raw))
	})
	r.PUT("/*path", func(ctx *gin.Context) {
		p := ctx.Param("path")
		raw, _ := ctx.GetRawData()
		ctx.String(200, "[PUT] Path: %s, Data: %s\n", p, string(raw))
	})
	r.DELETE("/*path", func(ctx *gin.Context) {
		p := ctx.Param("path")
		raw, _ := ctx.GetRawData()
		ctx.String(200, "[DELETE] Path: %s, Data: %s\n", p, string(raw))
	})
	r.Run(":9520")
}
