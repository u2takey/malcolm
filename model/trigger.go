package model

//Trigger type
const (
	TriggerTypeManual   = "manual"
	TriggerTypeCron     = "cron"
	TriggerTypeWeebhook = "webhook"
)

//Trigger template type
var (
	TriggerManual = TriggerConfigTemplate{
		Type:   "manual",
		Schema: "", // schema for TriggerManual is set by user
	}
	TriggerCron = TriggerConfigTemplate{
		Type: "cron",
		Schema: `{
			"description": "webhook setting",
			"type": "object",
			"properties": {
				"schedule": {
				"type": "string"
				}
		    }
		}`,
	}
	TriggerWebHook = TriggerConfigTemplate{
		Type: "webhook",
		Schema: `{
			"description": "webhook setting",
			"type": "object",
			"properties": {
				"matchevent": {
				"type": "string"
				},
				"matchbranch": {
				"type": "string"
				}
		    }
		}`,
	}
	AvaliablesTriggerTemplate = []TriggerConfigTemplate{TriggerManual, TriggerCron, TriggerWebHook}
)

// TriggerConfigTemplate is trigger setting template
type TriggerConfigTemplate struct {
	Type   string `bson:"type,omitempty"`
	Schema string `bson:"schema,omitempty"`
}

// TriggerConfig is trigger config in pipeline
type TriggerConfig struct {
	Cron    *CronTrigger
	Mannual *ManualTrigger
	WebHook *WebhookTrigger
}

// Trigger is trigger params in a building
type Trigger interface {
}

// ManualTrigger ..
type ManualTrigger struct {
	Author string
	Params map[string]string
}

// CronTrigger ..
type CronTrigger struct {
	Schedule string
}

// WebhookTrigger ..
type WebhookTrigger struct {
	Repo     string
	Event    string
	Branch   string
	CommitID string
	Comment  string
	Author   string
}

// #todo this is `easy` version of trigger option
type option struct {
	Help    string   `bson:"help,omitempty"`
	Key     string   `bson:"key,omitempty"`
	Default string   `bson:"default,omitempty"`
	Choices []string `bson:"choices,omitempty"`
}
