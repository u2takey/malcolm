package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/arlert/malcolm/model"
	"github.com/arlert/malcolm/server/middleware/header"
	"github.com/arlert/malcolm/server/service"
	_ "github.com/arlert/malcolm/utils/loghook"
	"github.com/arlert/malcolm/utils/reqlog"
)

// Load loads the router
func Load(cfg *model.Config) http.Handler {

	logrus.Debugf("\n\nLoad with config:\n %+v\n\n", cfg)

	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(header.NoCache)
	e.Use(header.Options)
	e.Use(header.Version)
	e.Use(reqlog.ReqLoggerMiddleware(logrus.New(), time.RFC3339, true))

	svc := service.New(cfg)

	e.GET("ping", svc.GetPing)
	v1group := e.Group("/v1")
	{
		//-----------------------------------------------------------------
		// job
		v1group.POST("/job", svc.PostJob)            // create job
		v1group.PATCH("/job/:jobid", svc.PatchJob)   // update job
		v1group.GET("/job", svc.GetJob)              // get job
		v1group.GET("/job/:jobid", svc.GetJob)       // get job
		v1group.DELETE("/job/:jobid", svc.DeleteJob) // delete job

		//-----------------------------------------------------------------
		// build
		v1group.POST("/job/:jobid/build", svc.PostBuild)            // trigger build
		v1group.PATCH("/job/:jobid/build/:buildid", svc.PatchBuild) // pause/continue/stop build
		v1group.GET("/job/:jobid/build/:buildid", svc.GetBuild)     // get build

		//-----------------------------------------------------------------
		// plugin
		v1group.POST("/plugin", svc.PostPlugin)                 // create build
		v1group.PATCH("/plugin/:identifier", svc.PatchPlugin)   // update build
		v1group.GET("/plugin/:identifier", svc.GetPlugin)       // get build
		v1group.DELETE("/plugin/:identifier", svc.DeletePlugin) // delete build

		//-----------------------------------------------------------------
		// log & message
		v1group.GET("/job/:jobid/build/:buildid/log/:logid", svc.GetLog)  // get whole build / single step log
		v1group.GET("/job/:jobid/build/:buildid/message", svc.GetMessage) // allow build step send message to master

		//-----------------------------------------------------------------
		// hook
		v1group.POST("/hook", svc.GetHook) // hook for git repo
	}

	return e
}
