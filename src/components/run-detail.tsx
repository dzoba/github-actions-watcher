import React from "react";
import { Box, Text, useInput } from "ink";
import type { RunDetail } from "../lib/types.js";
import { StatusBadge } from "./status-badge.js";
import { duration } from "../lib/format.js";

interface Props {
  detail: RunDetail | null;
  loading: boolean;
  error: string | null;
  onBack: () => void;
  onOpen: () => void;
  onRefresh: () => void;
  isActive: boolean;
  scrollOffset: number;
  onScroll: (offset: number) => void;
}

export function RunDetailView({
  detail,
  loading,
  error,
  onBack,
  onOpen,
  onRefresh,
  isActive,
  scrollOffset,
  onScroll,
}: Props) {
  useInput(
    (input, key) => {
      if (key.escape) {
        onBack();
      } else if (input === "o") {
        onOpen();
      } else if (input === "r") {
        onRefresh();
      } else if (key.upArrow) {
        onScroll(Math.max(0, scrollOffset - 1));
      } else if (key.downArrow) {
        onScroll(scrollOffset + 1);
      }
    },
    { isActive },
  );

  if (loading && !detail) {
    return <Text dimColor>Loading run details...</Text>;
  }

  if (error) {
    return <Text color="red">Error: {error}</Text>;
  }

  if (!detail) {
    return <Text dimColor>No detail available.</Text>;
  }

  // Build flat list of lines to render, then apply scroll offset
  const lines: React.ReactNode[] = [];

  lines.push(
    <Box key="title" gap={1}>
      <StatusBadge status={detail.status} conclusion={detail.conclusion} />
      <Text bold>{detail.displayTitle}</Text>
    </Box>,
  );
  lines.push(
    <Box key="meta" gap={1}>
      <Text dimColor>
        {detail.workflowName} #{detail.number} on {detail.headBranch} ({detail.event})
      </Text>
      {loading && <Text dimColor> fetching...</Text>}
    </Box>,
  );
  lines.push(<Text key="sep" dimColor>{"---"}</Text>);

  for (const job of detail.jobs) {
    lines.push(
      <Box key={`job-${job.databaseId}`} gap={1} marginTop={1}>
        <StatusBadge status={job.status} conclusion={job.conclusion} />
        <Text bold>{job.name}</Text>
        {job.startedAt && (
          <Text dimColor>({duration(job.startedAt, job.completedAt)})</Text>
        )}
      </Box>,
    );

    for (const step of job.steps) {
      lines.push(
        <Box key={`step-${job.databaseId}-${step.number}`} gap={1} paddingLeft={2}>
          <StatusBadge status={step.status} conclusion={step.conclusion} />
          <Text>{step.name}</Text>
          {step.startedAt && step.completedAt && (
            <Text dimColor>({duration(step.startedAt, step.completedAt)})</Text>
          )}
        </Box>,
      );
    }
  }

  const visibleLines = lines.slice(scrollOffset);

  return <Box flexDirection="column">{visibleLines}</Box>;
}
