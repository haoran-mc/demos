package main

import (
	"net/http"

	"github.com/haoran-mc/go_pkgs/go_web/gee"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "hello gee\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"haoran"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
