export type RunStatus =
  | "completed"
  | "in_progress"
  | "queued"
  | "requested"
  | "waiting"
  | "pending";

export type RunConclusion =
  | "success"
  | "failure"
  | "cancelled"
  | "skipped"
  | "timed_out"
  | "action_required"
  | "neutral"
  | "stale"
  | null;

export type StepStatus =
  | "completed"
  | "in_progress"
  | "queued"
  | "pending"
  | "waiting";

export type StepConclusion =
  | "success"
  | "failure"
  | "cancelled"
  | "skipped"
  | null;

export interface WorkflowRun {
  databaseId: number;
  displayTitle: string;
  event: string;
  headBranch: string;
  name: string;
  number: number;
  status: RunStatus;
  conclusion: RunConclusion;
  createdAt: string;
  updatedAt: string;
  url: string;
  workflowName: string;
}

export interface Step {
  name: string;
  status: StepStatus;
  conclusion: StepConclusion;
  number: number;
  startedAt: string;
  completedAt: string;
}

export interface Job {
  name: string;
  status: StepStatus;
  conclusion: StepConclusion;
  startedAt: string;
  completedAt: string;
  steps: Step[];
  url: string;
  databaseId: number;
}

export interface RunDetail {
  databaseId: number;
  displayTitle: string;
  event: string;
  headBranch: string;
  name: string;
  number: number;
  status: RunStatus;
  conclusion: RunConclusion;
  createdAt: string;
  updatedAt: string;
  url: string;
  workflowName: string;
  jobs: Job[];
}

export type View = "list" | "detail" | "repo-input";
