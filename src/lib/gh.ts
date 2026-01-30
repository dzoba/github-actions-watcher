import { execFile } from "node:child_process";
import { promisify } from "node:util";
import type { WorkflowRun, RunDetail } from "./types.js";

const execFileAsync = promisify(execFile);

async function gh<T>(args: string[]): Promise<T> {
  const { stdout } = await execFileAsync("gh", args, {
    maxBuffer: 10 * 1024 * 1024,
  });
  return JSON.parse(stdout) as T;
}

const RUN_FIELDS = [
  "databaseId",
  "displayTitle",
  "event",
  "headBranch",
  "name",
  "number",
  "status",
  "conclusion",
  "createdAt",
  "updatedAt",
  "url",
  "workflowName",
].join(",");

export async function fetchRuns(repo: string): Promise<WorkflowRun[]> {
  return gh<WorkflowRun[]>([
    "run",
    "list",
    "--repo",
    repo,
    "--json",
    RUN_FIELDS,
    "--limit",
    "20",
  ]);
}

export async function fetchRunDetail(
  repo: string,
  runId: number,
): Promise<RunDetail> {
  return gh<RunDetail>([
    "run",
    "view",
    String(runId),
    "--repo",
    repo,
    "--json",
    `${RUN_FIELDS},jobs`,
  ]);
}

export async function detectRepo(): Promise<string> {
  const { stdout } = await execFileAsync("git", [
    "remote",
    "get-url",
    "origin",
  ]);
  const url = stdout.trim();
  // Handle SSH: git@github.com:owner/repo.git
  const sshMatch = url.match(/git@github\.com:(.+?)(?:\.git)?$/);
  if (sshMatch) return sshMatch[1]!;
  // Handle HTTPS: https://github.com/owner/repo.git
  const httpsMatch = url.match(/github\.com\/(.+?)(?:\.git)?$/);
  if (httpsMatch) return httpsMatch[1]!;
  throw new Error(`Could not parse repo from remote URL: ${url}`);
}
