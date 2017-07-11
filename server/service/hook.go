package service

import (
	"github.com/gin-gonic/gin"

	req "github.com/u2takey/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// hook
// hook for git repo
func (s *Service) GetHook(c *gin.Context) {
	req.Entry(c).Debug("GetHook")
	c.String(200, "pong")
}
