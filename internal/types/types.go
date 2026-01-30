package types

// View represents the current screen.
type View int

const (
	ViewList View = iota
	ViewDetail
	ViewRepoInput
	ViewWelcome
)

// RunStatus is the status of a workflow run.
type RunStatus string

const (
	StatusCompleted  RunStatus = "completed"
	StatusInProgress RunStatus = "in_progress"
	StatusQueued     RunStatus = "queued"
	StatusRequested  RunStatus = "requested"
	StatusWaiting    RunStatus = "waiting"
	StatusPending    RunStatus = "pending"
)

// RunConclusion is the final result of a completed run.
type RunConclusion string

const (
	ConclusionSuccess        RunConclusion = "success"
	ConclusionFailure        RunConclusion = "failure"
	ConclusionCancelled      RunConclusion = "cancelled"
	ConclusionSkipped        RunConclusion = "skipped"
	ConclusionTimedOut       RunConclusion = "timed_out"
	ConclusionActionRequired RunConclusion = "action_required"
	ConclusionNeutral        RunConclusion = "neutral"
	ConclusionStale          RunConclusion = "stale"
)

// WorkflowRun represents a single workflow run from gh CLI.
type WorkflowRun struct {
	DatabaseID   int           `json:"databaseId"`
	DisplayTitle string        `json:"displayTitle"`
	Event        string        `json:"event"`
	HeadBranch   string        `json:"headBranch"`
	Name         string        `json:"name"`
	Number       int           `json:"number"`
	Status       RunStatus     `json:"status"`
	Conclusion   RunConclusion `json:"conclusion"`
	CreatedAt    string        `json:"createdAt"`
	UpdatedAt    string        `json:"updatedAt"`
	URL          string        `json:"url"`
	WorkflowName string        `json:"workflowName"`
}

// Step represents a single step within a job.
type Step struct {
	Name        string        `json:"name"`
	Status      RunStatus     `json:"status"`
	Conclusion  RunConclusion `json:"conclusion"`
	Number      int           `json:"number"`
	StartedAt   string        `json:"startedAt"`
	CompletedAt string        `json:"completedAt"`
}

// Job represents a single job within a workflow run.
type Job struct {
	Name        string        `json:"name"`
	Status      RunStatus     `json:"status"`
	Conclusion  RunConclusion `json:"conclusion"`
	StartedAt   string        `json:"startedAt"`
	CompletedAt string        `json:"completedAt"`
	Steps       []Step        `json:"steps"`
	URL         string        `json:"url"`
	DatabaseID  int           `json:"databaseId"`
}

// RunDetail is a WorkflowRun with nested jobs.
type RunDetail struct {
	WorkflowRun
	Jobs []Job `json:"jobs"`
}
