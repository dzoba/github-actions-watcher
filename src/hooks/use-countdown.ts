import { useState, useEffect, useRef } from "react";

export function useCountdown(intervalMs: number, enabled: boolean, resetKey: number) {
  const [secondsLeft, setSecondsLeft] = useState(Math.floor(intervalMs / 1000));
  const lastResetRef = useRef(Date.now());

  useEffect(() => {
    if (!enabled) return;

    lastResetRef.current = Date.now();
    setSecondsLeft(Math.floor(intervalMs / 1000));

    const tick = setInterval(() => {
      const elapsed = Date.now() - lastResetRef.current;
      const remaining = Math.max(0, Math.ceil((intervalMs - elapsed) / 1000));
      setSecondsLeft(remaining);

      if (remaining === 0) {
        lastResetRef.current = Date.now();
        setSecondsLeft(Math.floor(intervalMs / 1000));
      }
    }, 1000);

    return () => clearInterval(tick);
  }, [intervalMs, enabled, resetKey]);

  return { secondsLeft };
}
