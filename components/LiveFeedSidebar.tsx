"use client";

import { useEffect, useMemo, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Activity, Radio } from "lucide-react";
import Badge from "@/components/ui/Badge";
import Card from "@/components/ui/Card";
import ProgressBar from "@/components/ui/ProgressBar";
import SectionLabel from "@/components/ui/SectionLabel";
import { createSSEConnection, getChain, getStats } from "@/lib/api";
import { useConnectionContext } from "@/lib/context";
import { TRENDING_DATA } from "@/lib/sandbox";
import type { TrendingClaim } from "@/lib/types";

interface LiveFeedSidebarProps {
  onPickClaim?: (claim: string) => void;
}

interface FeedItem extends TrendingClaim {
  id: string;
  timestamp: string;
}

const SANDBOX_STREAM: TrendingClaim[] = [
  ...TRENDING_DATA.tech,
  ...TRENDING_DATA.finance,
  ...TRENDING_DATA.politics,
  ...TRENDING_DATA.sports,
];

function verdictVariant(verdict: string): "FALSE" | "MISLEADING" | "TRUE" | "UNVERIFIABLE" {
  const upper = verdict.toUpperCase();
  if (upper === "FALSE") return "FALSE";
  if (upper === "MISLEADING") return "MISLEADING";
  if (upper === "TRUE") return "TRUE";
  return "UNVERIFIABLE";
}

function riskToPercent(score: number): number {
  if (score <= 1) return Math.max(0, Math.min(100, score * 100));
  if (score <= 10) return Math.max(0, Math.min(100, (score / 10) * 100));
  return Math.max(0, Math.min(100, score));
}

function nowLabel() {
  return new Date().toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function normalizeIncoming(data: any): FeedItem {
  const claim = String(data.claim_text ?? data.claim ?? "Incoming claim");
  const verdict = String(data.verdict ?? "UNVERIFIABLE");
  const confidence = Number(data.confidence ?? 0);
  const riskScore = Number(data.risk_score ?? 0);
  const heatFromVirality = Math.max(12, Math.min(100, Number(data.virality ?? 0) / 1500));

  return {
    id: String(data.id ?? data.claim_id ?? `${claim}-${Date.now()}`),
    claim,
    verdict,
    confidence,
    risk_score: riskScore,
    risk_label: String(data.risk_label ?? data.risk_level ?? "LOW"),
    reason: String(data.reason ?? "Live integrity stream event"),
    sources: Array.isArray(data.sources)
      ? data.sources.map((x: unknown) => String(x))
      : Array.isArray(data.sources_checked)
        ? data.sources_checked.map((x: unknown) => String(x))
        : [],
    heat: Number.isFinite(Number(data.heat)) ? Number(data.heat) : heatFromVirality,
    platform: String(data.platform ?? "Unknown"),
    timestamp: data.timestamp ? new Date(data.timestamp).toLocaleTimeString() : nowLabel(),
  };
}

function toFeedItem(item: TrendingClaim, seed: number): FeedItem {
  return {
    ...item,
    id: `${item.platform}-${seed}`,
    timestamp: nowLabel(),
  };
}

export default function LiveFeedSidebar({ onPickClaim }: LiveFeedSidebarProps) {
  const [feed, setFeed] = useState<FeedItem[]>([]);
  const [sandboxMode, setSandboxMode] = useState(false);
  const { connected, setConnected, total, falseCount, chainCount, setMetrics } = useConnectionContext();

  useEffect(() => {
    let source: EventSource | null = null;
    let statsTimer: number | undefined;
    let sandboxTimer: number | undefined;
    let sandboxCursor = 0;
    let cancelled = false;

    const pushItem = (item: FeedItem) => {
      setFeed((prev) => [item, ...prev].slice(0, 12));
    };

    const syncStats = async () => {
      try {
        const [stats, chain] = await Promise.all([getStats(), getChain()]);
        if (cancelled) return;
        setMetrics({
          total: Number(stats.total ?? 0),
          falseCount: Number(stats.false_count ?? 0),
          chainCount: Array.isArray(chain) ? chain.length : Number(chain?.length ?? 0),
        });
      } catch {
        // Keep previous values.
      }
    };

    const startSandboxReplay = () => {
      if (sandboxTimer) return;
      setSandboxMode(true);
      sandboxTimer = window.setInterval(() => {
        const next = SANDBOX_STREAM[sandboxCursor % SANDBOX_STREAM.length];
        sandboxCursor += 1;
        pushItem(toFeedItem(next, sandboxCursor));
      }, 2200);
    };

    const stopSandboxReplay = () => {
      if (sandboxTimer) {
        window.clearInterval(sandboxTimer);
        sandboxTimer = undefined;
      }
      setSandboxMode(false);
    };

    const connectSSE = () => {
      try {
        source = createSSEConnection(
          (payload) => {
            if (cancelled) return;
            setConnected(true);
            stopSandboxReplay();
            pushItem(normalizeIncoming(payload));
          },
          () => {
            if (cancelled) return;
            setConnected(false);
            if (source && typeof source.close === "function") {
              source.close();
              source = null;
            }
            startSandboxReplay();
          },
        );
      } catch {
        setConnected(false);
        startSandboxReplay();
      }
    };

    connectSSE();
    syncStats();
    statsTimer = window.setInterval(syncStats, 5000);

    return () => {
      cancelled = true;
      if (source && typeof source.close === "function") source.close();
      if (statsTimer) window.clearInterval(statsTimer);
      if (sandboxTimer) window.clearInterval(sandboxTimer);
    };
  }, [setConnected, setMetrics]);

  const misleadingCount = useMemo(
    () => feed.filter((x) => x.verdict.toUpperCase() === "MISLEADING").length,
    [feed],
  );

  const totalDisplay = total > 0 ? total : feed.length;

  return (
    <aside id="live-feed" className="space-y-4">
      <SectionLabel>Live Intercepts</SectionLabel>

      <Card className="p-4 sm:p-5">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="font-heading text-lg font-bold text-[#111827]">Integrity Feed</h2>
            <p className="font-sans text-[10px] font-semibold uppercase tracking-wider text-[#9CA3AF]">
              {feed.length} Active Stream
            </p>
          </div>
          <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider ${connected ? "bg-[#D1FAE5] text-[#059669]" : "bg-[#F3F4F6] text-[#6B7280]"}`}>
            <span className={`h-1.5 w-1.5 rounded-full ${connected ? "bg-[#059669] animate-pulse" : "bg-[#9CA3AF]"}`} />
            {connected ? "LIVE" : sandboxMode ? "Mock" : "Offline"}
          </span>
        </div>

        <div className="space-y-3 max-h-[400px] overflow-y-auto pr-1 scrollbar-thin">
          <AnimatePresence initial={false}>
            {feed.map((item) => (
              <motion.button
                key={item.id}
                type="button"
                layout
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: 0.2 }}
                onClick={() => onPickClaim?.(item.claim)}
                className="group w-full rounded-xl border border-[#E5E7EB] bg-[#FAFAFB] p-3 text-left transition hover:border-[#D1D5DB] hover:bg-white hover:shadow-md"
              >
                <p className="line-clamp-2 text-xs font-medium leading-relaxed text-[#374151] group-hover:text-[#111827]">{item.claim}</p>
                <div className="mt-3 flex items-center justify-between">
                  <Badge variant={verdictVariant(item.verdict)}>{item.verdict}</Badge>
                  <span className="font-mono text-[9px] font-medium text-[#9CA3AF]">
                    {item.timestamp}
                  </span>
                </div>
                <div className="mt-3">
                  <ProgressBar
                    value={riskToPercent(item.risk_score)}
                    color={item.verdict.toUpperCase() === "TRUE" ? "#059669" : "#DC2626"}
                    className="h-1"
                  />
                </div>
              </motion.button>
            ))}
          </AnimatePresence>
        </div>

        <div className="grid grid-cols-2 gap-3 border-t border-[#E5E7EB] mt-6 pt-5">
          <div className="rounded-lg bg-[#F9FAFB] p-3 border border-[#F3F4F6]">
            <p className="text-[10px] font-bold uppercase tracking-wider text-[#9CA3AF]">Total</p>
            <p className="mt-1 text-base font-bold text-[#111827] font-heading">{totalDisplay}</p>
          </div>
          <div className="rounded-lg bg-[#F9FAFB] p-3 border border-[#F3F4F6]">
            <p className="text-[10px] font-bold uppercase tracking-wider text-[#DC2626]">False</p>
            <p className="mt-1 text-base font-bold text-[#DC2626] font-heading">{falseCount}</p>
          </div>
          <div className="rounded-lg bg-[#F9FAFB] p-3 border border-[#F3F4F6]">
            <p className="text-[10px] font-bold uppercase tracking-wider text-[#D97706]">Warning</p>
            <p className="mt-1 text-base font-bold text-[#D97706] font-heading">{misleadingCount}</p>
          </div>
          <div className="rounded-lg bg-[#F9FAFB] p-3 border border-[#F3F4F6]">
            <p className="text-[10px] font-bold uppercase tracking-wider text-[#2563EB]">Chain</p>
            <p className="mt-1 flex items-center gap-1.5 text-base font-bold text-[#2563EB] font-heading">
              <Activity className="h-4 w-4" />
              {chainCount}
            </p>
          </div>
        </div>
      </Card>
    </aside>
  );
}
