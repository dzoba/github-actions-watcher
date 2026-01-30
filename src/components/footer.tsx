import React from "react";
import { Box, Text } from "ink";
import type { View } from "../lib/types.js";

interface Props {
  view: View;
}

const keys: Record<View, string> = {
  list: "up/down: navigate | enter: details | s: switch repo | r: refresh | q: quit",
  detail: "up/down: scroll | esc: back | o: open in browser | r: refresh | q: quit",
  "repo-input": "enter: confirm | esc: cancel",
};

export function Footer({ view }: Props) {
  return (
    <Box marginTop={1}>
      <Text dimColor>{keys[view]}</Text>
    </Box>
  );
}
