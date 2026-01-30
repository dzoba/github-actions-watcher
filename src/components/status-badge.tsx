import React from "react";
import { Text } from "ink";
import type { RunStatus, RunConclusion, StepStatus, StepConclusion } from "../lib/types.js";

interface Props {
  status: RunStatus | StepStatus;
  conclusion: RunConclusion | StepConclusion;
}

export function StatusBadge({ status, conclusion }: Props) {
  if (status === "in_progress") {
    return <Text color="yellow">{"* running"}</Text>;
  }
  if (status === "queued" || status === "waiting" || status === "pending" || status === "requested") {
    return <Text color="gray">{"~ queued"}</Text>;
  }

  switch (conclusion) {
    case "success":
      return <Text color="green">{"+ passed"}</Text>;
    case "failure":
      return <Text color="red">{"x failed"}</Text>;
    case "cancelled":
      return <Text color="gray">{"- cancelled"}</Text>;
    case "skipped":
      return <Text color="gray">{"- skipped"}</Text>;
    case "timed_out":
      return <Text color="red">{"! timed out"}</Text>;
    case "action_required":
      return <Text color="yellow">{"! action req"}</Text>;
    default:
      return <Text color="gray">{"? unknown"}</Text>;
  }
}
