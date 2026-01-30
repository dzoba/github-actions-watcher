import { useState, useEffect } from "react";
import { detectRepo } from "../lib/gh.js";

export function useRepo() {
  const [repo, setRepo] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    detectRepo()
      .then((r) => {
        setRepo(r);
        setLoading(false);
      })
      .catch((err) => {
        setError(
          err instanceof Error ? err.message : "Failed to detect repo",
        );
        setLoading(false);
      });
  }, []);

  return { repo, setRepo, error, loading };
}
