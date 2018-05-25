package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":8888"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
