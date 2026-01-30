import React from "react";
import { Box, Text } from "ink";
import Spinner from "ink-spinner";

interface Props {
  repo: string;
  loading: boolean;
}

export function Header({ repo, loading }: Props) {
  return (
    <Box marginBottom={1}>
      <Text bold color="cyan">
        GitHub Actions
      </Text>
      <Text> - </Text>
      <Text bold>{repo}</Text>
      {loading && (
        <Text>
          {" "}
          <Spinner type="dots" />
        </Text>
      )}
    </Box>
  );
}
