package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"labix.org/v2/mgo/bson"

	mgoutil "github.com/arlert/malcolm/utils/mongo"
)

// Config is server config
type Config struct {
	MgoCfg   mgoutil.Config
	KubeAddr string
}

// ManualTriggerConfig config
type ManualTriggerConfig struct {
	Options []Option
}

type TriggerConfig struct {
	TriggerType string
}

// PipeConfig represent pipeline config
type PipeConfig struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Title       string          `bson:"title,omitempty"`
	Description string          `bson:"description,omitempty"`
	Trigger     []TriggerConfig `bson:"trigger,omitempty"`
	TaskGroups  []TaskGroup     `bson:"taskgroups,omitempty"`
	Services    []Task          `bson:"services,omitempty"`
	Matrix      MatrixEnv       `bson:"matrix,omitempty"`
	Created     time.Time       `bson:"created"`
	Updated     time.Time       `bson:"updated"`
}

// TaskGroup -> job -> pod -> onestep
type TaskGroup struct {
	Title         string        `bson:"title,omitempty"`
	Label         string        `bson:"title,omitempty"`
	Tasks         []Task        `bson:"tasks,omitempty"`
	Concurrent    []Task        `bson:"concurrent,omitempty"`
	Prerequisites Prerequisites `bson:"prerequisites,omitempty"`
}

// Prerequisites is tells when taskgroup should running
type Prerequisites struct {
	MatchExprs []string `bson:"matchexprs,omitempty"`
	// matchexpr : step.env match val
	RequireExpr string `bson:"requireexpr,omitempty"`
	// requireexpr : both/any/none
}

// MatrixEnv is for matrix work
type MatrixEnv map[string][]string

// Valid chech config param valid
func (pipe *PipeConfig) Valid() error {
	msg := []string{}
	for _, group := range pipe.TaskGroups {
		for _, task := range group.Tasks {
			err := task.Valid()
			if err != nil {
				msg = append(msg, fmt.Sprintf("error in config: %s / %s : [%s] ", group.Title, task.Title, task.Plugin))
			}
		}
	}
	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "\n"))
	}
	return nil
}

// Task -> single container
type Task struct {
	Title       string            `bson:"title,omitempty"`
	Plugin      string            `bson:"plugin,omitempty"`
	Environment map[string]string `bson:"environment,omitempty"` // use-> Options or Environment Options -> Environment
	Command     []string          `bson:"command,omitempty"`
	Args        []string          `bson:"args,omitempty"`
	PullPolicy  string            `bson:"pullPolicy,omitempty"`
	Ports       []int             `bson:"port,omitempty"`    // for service
	Timeout     int               `bson:"timeout,omitempty"` // timeout in minutes with default value
	// key -> path
	Credentials map[string]string `bson:"credentials,omitempty"`
}

// Valid check task param valid
func (w *Task) Valid() error {
	if w.Plugin == "" {
		return errors.New("pluginPath cannot be empty")
	}
	return nil
}
