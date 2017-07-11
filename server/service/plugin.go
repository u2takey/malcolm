package service

import (
	"github.com/gin-gonic/gin"

	req "github.com/arlert/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// plugin
// create build
func (s *Service) PostPlugin(c *gin.Context) {
	req.Entry(c).Debug("PostPlugin")
	c.String(200, "pong")
}

// update build
func (s *Service) PatchPlugin(c *gin.Context) {
	req.Entry(c).Debug("PatchPlugin")
	c.String(200, "pong")
}

// get build
func (s *Service) GetPlugin(c *gin.Context) {
	req.Entry(c).Debug("GetPlugin")
	c.String(200, "pong")
}

// delete build
func (s *Service) DeletePlugin(c *gin.Context) {
	req.Entry(c).Debug("DeletePlugin")
	c.String(200, "pong")
}
