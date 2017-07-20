package service

import (
	"github.com/arlert/malcolm/model"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"

	"github.com/arlert/malcolm/utils"
	req "github.com/arlert/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// build
// trigger build
func (s *Service) PostBuild(c *gin.Context) {
	req.Entry(c).Debug("PostBuild")
	pipeid := c.Param("pipeid")
	ctx := req.Context(c)
	build, err := s.pipem.RunPipe(ctx, pipeid)
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
		return
	}
	s.store.Cols.Build.Insert(build)
	R(c, build)
}

// pause/continue/stop build
func (s *Service) PatchBuild(c *gin.Context) {
	req.Entry(c).Debug("PatchBuild")
	c.String(200, "pong")
}

func (s *Service) GetBuildInQueue(c *gin.Context) {
	req.Entry(c).Debug("GetBuildInQueue")
	builds := s.pipem.Queue()
	R(c, builds)
}

// get build
func (s *Service) GetBuild(c *gin.Context) {
	req.Entry(c).Debug("GetBuild")
	buildid := c.Param("buildid")
	pipeid := c.Param("pipeid")
	page, size := utils.GetPaginationParams(c, 100)
	sel := bson.M{}
	builds := []model.Build{}
	if !bson.IsObjectIdHex(pipeid) {
		E(c, ErrInvalidParam.WithMessage("pipeid not valid"))
		return
	} else {
		sel["pipeid"] = bson.ObjectIdHex(pipeid)
	}
	if bson.IsObjectIdHex(buildid) {
		sel["_id"] = bson.ObjectIdHex(buildid)
	}
	err := s.store.Cols.Build.Find(sel).Sort("-created").Skip((page - 1) * size).Limit(size).All(&builds)
	if err != nil {
		E(c, ErrDB)
		return
	}
	R(c, builds)
}
