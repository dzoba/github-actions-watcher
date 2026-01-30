import React from "react";
import { Text } from "ink";
import { useCountdown } from "../hooks/use-countdown.js";

interface Props {
  intervalMs: number;
  enabled: boolean;
  resetKey: number;
}

export function Countdown({ intervalMs, enabled, resetKey }: Props) {
  const { secondsLeft } = useCountdown(intervalMs, enabled, resetKey);

  if (!enabled) return null;

  return <Text dimColor> | refresh in {secondsLeft}s</Text>;
}
