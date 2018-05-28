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

	// CORS restricted
	// Allows requests from "http://localhost:3000" origin with POST method
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000,https://playmoji-b7170.firebaseapp.com/"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	_ = server.New(e)

	// Start server
	e.Logger.Fatal(e.StartTLS(":8080", "cert.pem", "key.pem"))
}
