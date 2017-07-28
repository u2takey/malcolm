package model

type BuildAction string

const (
	ActionStart  BuildAction = "start"
	ActionPause  BuildAction = "pause"
	ActionResume BuildAction = "resume"
	ActionStop   BuildAction = "stop"
)

type BuildStatus struct {
	State       BuildState
	StateDetail StateDetail
	Message     string
}

type StepStatus struct {
	State       StepState
	StateDetail StateDetail
	Message     string
}

// WorkState is state of work
type BuildState string

const (
	BuildStatePending  BuildState = "pending"
	BuildStateRunning  BuildState = "running"
	BuildStatePaused   BuildState = "paused"
	BuildStatePausing  BuildState = "pausing"
	BuildStateComplete BuildState = "complete"
)

// StepState is state of a build step
type StepState string

const (
	StepStatePending  StepState = "pending"
	StepStateRunning  StepState = "running"
	StepStateComplete StepState = "complete"
)

type StateDetail string

const (
	StateCompleteDetailFailed   StateDetail = "failed"
	StateCompleteDetailCanceled StateDetail = "canceled"
	StateCompleteDetailSkipped  StateDetail = "skipped"
	StateCompleteDetailSuccess  StateDetail = "success"
)

const (
	DefaultNameSpace = "malcolm" // #todo ->param
	PodTypeJob       = "jobpod"
	PodTypeService   = "servicepod"
)
