package config

// ECSConfig holds configuration for the AWS ECS notifier.
type ECSConfig struct {
	Cluster    string `mapstructure:"cluster"`
	TaskDef    string `mapstructure:"task_def"`
	Container  string `mapstructure:"container"`
	LaunchType string `mapstructure:"launch_type"`
	Region     string `mapstructure:"region"`
}
