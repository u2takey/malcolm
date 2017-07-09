package jobmgr

import (
	"errors"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"labix.org/v2/mgo/bson"

	"github.com/u2takey/malcolm/model"
	"github.com/u2takey/malcolm/store"
	// "github.com/u2takey/malcolm/utils"
	req "github.com/u2takey/malcolm/utils/reqlog"
)

type JobMgr struct {
	joblock sync.RWMutex
	jobs    map[string]*model.JobConfig
	store   *store.Store
	engine  *Engine
}

func NewJobMgr(st *store.Store, eng *Engine) (j *JobMgr) {
	j = &JobMgr{
		store:  st,
		engine: eng,
	}
	err := j.ListJob()
	if err != nil {
		logrus.Fatalln("JobMgr.list error -", err)
	}
	return j
}

func (m *JobMgr) ListJob() (err error) {
	jobl := []model.JobConfig{}
	err = m.store.Cols.Job.Find(bson.M{}).All(&jobl)
	if err != nil {
		return
	}
	m.joblock.Lock()
	defer m.joblock.Unlock()
	jobm := make(map[string]*model.JobConfig, 0)
	for _, job := range jobl {
		jobm[job.ID.Hex()] = &job
	}
	m.jobs = jobm
	logrus.Infof("list jobs, count : %d \n", len(jobm))
	return
}

func (m *JobMgr) AddJob(c context.Context, job *model.JobConfig) {
	m.joblock.Lock()
	defer m.joblock.Unlock()
	m.jobs[job.ID.Hex()] = job
}

func (m *JobMgr) RemoveJob(c context.Context, jobid string) {
	m.joblock.Lock()
	defer m.joblock.Unlock()
	if _, ok := m.jobs[jobid]; ok {
		delete(m.jobs, jobid)
	} else {
		req.Entry(c).Infof("remove job  %s warning : notexsit \n", jobid)
	}
}

func (m *JobMgr) ReplaceJob(c context.Context, jobid string, job *model.JobConfig) {
	m.joblock.Lock()
	defer m.joblock.Unlock()
	m.jobs[job.ID.Hex()] = job
	if _, ok := m.jobs[jobid]; !ok {
		req.Entry(c).Infof("replace job  %s warning : notexsit \n", jobid)
	}
}

func (m *JobMgr) RunJob(c context.Context, jobid string) (err error) {
	m.joblock.RLock()
	defer m.joblock.RUnlock()
	if job, ok := m.jobs[jobid]; ok {
		err = m.runJob(c, job, make(chan bool))
	} else {
		req.Entry(c).Infof("RunJob job  %s warning : notexsit \n", jobid)
		err = errors.New("Job not exsit")
	}
	return
}

func (m *JobMgr) runJob(c context.Context, job *model.JobConfig, cancel <-chan bool) (err error) {
	if job, ok := m.jobs[job.ID.Hex()]; ok {
		template := &buildtemplate{}
		build := template.ConfigToBuild(job)
		go func() {
			pipe := NewPipeline(m.engine, &build.Works[0])
			pipe.Exec()

			timeout := time.After(time.Duration(60) * time.Minute)

			for {
				select {
				case <-pipe.Done():
					logrus.Debug("done")
					return
				case <-cancel:
					pipe.Stop()
					logrus.Debug("cancel")
					return
				case <-timeout:
					pipe.Stop()
					logrus.Debug("timeout")
					return
				case <-pipe.Step():
					logrus.Debugf("finish : \n %+v \n ", pipe.Curwork())
					if pipe.Curwork().Status == "fail" {
						pipe.Skip()
					} else {
						pipe.Exec()
					}
				case meassge := <-pipe.Msg():
					logrus.Debugf("meassge:%s error: %s", string(meassge.data), meassge.err.Error())
				}
			}
		}()
	} else {
		req.Entry(c).Infof("RunJob job  %s warning : notexsit \n", job.ID.Hex())
		return errors.New("Job not exsit")
	}
	return
}
