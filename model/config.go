package model

type Config struct {
}

type TriggerConfig struct {
	TriggerType string `json:"triggerType"`
}

// JobConfig represent pipeline config
// JobConfig contains metadata and workconfig
type JobConfig struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Trigger         []TriggerConfig `json:"trigger"`
	ConcurrentBuild bool            `json:"concurrentBuild"`
	Scm             []WorkConfig    `json:"scm"`
	PreBuild        []WorkConfig    `json:"preBuild"`
	Build           []WorkConfig    `json:"build"`
	AfterBuild      []WorkConfig    `json:"afterBuild"`
	// Publish         []WorkConfig    `json:"publish"`
	Notify  []WorkConfig `json:"notify"`
	Service []WorkConfig `json:"service"`
}

// WorkConfig represent pipeline step config
// WorkConfig contains running data which will be convert into Work
type WorkConfig struct {
	Name            string   `json:"name"`
	Image           string   `json:"image,omitempty"`
	Command         []string `json:"command,omitempty"`
	Args            []string `json:"args,omitempty"`
	Env             []EnvVar `json:"env,omitempty"`
	ImagePullPolicy string   `json:"imagePullPolicy,omitempty"`
	// use secret or configmap
	CredentialsName string `json:"credentialsName,omitempty"`
	CredentialsPath string `json:"credentialsPath,omitempty"`

	// -- auto set if needed --
	// Resources ResourceRequirements `json:"resources,omitempty" `
	// VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
	// Lifecycle       *Lifecycle `json:"lifecycle,omitempty"`
	// WorkingDir      string     `json:"workingDir,omitempty"`

	// -- not used ---
	// SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Stdin bool `json:"stdin,omitempty" `
	// StdinOnce bool `json:"stdinOnce,omitempty" `
	// TTY bool `json:"tty,omitempty" `
	// TerminationMessagePath string `json:"terminationMessagePath,omitempty" `
	// TerminationMessagePolicy TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty"`
	// ReadinessProbe *Probe `json:"readinessProbe,omitempty" `
	// LivenessProbe *Probe `json:"livenessProbe,omitempty" `
	// EnvFrom        []EnvFromSource     `json:"envFrom,omitempty"`
	// Ports []ContainerPort `json:"ports,omitempty"`
}
