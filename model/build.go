package model

import (
	"time"

	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"labix.org/v2/mgo/bson"
)

// ManualTrigger ...: Trigger Represent trigger method to start a build
type ManualTrigger struct {
	TriggerParam map[string]string // for a trigger
}

type Trigger interface {
}

// Build is a running/runned instance of pipe
type Build struct {
	ID          bson.ObjectId `bson:"_id"`
	PipeID      bson.ObjectId `bson:"pipeid" index:"index"`
	Trigger     Trigger       `bson:"trigger,omitempty"`
	Title       string        `bson:"title,omitempty"`
	Description string        `bson:"description,omitempty"`
	Project     string        `bson:"project,omitempty" index:"index"`
	Status      BuildStatus   `bson:"status,omitempty" index:"index"`
	Started     time.Time     `bson:"started"`
	Finished    time.Time     `bson:"finished"`
	Updated     time.Time     `bson:"updated"`
	Works       []*Work       `bson:"works,omitempty"`
	Author      string        `bson:"author,omitempty"`
}

// Work is a instance of build, a single build may trigger multiple work
type Work struct {
	WorkNo      int         `bson:"workno"`
	Title       string      `bson:"title,omitempty"`
	Description string      `bson:"description,omitempty"`
	Status      WorkStatus  `bson:"status,omitempty"`
	Started     time.Time   `bson:"started"`
	Finished    time.Time   `bson:"finished"`
	Steps       []*WorkStep `bson:"steps,omitempty"`
}

// WorkStep is a step during a Work
type WorkStep struct {
	StepNo   int          `bson:"stepno"`
	Title    string       `bson:"title,omitempty"`
	Status   StepStatus   `bson:"status,omitempty"`
	Started  time.Time    `bson:"started"`
	Finished time.Time    `bson:"finished"`
	Config   *TaskGroup   `json:"-" bson:"-"`
	K8sjob   *batchv1.Job `json:"job,omitempty" bson:"job,omitempty"`
}
