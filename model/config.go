package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"labix.org/v2/mgo/bson"

	mgoutil "github.com/u2takey/malcolm/utils/mongo"
)

type Config struct {
	MgoCfg   mgoutil.Config
	KubeAddr string
}

type ManualTriggerConfig struct {
	Options []Option
}

type TriggerConfig struct {
	TriggerType string
}

// JobConfig represent pipeline config
// JobConfig contains metadata and workconfig
type JobConfig struct {
	ID              bson.ObjectId   `bson:"_id,omitempty"`
	Name            string          `bson:"name,omitempty"`
	Description     string          `bson:"description,omitempty"`
	ConcurrentBuild bool            `bson:"concurrentBuild"`
	Trigger         []TriggerConfig `bson:"trigger,omitempty"`
	Scm             []WorkConfig    `bson:"scm,omitempty"`
	PreBuild        []WorkConfig    `bson:"preBuild,omitempty"`
	Build           []WorkConfig    `bson:"build,omitempty"`
	AfterBuild      []WorkConfig    `bson:"afterBuild,omitempty"`
	Notify          []WorkConfig    `bson:"notify,omitempty"`
	Service         []WorkConfig    `bson:"service,omitempty"`
	Created         time.Time       `bson:"created"`
	Updated         time.Time       `bson:"udated"`
	// Publish         []WorkConfig    `json:"publish"`
}

func (job *JobConfig) Valid() error {
	msg := []string{}
	workconfigs := map[string][]WorkConfig{
		"scm":        job.Scm,
		"prebuild":   job.PreBuild,
		"build":      job.Build,
		"afterbuild": job.AfterBuild,
		"notify":     job.Notify,
		"service":    job.Service,
	}

	for key, val := range workconfigs {
		for index, work := range val {
			work.Type = key // add type
			err := work.Valid()
			if err != nil {
				msg = append(msg, fmt.Sprintf("error in config: %s %s/%d ", work.Title, key, index))
			}
		}
	}
	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "\n"))
	} else {
		return nil
	}
}

// WorkConfig represent pipeline step config
// WorkConfig contains running data which will be convert into Work
type WorkConfig struct {
	Title       string            `bson:"title,omitempty"`
	Type        string            `bson:"type,omitempty"`
	Plugin      string            `bson:"plugin,omitempty"`
	Environment map[string]string `bson:"environment,omitempty"` // use-> Options or Environment Options -> Environment
	Command     []string          `bson:"command,omitempty"`
	Args        []string          `bson:"args,omitempty"`
	PullPolicy  string            `bson:"pullPolicy,omitempty"`
	// use secret or configmap
	CredentialsName string `bson:"credentialsName,omitempty"`
	CredentialsPath string `bson:"credentialsPath,omitempty"`

	// -- auto set if needed --
	// Resources ResourceRequirements `json:"resources,omitempty" `
	// VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
	// Lifecycle       *Lifecycle `json:"lifecycle,omitempty"`
	// WorkingDir      string     `json:"workingDir,omitempty"`

	// -- not used ---
	// SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Stdin bool `json:"stdin,omitempty" `
	// StdinOnce bool `json:"stdinOnce,omitempty" `
	// TTY bool `json:"tty,omitempty" `
	// TerminationMessagePath string `json:"terminationMessagePath,omitempty" `
	// TerminationMessagePolicy TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty"`
	// ReadinessProbe *Probe `json:"readinessProbe,omitempty" `
	// LivenessProbe *Probe `json:"livenessProbe,omitempty" `
	// EnvFrom        []EnvFromSource     `json:"envFrom,omitempty"`
	// Ports []ContainerPort `json:"ports,omitempty"`
}

func (w *WorkConfig) Valid() error {
	if w.Plugin == "" {
		return errors.New("pluginPath cannot be empty")
	}
	return nil
}
