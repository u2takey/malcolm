package service

import (
	"github.com/arlert/malcolm/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

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
	build, err := s.pipem.BuildAction(ctx, "", pipeid, model.ActionStart)
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
	query := c.Query("action")
	pipeid := c.Param("pipeid")
	buildid := c.Param("buildid")
	action := model.BuildAction(query)
	ctx := req.Context(c)
	if action == model.ActionPause || action == model.ActionResume || action == model.ActionStop {
		build, err := s.pipem.BuildAction(ctx, buildid, pipeid, action)
		if err != nil {
			c.AbortWithError(400, err)
		} else {
			R(c, build)
		}
	} else {
		c.AbortWithStatus(400)
	}
}

func (s *Service) GetBuildInQueue(c *gin.Context) {
	req.Entry(c).Debug("GetBuildInQueue")
	ctx := req.Context(c)
	builds, err := s.pipem.BuildQueue(ctx)
	if err != nil {
		req.Entry(c).Error(err)
	}
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
	singleBuild := false
	if !bson.IsObjectIdHex(pipeid) {
		E(c, ErrInvalidParam.WithMessage("pipeid not valid"))
		return
	} else {
		sel["pipeid"] = bson.ObjectIdHex(pipeid)
	}
	if bson.IsObjectIdHex(buildid) {
		singleBuild = true
		sel["_id"] = bson.ObjectIdHex(buildid)
	}
	err := s.store.Cols.Build.Find(sel).Sort("-created").Skip((page - 1) * size).Limit(size).All(&builds)
	if err != nil || (singleBuild && len(builds) != 1) {
		E(c, ErrDB)
		return
	}
	if singleBuild {
		R(c, builds[0])
	} else {
		R(c, builds)
	}
}
