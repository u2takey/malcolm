package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	v1 "k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
)

// ManualTrigger ...: Trigger Represent trigger method to start a build
type ManualTrigger struct {
	TriggerParam map[string]string // for a trigger
}

type Trigger interface {
}

// Build is a running/runned instance of pipeline
type Build struct {
	ID          bson.ObjectId             `bson:"_id"`
	PipeID      bson.ObjectId             `bson:"pipeid" index:"index"`
	Trigger     Trigger                   `bson:"trigger,omitempty"`
	Title       string                    `bson:"title,omitempty"`
	Description string                    `bson:"description,omitempty"`
	Project     string                    `bson:"project,omitempty" index:"index"`
	Status      BuildStatus               `bson:"status,omitempty" index:"index"`
	Created     time.Time                 `bson:"created"`
	Started     time.Time                 `bson:"started"`
	Finished    time.Time                 `bson:"finished"`
	Updated     time.Time                 `bson:"updated"`
	CurrentStep int                       `bson:"currentStep,omitempty"`
	Steps       []*WorkStep               `bson:"steps,omitempty"`
	Author      string                    `bson:"author,omitempty"`
	Dirty       bool                      `bson:"-" json:"-"`
	Volumn      *v1.PersistentVolumeClaim `bson:"volumn,omitempty" json:"volumn,omitempty"`
}

// WorkStep is a step during a Build
type WorkStep struct {
	Title    string       `bson:"title,omitempty"`
	StepNo   int          `bson:"stepno,omitempty"`
	Status   StepStatus   `bson:"status,omitempty"`
	Started  time.Time    `bson:"started"`
	Finished time.Time    `bson:"finished"`
	Config   *TaskGroup   `bson:"-" json:"-"`
	K8sjob   *batchv1.Job `bson:"job,omitempty" json:"job,omitempty"`
}
