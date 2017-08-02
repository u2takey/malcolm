package model

// BuildAction is trigger input
type BuildAction string

// BuildAction types
const (
	ActionStart  BuildAction = "start"
	ActionPause  BuildAction = "pause"
	ActionResume BuildAction = "resume"
	ActionStop   BuildAction = "stop"
)

// BuildStatus represent building status
type BuildStatus struct {
	State       BuildState
	StateDetail StateDetail
	Message     string
}

// StepStatus represent building step status
type StepStatus struct {
	State       StepState
	StateDetail StateDetail
	Message     string
}

// BuildState is general state info of build
type BuildState string

// BuildState types
const (
	BuildStatePending  BuildState = "pending"
	BuildStateRunning  BuildState = "running"
	BuildStatePaused   BuildState = "paused"
	BuildStatePausing  BuildState = "pausing"
	BuildStateComplete BuildState = "complete"
)

// StepState is state of a build step
type StepState string

// StepState type
const (
	StepStatePending  StepState = "pending"
	StepStateRunning  StepState = "running"
	StepStateComplete StepState = "complete"
)

// StateDetail is status detail of build or step state
type StateDetail string

// StateDetail types
const (
	StateCompleteDetailFailed   StateDetail = "failed"
	StateCompleteDetailCanceled StateDetail = "canceled"
	StateCompleteDetailSkipped  StateDetail = "skipped"
	StateCompleteDetailSuccess  StateDetail = "success"
)

// consts
const (
	DefaultNameSpace = "malcolm" // #todo ->param
	PodTypeJob       = "jobpod"
	PodTypeService   = "servicepod"
)
