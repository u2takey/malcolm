package model

const (
	EventPush   = "push"
	EventPull   = "pull_request"
	EventTag    = "tag"
	EventDeploy = "deployment"
)

const (
	StatusSkipped = "skipped"
	StatusPending = "pending"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusKilled  = "killed"
	StatusPaused  = "paused"
	StatusError   = "error"
)

type Option struct {
	Must        bool     `json:"must"`
	Key         string   `json:"key"`
	Default     string   `json:"default"`
	Choose      []string `json:"choose"`
	MultiChoose bool     `json:"multiChoose"`
}
