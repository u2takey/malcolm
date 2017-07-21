package pipemgr

import (
	"github.com/Sirupsen/logrus"

	"github.com/arlert/malcolm/model"
)

var (
	DefaultBufferSize = 100
)

type signal struct{}
type meassge struct {
	data string
	err  error
}

type Pipeline struct {
	index     int
	work      *model.Work
	stepmsg   chan (*meassge)
	stepstart chan (signal)
	stepdone  chan (signal)
	done      chan (signal)
	engine    *Engine
}

func NewPipeline(eng *Engine, work *model.Work) *Pipeline {

	pipeline := Pipeline{
		engine:    eng,
		work:      work,
		stepmsg:   make(chan *meassge, DefaultBufferSize),
		stepstart: make(chan signal),
		stepdone:  make(chan signal),
		done:      make(chan signal),
	}
	return &pipeline
}

func (p *Pipeline) Done() <-chan signal {
	return p.done
}

func (p *Pipeline) StepMsg() <-chan (*meassge) {
	return p.stepmsg
}

func (p *Pipeline) StepStart() <-chan signal {
	return p.stepstart
}

func (p *Pipeline) StepDone() <-chan signal {
	return p.stepdone
}

func (p *Pipeline) CurStep() *model.WorkStep {
	if p.index < len(p.work.Steps) && p.index >= 0 {
		return p.work.Steps[p.index]
	}
	return nil
}

func (p *Pipeline) NextStep() *model.WorkStep {
	if p.index < len(p.work.Steps)-1 {
		return p.work.Steps[p.index+1]
	}
	return nil
}

func (p *Pipeline) Finished() bool {
	return p.index >= len(p.work.Steps)-1
}

// Exec executes the current step.
func (p *Pipeline) Exec() {
	go func() {
		curstep := p.CurStep()
		p.stepstart <- signal{}
		if curstep != nil {
			p.exec(curstep)
		}
		p.stepdone <- signal{}
		p.next()
	}()
}

// Skip skips the current step.
func (p *Pipeline) Skip() {
	go func() {
		p.next()
	}()
}

func (p *Pipeline) Setup() error {
	return nil
}

// Do clean up
func (p *Pipeline) Teardown() {
}

func (p *Pipeline) Stop() {
	go func() {
		// kill current work ?
		p.done <- signal{}
	}()
}

func (p *Pipeline) next() {
	logrus.Debugf("next : %s ", p.CurStep().Config.Title)
	if p.Finished() {
		p.done <- signal{}
	} else {
		p.index++
	}
}

func (p *Pipeline) exec(work *model.WorkStep) {
	p.engine.RunSyncJobch(work, p.stepmsg)
}
