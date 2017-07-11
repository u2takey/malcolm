package jobmgr

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
	index  int
	work   *model.Work
	msg    chan (*meassge)
	step   chan (signal)
	done   chan (signal)
	engine *Engine
}

func NewPipeline(eng *Engine, work *model.Work) *Pipeline {

	pipeline := Pipeline{
		engine: eng,
		work:   work,
		msg:    make(chan *meassge, DefaultBufferSize),
		step:   make(chan signal),
		done:   make(chan signal),
	}
	return &pipeline
}

func (p *Pipeline) Done() <-chan signal {
	return p.done
}

func (p *Pipeline) Msg() <-chan (*meassge) {
	return p.msg
}

func (p *Pipeline) Step() <-chan signal {
	return p.step
}

func (p *Pipeline) Curwork() *model.WorkStep {
	if p.index < len(p.work.Steps) && p.index >= 0 {
		return p.work.Steps[p.index]
	}
	return nil
}

func (p *Pipeline) Finished() bool {
	return p.index >= len(p.work.Steps)-1
}

// Exec executes the current step.
func (p *Pipeline) Exec() {
	go func() {
		curwork := p.Curwork()
		if curwork != nil {
			p.exec(curwork)
		}
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
	logrus.Debug("next", p.Curwork())
	if p.Finished() {
		p.done <- signal{}
	} else {
		p.index += 1
		p.step <- signal{}
	}
}

func (p *Pipeline) exec(work *model.WorkStep) {
	p.engine.RunSyncJobch(work, p.msg)
}
