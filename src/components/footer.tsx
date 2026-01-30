import React from "react";
import { Box, Text } from "ink";
import type { View } from "../lib/types.js";

interface Props {
  view: View;
  intervalSeconds: number;
}

const keys: Record<View, string> = {
  list: "up/down: navigate | enter: details | s: switch repo | r: refresh | q: quit",
  detail: "up/down: scroll | esc: back | o: open in browser | r: refresh | q: quit",
  "repo-input": "enter: confirm | esc: cancel",
};

export function Footer({ view, intervalSeconds }: Props) {
  const showInterval = view === "list" || view === "detail";

  return (
    <Box marginTop={1}>
      <Text dimColor>{keys[view]}</Text>
      {showInterval && (
        <Text dimColor> | auto-refresh: {intervalSeconds}s</Text>
      )}
    </Box>
  );
}
