package main

import (
	"os"
	"spotify-playlist-generator/server"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// User middleware for logging
	e.Use(middleware.Logger())
	e.Logger.SetOutput(os.Stdout)

	_ = server.New(e)

	// Start server
	e.Logger.Fatal(e.Start(":8888"))
}
