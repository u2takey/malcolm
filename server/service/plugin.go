package service

import (
	"github.com/gin-gonic/gin"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/arlert/malcolm/model"
	"github.com/arlert/malcolm/utils"
	req "github.com/arlert/malcolm/utils/reqlog"
)

//-----------------------------------------------------------------
// plugin
// create build
func (s *Service) PostPlugin(c *gin.Context) {
	req.Entry(c).Debug("PostPlugin")
	plugin, err := BindValidPlugin(c)
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
		return
	}
	plugin.Updated = time.Now()
	err = s.store.Cols.Plugin.Insert(plugin)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	R(c, plugin)
}

// update build
func (s *Service) PatchPlugin(c *gin.Context) {
	req.Entry(c).Debug("PatchPlugin")
	plugin, err := BindValidPlugin(c)
	sel := bson.M{"_id": plugin.ID}
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
		return
	}
	plugin.Updated = time.Now()
	err = s.store.Cols.Plugin.Update(sel, plugin)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	R(c, plugin)
}

// get build
func (s *Service) GetPlugin(c *gin.Context) {
	req.Entry(c).Debug("GetPlugin")
	pluginid := c.Param("pluginid")
	page, size := utils.GetPaginationParams(c, 100)

	var sel bson.M
	plugins := []model.Plugin{}
	if pluginid == "" {
		sel = bson.M{}
	} else if !bson.IsObjectIdHex(pluginid) {
		E(c, ErrInvalidParam.WithMessage("pluginid not valid"))
	} else {
		sel = bson.M{"_id": bson.ObjectIdHex(pluginid)}
	}

	err := s.store.Cols.Plugin.Find(sel).Sort("-created").Skip((page - 1) * size).Limit(size).All(&plugins)
	if err != nil {
		E(c, ErrDB)
		return
	}
	R(c, plugins)
}

// delete build
func (s *Service) DeletePlugin(c *gin.Context) {
	req.Entry(c).Debug("DeletePlugin")
	pluginid := c.Param("pluginid")
	if !bson.IsObjectIdHex(pluginid) {
		E(c, ErrInvalidParam)
		return
	}
	sel := bson.M{"_id": bson.ObjectIdHex(pluginid)}
	err := s.store.Cols.Plugin.Remove(sel)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	R(c, gin.H{})
}

//-----------------------------------------------------------------
// BindPipe
func BindValidPlugin(c *gin.Context) (plugin *model.Plugin, err error) {
	plugin = &model.Plugin{}
	c.BindJSON(plugin)
	err = plugin.Valid()
	return
}
