package model

import (
	mgoutil "github.com/arlert/malcolm/utils/mongo"
)

// Config is server config
type Config struct {
	MgoCfg   mgoutil.Config
	KubeAddr string
}

type Resource struct {
	Cpu    string
	Memory string
}

var DefaultBuildConfig = &GeneralBuildingConfig{
	Project:            DefaultNameSpace,
	WorkTimeoutDefault: 60,
	StepTimeoutDefault: 60,
	WorkSpace:          "/workspace",
	StorageSize:        "1Gi",
}

// GeneralBuildingConfig is project based building config
type GeneralBuildingConfig struct {
	Project              string
	WorkTimeoutDefault   int // in minute
	StepTimeoutDefault   int // in minute
	ResourceLimitDefault Resource
	ResourceLimitRequest Resource
	WorkSpace            string
	StorageClass         string
	StorageSize          string
}

// ManualTriggerConfig config
type ManualTriggerConfig struct {
	Options []Option
}

type TriggerConfig struct {
	TriggerType string
}
