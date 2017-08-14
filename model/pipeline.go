package model

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Pipeline represent pipeline config
type Pipeline struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Title        string        `bson:"title,omitempty"`
	Description  string        `bson:"description,omitempty"`
	WorkSpace    string        `bson:"workspace,omitempty"`
	StorageClass string        `bson:"storageClass,omitempty"`
	StorageSize  string        `bson:"storageSize,omitempty"`
	Trigger      TriggerConfig `bson:"trigger,omitempty"` // manual trigger default
	TaskGroups   []TaskGroup   `bson:"taskgroups,omitempty"`
	Services     []Task        `bson:"services,omitempty"`
	Matrix       MatrixEnv     `bson:"matrix,omitempty"`
	Created      time.Time     `bson:"created"`
	Updated      time.Time     `bson:"updated"`
	Timeout      int           `bson:"timeout,omitempty"` // timeout in minutes with default value
}

// TaskGroup -> job -> pod -> onestep
type TaskGroup struct {
	Title      string      `bson:"title,omitempty"`
	Label      string      `bson:"label,omitempty"`
	PreTasks   []Task      `bson:"pretasks,omitempty"`
	Tasks      []Task      `bson:"tasks,omitempty"`
	Constraint *Constraint `bson:"Constraint,omitempty"`
	Timeout    int         `bson:"timeout,omitempty"` // timeout in minutes with default value
}

// Prerequisites is tells when taskgroup should running
type Constraint struct {
	MatchState                         // last step.statedetail match
	MatchEnvs        map[string]string // step.env match val
	MatchExpressions []LabelSelectorRequirement
}

// MatchState :last step status
type MatchState string

const (
	MatchStateFail    MatchState = "Fail"
	MatchStateSuccess MatchState = "Success"
	MatchStateAlways  MatchState = "Always"
)

type LabelSelectorRequirement struct {
	Key      string
	Operator LabelSelectorOperator
	Values   []string
}

type LabelSelectorOperator string

const (
	LabelSelectorOpIn           LabelSelectorOperator = "In"
	LabelSelectorOpNotIn        LabelSelectorOperator = "NotIn"
	LabelSelectorOpExists       LabelSelectorOperator = "Exists"
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
)

// Task -> single container
type Task struct {
	Title       string                 `bson:"title,omitempty"`
	Plugin      string                 `bson:"plugin,omitempty"`
	Environment map[string]interface{} `bson:"environment,omitempty"` // plugin config -> environment
	Command     []string               `bson:"command,omitempty"`
	Privileged  bool                   `bson:"privileged,omitempty"`
	Args        []string               `bson:"args,omitempty"`
	PullPolicy  string                 `bson:"pullPolicy,omitempty"`
	Ports       []int                  `bson:"port,omitempty"` // for service
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
	if pipe.WorkSpace == "" {
		pipe.WorkSpace = config.WorkSpace
	}
	if pipe.StorageClass != "" && pipe.StorageSize == "" {
		pipe.StorageSize = config.StorageSize
	}
	return
}

// ValidAndSetDefault check config param valid and set default
func (g *TaskGroup) ValidAndSetDefault(config *GeneralBuildingConfig) (err error) {
	if g.Timeout == 0 {
		g.Timeout = config.StepTimeoutDefault
	}
	if g.Constraint == nil {
		g.Constraint = &Constraint{
			MatchState: config.ConstriantStateDefault,
		}
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
