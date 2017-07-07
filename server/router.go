package server

import (
	"net/http"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"

	"github.com/u2takey/malcolm/model"
	"github.com/u2takey/malcolm/server/middleware/header"
	"github.com/u2takey/malcolm/server/service"
)

// Load loads the router
func Load(cfg *model.Config) http.Handler {

	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(header.NoCache)
	e.Use(header.Options)
	e.Use(header.Version)
	e.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	svc := service.New(cfg)
	e.GET("ping", svc.GetPing)
	v1group := e.Group("/v1")
	{

		v1group.POST("/job", e.NoRoute)          // create
		v1group.PATCH("/job/:jobid", e.NoRoute)  // update
		v1group.GET("/job/:jobid", e.NoRoute)    // get
		v1group.DELETE("/job/:jobid", e.NoRoute) // delete

		v1group.POST("/job/:jobid/build/", e.NoRoute)          // trigger build
		v1group.PATCH("/job/:jobid/build/:buildid", e.NoRoute) // pause/continue/stop build
		v1group.GET("/job/:jobid/build/:buildid", e.NoRoute)   // get build

		v1group.GET("/job/:jobid/build/:buildid/log/:logid", e.NoRoute) // get whole build / single step log

		v1group.GET("/job/:jobid/build/:buildid/message", e.NoRoute) // allow build step send message to master

		v1group.POST("/hook", e.NoRoute) // hook for git repo
	}

	return e
}
