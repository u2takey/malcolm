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
// job
// create job
func (s *Service) PostJob(c *gin.Context) {
	req.Entry(c).Debug("PostJob")
	job, err := BindValidJob(c)
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
		return
	}
	job.ID = ""
	job.Created = time.Now()
	job.Updated = time.Now()
	err = s.store.Cols.Job.Insert(job)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	R(c, job)
}

// update job
func (s *Service) PatchJob(c *gin.Context) {
	req.Entry(c).Debug("PatchJob")
	job, err := BindValidJob(c)
	sel := bson.M{"_id": job.ID}
	if err != nil {
		E(c, ErrInvalidParam.WithMessage(err.Error()))
		return
	}
	job.Updated = time.Now()
	err = s.store.Cols.Job.Update(sel, job)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	R(c, job)
}

// get job
func (s *Service) GetJob(c *gin.Context) {
	req.Entry(c).Debug("GetJob")
	jobid := c.Param("jobid")
	page, size := utils.GetPaginationParams(c, 100)

	var sel bson.M
	jobs := []model.JobConfig{}
	if jobid == "" {
		sel = bson.M{}
	} else if !bson.IsObjectIdHex(jobid) {
		E(c, ErrInvalidParam.WithMessage("jobid not valid"))
	} else {
		sel = bson.M{"_id": bson.ObjectIdHex(jobid)}
	}

	err := s.store.Cols.Job.Find(sel).Sort("-created").Skip((page - 1) * size).Limit(size).All(&jobs)
	if err != nil {
		E(c, ErrDB)
		return
	}
	R(c, jobs)
}

// delete job
func (s *Service) DeleteJob(c *gin.Context) {
	req.Entry(c).Debug("DeleteJob")
	jobid := c.Param("jobid")
	if !bson.IsObjectIdHex(jobid) {
		E(c, ErrInvalidParam)
		return
	}
	sel := bson.M{"_id": bson.ObjectIdHex(jobid)}
	err := s.store.Cols.Job.Remove(sel)
	if err != nil {
		E(c, ErrDB.WithMessage(err.Error()))
		return
	}
	R(c, gin.H{})
}

//-----------------------------------------------------------------
// BindJob
func BindValidJob(c *gin.Context) (job *model.JobConfig, err error) {
	job = &model.JobConfig{}
	c.BindJSON(job)
	err = job.Valid()
	return
}
