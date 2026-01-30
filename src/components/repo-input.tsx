import React, { useState } from "react";
import { Box, Text, useInput } from "ink";
import TextInput from "ink-text-input";

interface Props {
  currentRepo: string;
  onConfirm: (repo: string) => void;
  onCancel: () => void;
  isActive: boolean;
}

export function RepoInput({ currentRepo, onConfirm, onCancel, isActive }: Props) {
  const [value, setValue] = useState(currentRepo);

  useInput(
    (_input, key) => {
      if (key.escape) {
        onCancel();
      }
    },
    { isActive },
  );

  return (
    <Box flexDirection="column" marginTop={1}>
      <Text bold>Switch repository:</Text>
      <Box gap={1}>
        <Text color="cyan">owner/repo:</Text>
        <TextInput
          value={value}
          onChange={setValue}
          onSubmit={(val) => {
            const trimmed = val.trim();
            if (trimmed && trimmed.includes("/")) {
              onConfirm(trimmed);
            }
          }}
          focus={isActive}
        />
      </Box>
    </Box>
  );
}
