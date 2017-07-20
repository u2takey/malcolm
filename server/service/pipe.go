package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"

	"github.com/arlert/malcolm/model"
	"github.com/arlert/malcolm/utils"
	req "github.com/arlert/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// pipe
// create pipe
func (s *Service) PostPipe(c *gin.Context) {
	req.Entry(c).Debug("PostPipe")
	pipe, err := BindValidPipe(c)
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
		return
	}
	pipe.ID = bson.NewObjectId()
	pipe.Created = time.Now()
	pipe.Updated = time.Now()
	err = s.store.Cols.Pipe.Insert(pipe)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	s.pipem.AddPipe(req.Context(c), pipe)
	R(c, pipe)
}

// update pipe
func (s *Service) PatchPipe(c *gin.Context) {
	req.Entry(c).Debug("PatchPipe")
	pipe, err := BindValidPipe(c)
	sel := bson.M{"_id": pipe.ID}
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
		return
	}
	pipe.Updated = time.Now()
	err = s.store.Cols.Pipe.Update(sel, pipe)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	s.pipem.AddPipe(req.Context(c), pipe)
	R(c, pipe)
}

// get pipe
func (s *Service) GetPipe(c *gin.Context) {
	req.Entry(c).Debug("GetPipe")
	pipeid := c.Param("pipeid")
	page, size := utils.GetPaginationParams(c, 100)

	var sel bson.M
	pipes := []model.PipeConfig{}
	if pipeid == "" {
		sel = bson.M{}
	} else if !bson.IsObjectIdHex(pipeid) {
		E(c, ErrInvalidParam.WithMessage("pipeid not valid"))
		return
	} else {
		sel = bson.M{"_id": bson.ObjectIdHex(pipeid)}
	}

	err := s.store.Cols.Pipe.Find(sel).Sort("-created").Skip((page - 1) * size).Limit(size).All(&pipes)
	if err != nil {
		E(c, ErrDB)
		return
	}
	R(c, pipes)
}

// delete pipe
func (s *Service) DeletePipe(c *gin.Context) {
	req.Entry(c).Debug("DeletePipe")
	pipeid := c.Param("pipeid")
	if !bson.IsObjectIdHex(pipeid) {
		E(c, ErrInvalidParam)
		return
	}
	sel := bson.M{"_id": bson.ObjectIdHex(pipeid)}
	err := s.store.Cols.Pipe.Remove(sel)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	s.pipem.RemovePipe(req.Context(c), pipeid)
	R(c, gin.H{})
}

//-----------------------------------------------------------------
// BindPipe
func BindValidPipe(c *gin.Context) (pipe *model.PipeConfig, err error) {
	pipe = &model.PipeConfig{}
	c.BindJSON(pipe)
	err = pipe.Valid()
	return
}
