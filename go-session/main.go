package main

import (
	"net/http"

	"github.com/haoran-mc/demos/go-session/session_store"
	"github.com/labstack/echo/v4"
)

type User struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

const (
	name string = "name"
	pass string = "pass"
)

func main() {
	session_store.InitCookieStore()
	session_store.InitFilesystemStore()

	e := echo.New()

	e.POST("/user", func(c echo.Context) error {
		user := User{}
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid request",
			})
		}

		session, _ := session_store.GetFileSystemSession(c.Request())
		session.Values[name] = user.Name
		session.Values[pass] = user.Pass
		session.Save(c.Request(), c.Response())

		return c.JSON(http.StatusOK, map[string]any{
			"message":    "ok",
			"session.ID": session.ID,
		})
	})

	e.GET("/user", func(c echo.Context) error {
		session, _ := session_store.GetFileSystemSession(c.Request())
		n := session.Values[name]
		p := session.Values[pass]
		data := session.Values["data"]

		return c.JSON(http.StatusOK, map[string]any{
			"message": "ok",
			"name":    n,
			"pass":    p,
			"data":    data,
			"id":      session.ID,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
