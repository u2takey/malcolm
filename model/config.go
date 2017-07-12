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
	ID              bson.ObjectId   `bson:"_id,omitempty"`
	Name            string          `bson:"name,omitempty"`
	Description     string          `bson:"description,omitempty"`
	ConcurrentBuild bool            `bson:"concurrentBuild"`
	Trigger         []TriggerConfig `bson:"trigger,omitempty"`
	Tasks           []Task          `bson:"pipeline,omitempty"`
	Service         Task            `bson:"service,omitempty"`
	Created         time.Time       `bson:"created"`
	Updated         time.Time       `bson:"udated"`
}

func (pipe *PipeConfig) Valid() error {
	msg := []string{}
	for _, task := range pipe.Tasks {
		err := task.Valid()
		if err != nil {
			msg = append(msg, fmt.Sprintf("error in config: %s / %s : [%s] ", task.Title, task.Type, task.Plugin))
		}
	}
	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "\n"))
	} else {
		return nil
	}
}

// Task represent pipeline step config
type Task struct {
	Title       string            `bson:"title,omitempty"`
	Type        string            `bson:"type,omitempty"`
	Plugin      string            `bson:"plugin,omitempty"`
	Environment map[string]string `bson:"environment,omitempty"` // use-> Options or Environment Options -> Environment
	Command     []string          `bson:"command,omitempty"`
	Args        []string          `bson:"args,omitempty"`
	PullPolicy  string            `bson:"pullPolicy,omitempty"`
	// key -> path
	Credentials map[string]string `bson:"credentials,omitempty"`
}

func (w *Task) Valid() error {
	if w.Plugin == "" {
		return errors.New("pluginPath cannot be empty")
	}
	return nil
}
