package config

// CodePipelineConfig holds settings for the AWS CodePipeline notifier.
type CodePipelineConfig struct {
	// JobID is the CodePipeline job ID to report failure against.
	JobID string `yaml:"job_id"`
	// Region is the AWS region where the pipeline resides.
	Region string `yaml:"region"`
}
