package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"labix.org/v2/mgo/bson"

	mgoutil "github.com/arlert/malcolm/utils/mongo"
)

// ---------------------------------------------------------------------------
// server config
type Config struct {
	MgoCfg   mgoutil.Config
	KubeAddr string
}

// ---------------------------------------------------------------------------
// user config
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

// taskgroup -> job -> pod -> onestep
type TaskGroup struct {
	Title      string `bson:"title,omitempty"`
	Tasks      []Task `bson:"tasks,omitempty"`
	Concurrent []Task `bson:"concurrent,omitempty"`
}

type MatrixEnv map[string][]string

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
	} else {
		return nil
	}
}

// task -> single container
type Task struct {
	Title       string            `bson:"title,omitempty"`
	Plugin      string            `bson:"plugin,omitempty"`
	Environment map[string]string `bson:"environment,omitempty"` // use-> Options or Environment Options -> Environment
	Command     []string          `bson:"command,omitempty"`
	Args        []string          `bson:"args,omitempty"`
	PullPolicy  string            `bson:"pullPolicy,omitempty"`
	Ports       []int             `bson:"port,omitempty"` // for service
	// key -> path
	Credentials map[string]string `bson:"credentials,omitempty"`
}

func (w *Task) Valid() error {
	if w.Plugin == "" {
		return errors.New("pluginPath cannot be empty")
	}
	return nil
}
