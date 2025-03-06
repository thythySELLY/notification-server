package main

import (
	"net/http"

	"notification-server/config"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize MongoDB connection
	config.LoadEnv()
	config.InitMongoDB()
	config.InitRedis()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
