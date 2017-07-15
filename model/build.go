package model

import (
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"labix.org/v2/mgo/bson"
	"time"
)

// Trigger Represent trigger method to start a build
type ManualTrigger struct {
	TriggerParam map[string]string // for a trigger
}

type Trigger interface {
}

type Build struct {
	ID              bson.ObjectId `bson:"_id"`
	PipeID          bson.ObjectId `bson:"pipeid" index:"index"`
	Trigger         Trigger       `bson:"trigger,omitempty"`
	Title           string        `bson:"title,omitempty"`
	Description     string        `bson:"description,omitempty"`
	ConcurrentBuild bool          `bson:"concurrentBuild"`
	Project         string        `bson:"project,omitempty" index:"index"`
	Status          string        `bson:"status,omitempty" index:"index"`
	Created         time.Time     `bson:"created"`
	Started         time.Time     `bson:"started"`
	Finished        time.Time     `bson:"finished"`
	Message         []byte        `bson:"message,omitempty"`
	Works           []Work        `bson:"works,omitempty"`
	Author          string        `bson:"author,omitempty"`
	Error           string        `bson:"error,omitempty"`
}

type Work struct {
	WorkNo      int         `bson:"workno"`
	Title       string      `bson:"title,omitempty"`
	Description string      `bson:"description,omitempty"`
	Status      string      `bson:"status,omitempty"`
	Started     time.Time   `bson:"started"`
	Finished    time.Time   `bson:"finished"`
	Message     []byte      `bson:"message,omitempty"`
	Error       string      `bson:"error,omitempty"`
	Steps       []*WorkStep `bson:"steps,omitempty"`
}

type WorkStep struct {
	StepNo   int          `bson:"stepno"`
	Title    string       `bson:"title,omitempty"`
	Status   string       `bson:"status,omitempty"`
	Started  time.Time    `bson:"started"`
	Finished time.Time    `bson:"finished"`
	Message  []byte       `bson:"message,omitempty"`
	Error    string       `bson:"error,omitempty"`
	Config   *TaskGroup   `json:"-" bson:"-"`
	K8sjob   *batchv1.Job `json:"job,omitempty" bson:"job,omitempty"`
}

// log will not be saved into db
// type Log struct {
// 	ID          bson.ObjectId `bson:"_id"`
// 	BuildID     int64         `bson:"buildid" index:"index"`
// 	Title       string        `bson:"title,omitempty"`
// 	Description string        `bson:"description,omitempty"`
// 	Step        string        `bson:"step"`
// 	Data        []byte        `bson:"data,omitempty"`
// 	Created     time.Time     `bson:"created"`
// }

// Constraint represent step constraint between steps
// type Constraint struct {
//  Include []string
//  Exclude []string
// }

// type Constraints struct {
//  Repo        Constraint
//  Ref         Constraint
//  Refspec     Constraint
//  Environment Constraint
//  Event       Constraint
//  Branch      Constraint
//  Status      Constraint
// }
