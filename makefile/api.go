package main

import (
	"os"

	helloworld "github.com/haoran-mc/build_tools/makefile/src/hello-world"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	e := echo.New()
	e.GET("/", helloworld.HelloWorld)
	e.Start(os.Getenv("ADDR"))
}
