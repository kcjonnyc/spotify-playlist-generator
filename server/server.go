package server

import (
    "net/http"

    "github.com/labstack/echo"
)

type Server struct {
    e   *echo.Echo
}

type Error struct {
    Error       string `json:"error"`
}

type Message struct {
    Message     string `json:"message"`
}

func New(e *echo.Echo) (s *Server) {
    s = &Server{
		e:    e,
	}

    // Routes
	e.GET("/status", s.health)
	e.POST("/playlists", s.generatePlaylist)

    return
}

// Simple health check endpoint
func (s *Server) health(c echo.Context) error {
	return c.JSON(http.StatusOK, "Looking good, up and running :)")
}
