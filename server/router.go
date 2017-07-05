package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/u2takey/malcolm/server/middleware/header"
)

// Load loads the router
func Load(middleware ...gin.HandlerFunc) http.Handler {

	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(header.NoCache)
	e.Use(header.Options)
	e.Use(header.Version)
	e.Use(middleware...)

	e.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return e
}
