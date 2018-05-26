package main

import (
	"os"
	"spotify-playlist-generator/server"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// Echo instance
	e := echo.New()

	// User middleware for logging
	e.Use(middleware.Logger())
	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetLevel(log.DEBUG)

	_ = server.New(e)

	// Start server
	e.Logger.Fatal(e.Start(":8888"))
}
