import React from "react";
import { Box, Text, useInput, useStdout } from "ink";
import type { WorkflowRun } from "../lib/types.js";
import { StatusBadge } from "./status-badge.js";
import { relativeTime, truncate } from "../lib/format.js";

interface Props {
  runs: WorkflowRun[];
  selectedIndex: number;
  onSelect: (index: number) => void;
  onEnter: (run: WorkflowRun) => void;
  isActive: boolean;
}

export function RunList({ runs, selectedIndex, onSelect, onEnter, isActive }: Props) {
  const { stdout } = useStdout();
  const cols = stdout?.columns ?? 120;

  useInput(
    (input, key) => {
      if (key.upArrow) {
        onSelect(Math.max(0, selectedIndex - 1));
      } else if (key.downArrow) {
        onSelect(Math.min(runs.length - 1, selectedIndex + 1));
      } else if (key.return) {
        const run = runs[selectedIndex];
        if (run) onEnter(run);
      }
    },
    { isActive },
  );

  if (runs.length === 0) {
    return <Text dimColor>No workflow runs found.</Text>;
  }

  // Progressive column hiding based on terminal width
  // Full: selector(2) + status(14) + workflow(20) + branch(18) + title(flex) + time(8) ~ 62+ cols
  const showTime = cols >= 70;
  const showBranch = cols >= 55;
  const showWorkflow = cols >= 40;
  const titleMax = Math.max(15, cols - (2 + 14 + (showWorkflow ? 21 : 0) + (showBranch ? 19 : 0) + (showTime ? 9 : 0)));

  return (
    <Box flexDirection="column">
      {runs.map((run, i) => {
        const selected = i === selectedIndex;
        return (
          <Box key={run.databaseId} gap={1}>
            <Text>{selected ? ">" : " "}</Text>
            <Box width={14}>
              <StatusBadge status={run.status} conclusion={run.conclusion} />
            </Box>
            {showWorkflow && (
              <Box width={20}>
                <Text color="blue">{truncate(run.workflowName, 20)}</Text>
              </Box>
            )}
            {showBranch && (
              <Box width={18}>
                <Text color="magenta">{truncate(run.headBranch, 18)}</Text>
              </Box>
            )}
            <Box flexGrow={1}>
              <Text>{truncate(run.displayTitle, titleMax)}</Text>
            </Box>
            {showTime && (
              <Box width={8}>
                <Text dimColor>{relativeTime(run.createdAt)}</Text>
              </Box>
            )}
          </Box>
        );
      })}
    </Box>
  );
}
