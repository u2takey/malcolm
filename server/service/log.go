package service

import (
	"github.com/gin-gonic/gin"

	req "github.com/arlert/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// log
// get whole build / single step log
func (s *Service) GetLog(c *gin.Context) {
	req.Entry(c).Debug("GetLog")
	c.String(200, "pong")
}

//-----------------------------------------------------------------
// message
// allow build step send message to master
func (s *Service) GetMessage(c *gin.Context) {
	req.Entry(c).Debug("GetMessage")
	c.String(200, "pong")
}
