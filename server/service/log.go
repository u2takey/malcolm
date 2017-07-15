package service

import (
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"

	"github.com/arlert/malcolm/model"
	req "github.com/arlert/malcolm/utils/reqlog"
	"k8s.io/kubernetes/pkg/util/flushwriter"
)

//-----------------------------------------------------------------
// log
// get whole build / single step log
func (s *Service) GetLog(c *gin.Context) {
	req.Entry(c).Debug("GetLog")
	buildid := c.Param("buildid")
	pipeid := c.Param("pipeid")
	sel := bson.M{}
	build := &model.Build{}
	if !bson.IsObjectIdHex(pipeid) || !bson.IsObjectIdHex(buildid) {
		E(c, ErrInvalidParam.WithMessage("pipeid or buildid not valid"))
		return
	}
	sel["pipeid"] = bson.ObjectIdHex(pipeid)
	sel["_id"] = bson.ObjectIdHex(buildid)

	err := s.store.Cols.Build.Find(sel).One(build)
	if err != nil || len(build.Works) == 0 {
		E(c, ErrDB)
		return
	}
	c.Header("Transfer-Encoding", "chunked")
	writer := flushwriter.Wrap(c.Writer)
	err = s.logmgr.GetLog(&build.Works[0], writer)
	if err != nil {
		c.AbortWithError(400, err)
	}
}

//-----------------------------------------------------------------
// message
// allow build step send message to master
func (s *Service) GetMessage(c *gin.Context) {
	req.Entry(c).Debug("GetMessage")
	c.String(200, "pong")
}
