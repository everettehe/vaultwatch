package config

// JiraConfig holds configuration for the Jira notifier.
type JiraConfig struct {
	// BaseURL is the base URL of the Jira instance, e.g. https://myorg.atlassian.net
	BaseURL string `yaml:"base_url"`
	// Token is the Jira API token used for authentication.
	Token string `yaml:"token"`
	// Project is the Jira project key where issues will be created.
	Project string `yaml:"project"`
	// IssueType is the type of Jira issue to create. Defaults to "Task".
	IssueType string `yaml:"issue_type"`
}
