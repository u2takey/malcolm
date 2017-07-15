package pipemgr

import (
	"errors"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"labix.org/v2/mgo/bson"

	"github.com/arlert/malcolm/model"
	"github.com/arlert/malcolm/store"
	// "github.com/arlert/malcolm/utils"
	req "github.com/arlert/malcolm/utils/reqlog"
)

type PipeMgr struct {
	pipelock sync.RWMutex
	pipes    map[string]*model.PipeConfig
	store    *store.Store
	engine   *Engine
}

func NewPipeMgr(st *store.Store, eng *Engine) (j *PipeMgr) {
	j = &PipeMgr{
		store:  st,
		engine: eng,
	}
	err := j.ListPipe()
	if err != nil {
		logrus.Fatalln("PipeMgr.list error -", err)
	}
	return j
}

func (m *PipeMgr) ListPipe() (err error) {
	pipel := []model.PipeConfig{}
	err = m.store.Cols.Pipe.Find(bson.M{}).All(&pipel)
	if err != nil {
		return
	}
	m.pipelock.Lock()
	defer m.pipelock.Unlock()
	pipem := make(map[string]*model.PipeConfig, 0)
	for _, pipe := range pipel {
		pipem[pipe.ID.Hex()] = &pipe
	}
	m.pipes = pipem
	logrus.Infof("list pipes, count : %d \n", len(pipem))
	return
}

func (m *PipeMgr) AddPipe(c context.Context, pipe *model.PipeConfig) {
	m.pipelock.Lock()
	defer m.pipelock.Unlock()
	m.pipes[pipe.ID.Hex()] = pipe
}

func (m *PipeMgr) RemovePipe(c context.Context, pipeid string) {
	m.pipelock.Lock()
	defer m.pipelock.Unlock()
	if _, ok := m.pipes[pipeid]; ok {
		delete(m.pipes, pipeid)
	} else {
		req.Entry(c).Infof("remove pipe  %s warning : notexsit \n", pipeid)
	}
}

func (m *PipeMgr) RunPipe(c context.Context, pipeid string) (build *model.Build, err error) {
	m.pipelock.RLock()
	defer m.pipelock.RUnlock()
	if pipe, ok := m.pipes[pipeid]; ok {
		build, err = m.runPipe(c, pipe, make(chan bool))
	} else {
		req.Entry(c).Infof("RunPipe pipe  %s warning : notexsit \n", pipeid)
		err = errors.New("Pipe not exsit")
	}
	return
}

func (m *PipeMgr) runPipe(c context.Context, pipe *model.PipeConfig, cancel <-chan bool) (build *model.Build, err error) {
	if pipe, ok := m.pipes[pipe.ID.Hex()]; ok {
		template := &buildtemplate{}
		build, err = template.ConfigToBuild(pipe)
		if err != nil {
			req.Entry(c).Warning("error in ConfigToBuild", err)
			return nil, err
		}
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
					logrus.Debugf("finish : %s : %s  ", pipe.Curwork().Config.Title, pipe.Curwork().Config.Title)
					if pipe.Curwork().Status == "fail" {
						pipe.Skip()
					} else {
						pipe.Exec()
					}
				case meassge := <-pipe.Msg():
					logrus.Debugf("meassge:%s error: %s", meassge.data, meassge.err)
				}
			}
		}()
	} else {
		req.Entry(c).Infof("RunPipe pipe  %s warning : notexsit \n", pipe.ID.Hex())
		return nil, errors.New("Pipe not exsit")
	}
	return
}
