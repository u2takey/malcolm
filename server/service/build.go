package service

import (
	"github.com/gin-gonic/gin"
	_ "github.com/u2takey/malcolm/model"

	req "github.com/u2takey/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// build
// trigger build
func (s *Service) PostBuild(c *gin.Context) {
	req.Entry(c).Debug("PostBuild")
	jobid := c.Param("jobid")
	ctx := req.Context(c)
	err := s.jobm.RunJob(ctx, jobid)
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
	}
	R(c, gin.H{})
}

// pause/continue/stop build
func (s *Service) PatchBuild(c *gin.Context) {
	req.Entry(c).Debug("PatchBuild")
	c.String(200, "pong")
}

// get build
func (s *Service) GetBuild(c *gin.Context) {
	req.Entry(c).Debug("GetBuild")
	c.String(200, "pong")
}
