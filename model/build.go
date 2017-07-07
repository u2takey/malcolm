package model

// Trigger Represent trigger method to start a build
type ManualTrigger struct {
	TriggerParam map[string]string // for a trigger
}

type ManualTriggerConfig struct {
	Options []Option
}

type Build struct {
	ID       int64    `json:"id"            db:"id,pk"`
	Trigger  Trigger  `json:"trigger"       db:"trigger"`
	Status   string   `json:"status"        db:"status"`
	Enqueued int64    `json:"enqueued"      db:"enqueued"`
	Created  int64    `json:"created"       db:"created"`
	Started  int64    `json:"started"       db:"started"`
	Finished int64    `json:"finished"      db:"finished"`
	Message  string   `json:"message"       db:"message"`
	Author   string   `json:"author"        db:"author"`
	Errors   []string `json:"error"         db:"error"`
}

type Log struct {
	ID      int64  `db:"log_id,pk"`
	BuildID int64  `db:"log_id"`
	Title   string `db:"log_title"`
	Data    []byte `db:"log_data"`
	Created int64  `db:"created"`
}

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

type Plugin struct {
	Name        string
	Version     string
	Indentifier string
	Options     []Option
}

type Option struct {
	Must        bool
	Key         string
	Default     string
	Choose      []string
	MultiChoose bool
}
