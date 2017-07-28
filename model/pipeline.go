package model

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Pipeline represent pipeline config
type Pipeline struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Title       string          `bson:"title,omitempty"`
	Description string          `bson:"description,omitempty"`
	Trigger     []TriggerConfig `bson:"trigger,omitempty"`
	TaskGroups  []TaskGroup     `bson:"taskgroups,omitempty"`
	Services    []Task          `bson:"services,omitempty"`
	Matrix      MatrixEnv       `bson:"matrix,omitempty"`
	Created     time.Time       `bson:"created"`
	Updated     time.Time       `bson:"updated"`
	Timeout     int             `bson:"timeout,omitempty"` // timeout in minutes with default value
}

// TaskGroup -> job -> pod -> onestep
type TaskGroup struct {
	Title         string        `bson:"title,omitempty"`
	Label         string        `bson:"label,omitempty"`
	PreTasks      []Task        `bson:"pretasks,omitempty"`
	Tasks         []Task        `bson:"tasks,omitempty"`
	Prerequisites Prerequisites `bson:"prerequisites,omitempty"`
	Timeout       int           `bson:"timeout,omitempty"` // timeout in minutes with default value
}

// Prerequisites is tells when taskgroup should running
type Prerequisites struct {
	MatchExprs []string `bson:"matchexprs,omitempty"`
	// matchexpr : step.env match val
	RequireExpr string `bson:"requireexpr,omitempty"`
	// requireexpr : both/any/none
}

// Task -> single container
type Task struct {
	Title       string            `bson:"title,omitempty"`
	Plugin      string            `bson:"plugin,omitempty"`
	Environment map[string]string `bson:"environment,omitempty"` // use-> Options or Environment Options -> Environment
	Command     []string          `bson:"command,omitempty"`
	Args        []string          `bson:"args,omitempty"`
	PullPolicy  string            `bson:"pullPolicy,omitempty"`
	Ports       []int             `bson:"port,omitempty"` // for service
	// Timeout     int               `bson:"timeout,omitempty"` // timeout in minutes with default value
	// key -> path
	Credentials map[string]string `bson:"credentials,omitempty"`
}

// MatrixEnv is for matrix work
type MatrixEnv map[string][]string

// ValidAndSetDefault check config param valid and set default
func (pipe *Pipeline) ValidAndSetDefault(config *GeneralBuildingConfig) (err error) {
	if pipe.Timeout == 0 {
		pipe.Timeout = config.WorkTimeoutDefault
	}
	for _, group := range pipe.TaskGroups {
		err1 := group.ValidAndSetDefault(config)
		err = WarpErrors(err, err1)
	}
	return
}

// ValidAndSetDefault check config param valid and set default
func (g *TaskGroup) ValidAndSetDefault(config *GeneralBuildingConfig) (err error) {
	if g.Timeout == 0 {
		g.Timeout = config.StepTimeoutDefault
	}
	for _, task := range g.Tasks {
		err1 := task.ValidAndSetDefault(config)
		err = WarpErrors(err, err1)
	}
	return
}

// ValidAndSetDefault check task param valid and set default
func (w *Task) ValidAndSetDefault(config *GeneralBuildingConfig) (err error) {
	if w.Plugin == "" {
		return errors.New("pluginPath cannot be empty")
	}
	if w.Title == "" {
		w.Title = w.Plugin
	}
	return nil
}
