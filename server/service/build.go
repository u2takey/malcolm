package service

import (
	_ "github.com/arlert/malcolm/model"
	"github.com/gin-gonic/gin"

	req "github.com/arlert/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// build
// trigger build
func (s *Service) PostBuild(c *gin.Context) {
	req.Entry(c).Debug("PostBuild")
	pipeid := c.Param("pipeid")
	ctx := req.Context(c)
	err := s.pipem.RunPipe(ctx, pipeid)
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
