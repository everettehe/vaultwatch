package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// JiraNotifier creates Jira issues for expiring secrets.
type JiraNotifier struct {
	baseURL  string
	token    string
	project  string
	issueType string
	client  *http.Client
}

type jiraIssueRequest struct {
	Fields jiraFields `json:"fields"`
}

type jiraFields struct {
	Project   jiraKey `json:"project"`
	Summary   string  `json:"summary"`
	Description string `json:"description"`
	IssueType jiraKey `json:"issuetype"`
}

type jiraKey struct {
	Key string `json:"key"`
}

// NewJiraNotifier constructs a JiraNotifier.
func NewJiraNotifier(baseURL, token, project, issueType string) (*JiraNotifier, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("jira: base URL is required")
	}
	if token == "" {
		return nil, fmt.Errorf("jira: API token is required")
	}
	if project == "" {
		return nil, fmt.Errorf("jira: project key is required")
	}
	if issueType == "" {
		issueType = "Task"
	}
	return &JiraNotifier{
		baseURL:   baseURL,
		token:     token,
		project:   project,
		issueType: issueType,
		client:    &http.Client{},
	}, nil
}

// Notify creates a Jira issue for the given secret.
func (j *JiraNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	body := jiraIssueRequest{
		Fields: jiraFields{
			Project:     jiraKey{Key: j.project},
			Summary:     msg.Subject,
			Description: msg.Body,
			IssueType:   jiraKey{Key: j.issueType},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("jira: failed to marshal payload: %w", err)
	}
	url := j.baseURL + "/rest/api/2/issue"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("jira: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+j.token)
	resp, err := j.client.Do(req)
	if err != nil {
		return fmt.Errorf("jira: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("jira: unexpected status %d", resp.StatusCode)
	}
	return nil
}
