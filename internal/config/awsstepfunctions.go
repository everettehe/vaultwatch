package config

// StepFunctionsConfig holds configuration for the AWS Step Functions notifier.
type StepFunctionsConfig struct {
	// StateMachineARN is the ARN of the Step Functions state machine to trigger.
	StateMachineARN string `yaml:"state_machine_arn"`
}
