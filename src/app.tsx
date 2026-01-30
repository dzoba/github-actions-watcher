import React, { useState, useCallback, useEffect } from "react";
import { Box, Text, useApp, useInput } from "ink";
import Spinner from "ink-spinner";
import type { View, WorkflowRun, RunDetail } from "./lib/types.js";
import { fetchRuns, fetchRunDetail } from "./lib/gh.js";
import { useRepo } from "./hooks/use-repo.js";
import { usePolling } from "./hooks/use-polling.js";
import { Header } from "./components/header.js";
import { Footer } from "./components/footer.js";
import { RunList } from "./components/run-list.js";
import { RunDetailView } from "./components/run-detail.js";
import { RepoInput } from "./components/repo-input.js";

interface Props {
  interval: number;
}

export function App({ interval }: Props) {
  const { exit } = useApp();
  const { repo, setRepo, error: repoError, loading: repoLoading } = useRepo();

  const [view, setView] = useState<View>("list");
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [selectedRunId, setSelectedRunId] = useState<number | null>(null);
  const [detailScrollOffset, setDetailScrollOffset] = useState(0);

  // Detail state
  const [detail, setDetail] = useState<RunDetail | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailError, setDetailError] = useState<string | null>(null);

  const runsFetcher = useCallback(() => {
    if (!repo) return Promise.resolve([]);
    return fetchRuns(repo);
  }, [repo]);

  const {
    data: runs,
    setData: setRuns,
    loading: runsLoading,
    error: runsError,
    refresh: refreshRuns,
  } = usePolling<WorkflowRun[]>(runsFetcher, interval, !!repo);

  // Fetch detail
  const fetchDetail = useCallback(async () => {
    if (!repo || selectedRunId === null) return;
    setDetailLoading(true);
    setDetailError(null);
    try {
      const d = await fetchRunDetail(repo, selectedRunId);
      setDetail(d);
    } catch (err) {
      setDetailError(err instanceof Error ? err.message : "Failed to fetch detail");
    } finally {
      setDetailLoading(false);
    }
  }, [repo, selectedRunId]);

  // Auto-fetch detail when entering detail view
  useEffect(() => {
    if (view === "detail" && selectedRunId !== null) {
      fetchDetail();
    }
  }, [view, selectedRunId, fetchDetail]);

  // Poll detail view
  useEffect(() => {
    if (view !== "detail" || selectedRunId === null) return;
    const id = setInterval(fetchDetail, interval);
    return () => clearInterval(id);
  }, [view, selectedRunId, interval, fetchDetail]);

  // Global key handler for keys that don't belong to child components
  useInput(
    (input, key) => {
      if (input === "q") {
        exit();
      } else if (input === "s" && view === "list") {
        setView("repo-input");
      } else if (input === "r" && view === "list") {
        refreshRuns();
      }
    },
    { isActive: view === "list" || view === "detail" },
  );

  const handleDrillIn = useCallback((run: WorkflowRun) => {
    setSelectedRunId(run.databaseId);
    setDetail(null);
    setDetailScrollOffset(0);
    setView("detail");
  }, []);

  const handleBack = useCallback(() => {
    setView("list");
    setSelectedRunId(null);
    setDetail(null);
  }, []);

  const handleOpenInBrowser = useCallback(async () => {
    if (!detail?.url) return;
    const openModule = await import("open");
    openModule.default(detail.url);
  }, [detail]);

  const handleRepoConfirm = useCallback(
    (newRepo: string) => {
      setRepo(newRepo);
      setRuns(null);
      setSelectedIndex(0);
      setView("list");
    },
    [setRepo, setRuns],
  );

  const handleRepoCancel = useCallback(() => {
    setView("list");
  }, []);

  // Loading state
  if (repoLoading) {
    return (
      <Text>
        <Spinner type="dots" /> Detecting repository...
      </Text>
    );
  }

  if (repoError && !repo) {
    return (
      <Box flexDirection="column">
        <Text color="red">Error: {repoError}</Text>
        <Text dimColor>Run from a directory with a git remote, or press s to enter a repo manually.</Text>
        <RepoInput
          currentRepo=""
          onConfirm={(r) => {
            setRepo(r);
            setView("list");
          }}
          onCancel={() => exit()}
          isActive={true}
        />
      </Box>
    );
  }

  if (!repo) return null;

  return (
    <Box flexDirection="column">
      <Header repo={repo} loading={runsLoading || detailLoading} />

      {runsError && view === "list" && (
        <Text color="red">Error: {runsError}</Text>
      )}

      {view === "list" && runs && (
        <RunList
          runs={runs}
          selectedIndex={selectedIndex}
          onSelect={setSelectedIndex}
          onEnter={handleDrillIn}
          isActive={view === "list"}
        />
      )}

      {view === "detail" && (
        <RunDetailView
          detail={detail}
          loading={detailLoading}
          error={detailError}
          onBack={handleBack}
          onOpen={handleOpenInBrowser}
          onRefresh={fetchDetail}
          isActive={view === "detail"}
          scrollOffset={detailScrollOffset}
          onScroll={setDetailScrollOffset}
        />
      )}

      {view === "repo-input" && (
        <RepoInput
          currentRepo={repo}
          onConfirm={handleRepoConfirm}
          onCancel={handleRepoCancel}
          isActive={view === "repo-input"}
        />
      )}

      <Footer view={view} />
    </Box>
  );
}
