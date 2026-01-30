import { useState, useEffect, useCallback, useRef } from "react";

export function usePolling<T>(
  fetcher: () => Promise<T>,
  intervalMs: number,
  enabled: boolean,
) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const mountedRef = useRef(true);
  const lastJsonRef = useRef<string>("");
  const initialFetchDoneRef = useRef(false);

  const refresh = useCallback(async () => {
    // Only show loading state on first fetch, not background polls
    if (!initialFetchDoneRef.current) {
      setLoading(true);
    }
    try {
      const result = await fetcher();
      if (!mountedRef.current) return;

      // Only update state if data actually changed
      const json = JSON.stringify(result);
      if (json !== lastJsonRef.current) {
        lastJsonRef.current = json;
        setData(result);
      }
      if (error !== null) setError(null);
    } catch (err) {
      if (mountedRef.current) {
        setError(err instanceof Error ? err.message : "Fetch failed");
      }
    } finally {
      if (mountedRef.current) {
        if (!initialFetchDoneRef.current) {
          initialFetchDoneRef.current = true;
          setLoading(false);
        }
      }
    }
  }, [fetcher, error]);

  useEffect(() => {
    mountedRef.current = true;
    return () => {
      mountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    if (!enabled) return;
    initialFetchDoneRef.current = false;
    lastJsonRef.current = "";
    refresh();
    const id = setInterval(refresh, intervalMs);
    return () => clearInterval(id);
  }, [enabled, intervalMs, refresh]);

  return { data, setData, loading, error, refresh };
}
