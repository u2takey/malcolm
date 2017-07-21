package model

type BuildStatus string

const (
	BuildStatusRunning  BuildStatus = "running"
	BuildStatusComplete BuildStatus = "complete"
	BuildStatusPaused   BuildStatus = "paused"
)

type WorkStatus struct {
	State       WorkState
	StateDetail StateDetail
	StateReason string
	Message     string
}

type StepStatus struct {
	State       StepState
	StateDetail StateDetail
	StateReason string
	Message     string
}

// WorkState is state of work
type WorkState string

const (
	WorkStateRunning  WorkState = "running"
	WorkStatePaused   WorkState = "paused"
	WorkStateComplete WorkState = "complete"
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
	StateCompleteDetailFailed  StateDetail = "failed"
	StateCompleteDetailStopped StateDetail = "stopped"
	StateCompleteDetailSkipped StateDetail = "skipped"
	StateCompleteDetailSuccess StateDetail = "success"
)

const (
	DefaultNameSpace = "malcolm" // #todo ->param
	PodTypeJob       = "jobpod"
	PodTypeService   = "servicepod"
)
