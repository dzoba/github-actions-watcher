package gh

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/dzoba/github-actions-watcher/internal/types"
)

const runFields = "databaseId,displayTitle,event,headBranch,name,number,status,conclusion,createdAt,updatedAt,url,workflowName"

// FetchRepoList returns recently-pushed repos for the authenticated user.
func FetchRepoList() ([]types.PickerRepo, error) {
	out, err := exec.Command("gh", "repo", "list",
		"--json", "nameWithOwner,pushedAt",
		"--limit", "20",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("gh repo list failed: %w", err)
	}
	var repos []types.PickerRepo
	if err := json.Unmarshal(out, &repos); err != nil {
		return nil, fmt.Errorf("failed to parse repo list: %w", err)
	}
	return repos, nil
}

// FetchRuns returns the 20 most recent workflow runs for a repo.
func FetchRuns(repo string) ([]types.WorkflowRun, error) {
	out, err := exec.Command("gh", "run", "list",
		"--repo", repo,
		"--json", runFields,
		"--limit", "20",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("gh run list failed: %w", err)
	}
	var runs []types.WorkflowRun
	if err := json.Unmarshal(out, &runs); err != nil {
		return nil, fmt.Errorf("failed to parse runs: %w", err)
	}
	return runs, nil
}

// FetchRunDetail returns a single run with its jobs and steps.
func FetchRunDetail(repo string, runID int) (*types.RunDetail, error) {
	out, err := exec.Command("gh", "run", "view",
		strconv.Itoa(runID),
		"--repo", repo,
		"--json", runFields+",jobs",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("gh run view failed: %w", err)
	}
	var detail types.RunDetail
	if err := json.Unmarshal(out, &detail); err != nil {
		return nil, fmt.Errorf("failed to parse run detail: %w", err)
	}
	return &detail, nil
}
