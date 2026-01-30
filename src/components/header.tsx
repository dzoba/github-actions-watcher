import React from "react";
import { Box, Text } from "ink";

interface Props {
  repo: string;
}

export function Header({ repo }: Props) {
  return (
    <Box marginBottom={1}>
      <Text bold color="cyan">
        GitHub Actions
      </Text>
      <Text> - </Text>
      <Text bold>{repo}</Text>
    </Box>
  );
}
