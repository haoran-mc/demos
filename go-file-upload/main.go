package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func filePrint(file io.Reader) {
	bs, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("read file fail: " + err.Error())
	}
	fmt.Println(string(bs))
}

func process(c echo.Context) error {
	// 1. local file
	filePath := c.Request().FormValue("path")
	if filePath != "" {
		c.JSON(http.StatusOK, "upload file from local")
		return nil
	}

	// 2. POST multi part form
	if c.Request().MultipartForm != nil {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		for _, v := range form.File {
			for _, f := range v {
				file, e := f.Open()
				if e != nil {
					continue
				}
				defer file.Close()
				c.JSON(http.StatusOK, "upload file as MultipartForm")
				filePrint(file)
			}
		}
		return nil
	}

	// 3. POST Blob
	if strings.HasPrefix(c.Request().Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		// TODO 重新发送文件
		c.JSON(http.StatusOK, "upload file by blob")
		return nil
	}

	// 4. PUT Body
	filePrint(c.Request().Body)
	c.JSON(http.StatusOK, "upload file by body")
	return nil
}

func main() {
	e := echo.New()
	e.Server.ReadTimeout = 30 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.POST("/file", process)
	e.PUT("/file", process)
	e.POST("/file/:name", process)
	e.PUT("/file/:name", process)

	e.Start(":7782")
}
