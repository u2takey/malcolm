package pipemgr

import (
	"errors"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"

	"github.com/arlert/malcolm/model"
	"github.com/arlert/malcolm/server/cronmgr"
	"github.com/arlert/malcolm/store"
	// "github.com/arlert/malcolm/utils"
	"fmt"

	req "github.com/arlert/malcolm/utils/reqlog"
)

// PipeManagerInterface is interface for manage pipelines
type PipeManagerInterface interface {
	PipelineAdd(c context.Context, pipe *model.Pipeline) error
	PipelineRemove(c context.Context, pipeid string) error
	PipelineFind(c context.Context, pipeid string) (pipe *model.Pipeline, err error)

	BuildAction(c context.Context, buildid, pipeid string, action model.BuildAction) (build *model.Build, err error)
	BuildQueue(c context.Context) (builds []*model.Build, err error)

	Run()
}

// ControlEvent receive control event from external
type ControlEvent struct {
	action   model.BuildAction
	build    *model.Build
	pipeline *model.Pipeline
}

// BuildEventType is build event from engine
type BuildEventType int

const (
	eventNone BuildEventType = iota
	eventError
	eventStatusStart
	eventStatusDone
	eventStatusDeleted
)

// BuildEvent receive build event from internal
type BuildEvent struct {
	eventtype BuildEventType
	data      string
	buildid   string
}

// PipeManager implement PipeManagerInterface
type PipeManager struct {
	store     *store.Store
	pipes     map[string]*model.Pipeline
	builds    map[string]*model.Build
	pipelock  sync.RWMutex
	buildlock sync.RWMutex

	buildEvent   chan BuildEvent
	controlEvent chan ControlEvent

	engine  BuildEngineInterface
	cronmgr cronmgr.Interface
}

// NewPipeManager init a Pipemanager
func NewPipeManager(st *store.Store, eng BuildEngineInterface) (m *PipeManager) {
	m = &PipeManager{
		store:        st,
		engine:       eng,
		pipes:        make(map[string]*model.Pipeline, 0),
		builds:       make(map[string]*model.Build, 0),
		buildEvent:   make(chan BuildEvent, 100),
		controlEvent: make(chan ControlEvent, 100),
		cronmgr:      cronmgr.New(),
	}
	return m
}

//Run is main loop of PipeManager
func (m *PipeManager) Run() {
	err := m.listPipe()
	if err != nil {
		logrus.Fatalln("PipeManager.listPipe error -", err)
	}
	err = m.listBuild()
	if err != nil {
		logrus.Fatalln("PipeManager.listBuild error -", err)
	}
	t := time.NewTicker(time.Second)
	c := context.Background()

	m.cronmgr.Start()
	for _, pipe := range m.pipes {
		err1 := m.cronmgr.UpInsert(pipe)
		if err1 != nil {
			logrus.Warningln("cronmgr.UpInsert error -", err1)
		}
	}

	go m.engine.Start(m.buildEvent)

	go func() {
		for {
			var curstep *model.WorkStep
			select {
			case ev := <-m.buildEvent:
				// build step event from engine
				// update step status -> build stauts; trigger step run/skip
				build := m.buildFind(c, ev.buildid)
				if build == nil {
					logrus.Error("unexpected build event ", ev)
					break
				}
				build.Updated = time.Now()
				build.Dirty = true

				if build.CurrentStep < len(build.Steps) {
					curstep = build.Steps[build.CurrentStep]
				}
				if curstep == nil {
					logrus.Error("unexpected error:  curstep not found  ")
					break
				}
				switch ev.eventtype {
				case eventError:
					// update step status
					curstep.Status.State = model.StepStateComplete
					curstep.Status.StateDetail = model.StateCompleteDetailFailed
					curstep.Finished = time.Now()
					curstep.Status.Message = ev.data
					// get error goto nextstep
					m.nextstep(build)
				case eventStatusStart:
					// update step status
					curstep.Started = time.Now()
					curstep.Status.State = model.StepStateRunning

					// update build status
					if build.Status.State == model.BuildStatePending {
						build.Started = time.Now()
						build.Status.State = model.BuildStateRunning
					}
				case eventStatusDone:
					// update step status
					curstep.Status.State = model.StepStateComplete
					curstep.Status.StateDetail = model.StateCompleteDetailSuccess
					curstep.Finished = time.Now()

					// update build status -> goto next step
					m.nextstep(build)
					if build.Status.State == model.BuildStatePausing {
						build.Status.State = model.BuildStatePaused
					}
				case eventStatusDeleted:
					// todo
				}
			case ev := <-m.controlEvent:
				// control event from external, update build state
				build := ev.build
				build.Updated = time.Now()
				build.Dirty = true
				if build.CurrentStep < len(build.Steps) {
					curstep = build.Steps[build.CurrentStep]
				}
				switch ev.action {
				case model.ActionStart:
					if build.Status.State != model.BuildStatePending {
						logrus.Error("unexpected error:  start a job not in pending  ", build.Status.State)
						break
					}
					// start build
					err := m.engine.RunStep(build.Steps[build.CurrentStep])
					if err != nil {
						build.Status.State = model.BuildStateComplete
						build.Finished = time.Now()
						build.Status.StateDetail = model.StateCompleteDetailFailed
						build.Status.Message = err.Error()
					}
				case model.ActionResume:
					if build.Status.State != model.BuildStatePaused {
						logrus.Error("unexpected error:  resume a job not in paused  ", build.Status.State)
						break
					}
					if curstep == nil {
						logrus.Error("unexpected error:  curstep not found  ", build.CurrentStep)
						break
					}
					// start build
					err := m.engine.RunStep(curstep)
					if err != nil {
						build.Status.State = model.BuildStateComplete
						build.Finished = time.Now()
						build.Status.StateDetail = model.StateCompleteDetailFailed
						build.Status.Message = err.Error()
					} else {
						build.Status.State = model.BuildStatePending
					}
				case model.ActionPause:
					if build.Status.State != model.BuildStateRunning {
						logrus.Error("unexpected error:  pause a job not in running  ", build.Status.State)
						break
					}
					if build.CurrentStep >= len(build.Steps) {
						logrus.Error("unexpected error:  curstep not found  ", build.CurrentStep)
						break
					}
					if curstep == nil {
						logrus.Error("unexpected error:  curstep not found  ", build.CurrentStep)
						break
					}
					if curstep.Status.State == model.StepStateRunning {
						build.Status.State = model.BuildStatePausing
					} else {
						build.Status.State = model.BuildStatePaused
					}
				case model.ActionStop:
					build.Status.State = model.BuildStateComplete
					build.Status.StateDetail = model.StateCompleteDetailCanceled
					build.Finished = time.Now()
				}
			case pipeid := <-m.cronmgr.CronPipeChan():
				go func() {
					//#todo Trigger set to refractor
					build, err := m.BuildAction(context.Background(), "", pipeid, model.ActionStart)
					if err != nil {
						logrus.Error("BuildAction Start error ", err)
					} else {
						build.Trigger = &model.CronTrigger{}
					}
				}()
			case <-t.C:
				m.syncStatus()
			}
		}
	}()
}

func (m *PipeManager) nextstep(build *model.Build) {
next:
	build.CurrentStep++
	if build.CurrentStep >= len(build.Steps) {
		// finished
		build.Status.State = model.BuildStateComplete
		build.Finished = time.Now()
		fail := false
		for _, step := range build.Steps {
			if step.Status.StateDetail == model.StateCompleteDetailFailed {
				build.Status.StateDetail = model.StateCompleteDetailFailed
				fail = true
				break
			}
		}
		if !fail {
			build.Status.StateDetail = model.StateCompleteDetailSuccess
		}
	} else if build.Status.State == model.BuildStateRunning {
		nextstep := build.Steps[build.CurrentStep]
		if !matchConstraint(build) {
			nextstep.Status.State = model.StepStateComplete
			nextstep.Finished = time.Now()
			nextstep.Status.StateDetail = model.StateCompleteDetailSkipped
			nextstep.Status.Message = "constraint not match : skipped"
			goto next
		}

		// match and not finish
		err := m.engine.RunStep(nextstep)
		if err != nil {
			// run step fail -> complete whole build
			nextstep.Status.State = model.StepStateComplete
			nextstep.Finished = time.Now()
			nextstep.Status.StateDetail = model.StateCompleteDetailFailed
			nextstep.Status.Message = err.Error()
			goto next
		}
	}
}

// listPipe populate pipes from store to local cache
func (m *PipeManager) listPipe() (err error) {
	pipel := []model.Pipeline{}
	err = m.store.Cols.Pipe.Find(bson.M{}).All(&pipel)
	if err != nil {
		return
	}
	m.pipelock.Lock()
	defer m.pipelock.Unlock()
	for index, pipe := range pipel {
		pipeid := pipe.ID.Hex()
		m.pipes[pipeid] = &pipel[index]
	}
	logrus.Infof("list pipes success, count : %d \n", len(m.pipes))
	return
}

func (m *PipeManager) listBuild() (err error) {
	buildl := []model.Build{}
	err = m.store.Cols.Build.Find(bson.M{
		"status.state": bson.M{"$ne": model.BuildStateComplete},
	}).All(&buildl)
	if err != nil {
		return
	}
	m.buildlock.Lock()
	defer m.buildlock.Unlock()
	for index := range buildl {
		build := &buildl[index]
		build.Status.State = model.BuildStateComplete
		build.Finished = time.Now()
		if build.Started.Unix() == 0 {
			build.Started = time.Now()
		}
		build.Status.StateDetail = model.StateCompleteDetailCanceled
		buildid := build.ID.Hex()
		m.builds[buildid] = build
	}
	logrus.Infof("list builds success, count : %d \n", len(m.builds))
	return
}

// PipelineAdd add pipeline to cache
func (m *PipeManager) PipelineAdd(c context.Context, pipe *model.Pipeline) error {
	m.pipelock.Lock()
	defer m.pipelock.Unlock()
	m.pipes[pipe.ID.Hex()] = pipe
	err := m.cronmgr.UpInsert(pipe)
	if err != nil {
		return err
	}
	return nil
}

// PipelineRemove remove cached pipeline
func (m *PipeManager) PipelineRemove(c context.Context, pipeid string) (err error) {
	m.pipelock.Lock()
	defer m.pipelock.Unlock()
	if _, ok := m.pipes[pipeid]; ok {
		delete(m.pipes, pipeid)
		m.cronmgr.Delete(pipeid)
	} else {
		err = errors.New("pipeline not found in cache")
	}
	return
}

// PipelineFind return matched pipe in cache
func (m *PipeManager) PipelineFind(c context.Context, pipeid string) (pipe *model.Pipeline, err error) {
	m.pipelock.Lock()
	defer m.pipelock.Unlock()
	if pipe, ok := m.pipes[pipeid]; ok {
		return pipe, nil
	}
	err = errors.New("pipeline not found in cache")
	return
}

// buildAdd add build into cache
func (m *PipeManager) buildAdd(c context.Context, build *model.Build) {
	m.buildlock.Lock()
	defer m.buildlock.Unlock()
	m.builds[build.ID.Hex()] = build
}

// buildRemove delete build from cache
func (m *PipeManager) buildRemove(c context.Context, build *model.Build) {
	m.buildlock.Lock()
	defer m.buildlock.Unlock()
	delete(m.builds, build.ID.Hex())
}

// buildFind find build in cache
func (m *PipeManager) buildFind(c context.Context, id string) *model.Build {
	m.buildlock.Lock()
	defer m.buildlock.Unlock()
	if build, ok := m.builds[id]; ok {
		return build
	}
	return nil
}

// BuildQueue return builds in queue
func (m *PipeManager) BuildQueue(c context.Context) (res []*model.Build, err error) {
	m.buildlock.Lock()
	defer m.buildlock.Unlock()
	for _, val := range m.builds {
		res = append(res, val)
	}
	return
}

//getBuild convert pipeline config into build
func (m *PipeManager) getBuild(c context.Context, pipe *model.Pipeline) (build *model.Build, err error) {
	template := &buildtemplate{}
	build, err = template.ConfigToBuild(pipe)
	if err != nil {
		req.Entry(c).Warning("error in ConfigToBuild", err)
	}
	build.Created = time.Now()
	build.Status.State = model.BuildStatePending
	for _, step := range build.Steps {
		step.Status.State = model.StepStatePending
	}
	return
}

func (m *PipeManager) buildSetup(c context.Context, build *model.Build) (err error) {
	if build.Volumn != nil {
		err = m.engine.SetupVolumn(build.Volumn)
	}
	return
}

func (m *PipeManager) buildCleanup(c context.Context, build *model.Build) (err error) {
	if build.Volumn != nil {
		err = m.engine.CleanupVolumn(build.Volumn)
	}
	return
}

// BuildAction receive build action
func (m *PipeManager) BuildAction(c context.Context,
	buildid, pipeid string, action model.BuildAction) (build *model.Build, err error) {

	var pipe *model.Pipeline
	pipe, err = m.PipelineFind(c, pipeid)
	if err != nil {
		return
	}
	logrus.Debugln("BuildAction", action, pipeid, pipe.ID)
	if buildid != "" {
		build = m.buildFind(c, buildid)
		if build == nil {
			return nil, errors.New("build not exsit")
		}
	}

	switch action {
	// start a new build
	case model.ActionStart:
		if build != nil {
			err = errors.New("you cannot start a exsit build")
			return
		}
		build, err = m.getBuild(c, pipe)
		if err != nil {
			return
		}
		err = m.buildSetup(c, build)
		if err != nil {
			return
		}
		m.buildAdd(c, build)
	case model.ActionPause:
		if build == nil || build.Status.State != model.BuildStateRunning {
			err = fmt.Errorf("build not exist in queue or not in running state:%s", build.Status.State)
			return
		}
	case model.ActionStop:
		if build == nil || build.Status.State == model.BuildStateComplete {
			err = fmt.Errorf("build not exist in queue in complete state:%s", build.Status.State)
			return
		}
	case model.ActionResume:
		if build == nil || build.Status.State != model.BuildStatePaused {
			err = fmt.Errorf("build not exist in queue or not in paused state:%s", build.Status.State)
			return
		}
	}

	m.controlEvent <- ControlEvent{
		action:   action,
		build:    build,
		pipeline: pipe,
	}
	return
}

// syncStatus sync cached build status into store
func (m *PipeManager) syncStatus() {
	m.buildlock.Lock()
	defer m.buildlock.Unlock()
	c := context.Background()
	temp := make([]*model.Build, 0)
	for _, build := range m.builds {
		if build.Dirty == true {
			sel := bson.M{"_id": build.ID}
			_, err := m.store.Cols.Build.Upsert(sel, build)
			if err != nil {
				logrus.Debug(err)
			}
			// b, err := json.MarshalIndent(build, "", "	")
			// logrus.Debug(string(b))
			build.Dirty = false
			logrus.Debugf("syncStatus:%s, total:%d", build.ID, len(m.builds))
		}
		if build.Status.State == model.BuildStateComplete {
			temp = append(temp, build)
		}
	}

	for _, build := range temp {
		err := m.buildCleanup(c, build)
		if err != nil {
			logrus.Errorln("buildCleanup error ", err)
		}
		delete(m.builds, build.ID.Hex())
	}
}
