"use client";

import { useEffect, useMemo, useState } from "react";
import { Briefcase, Cpu, Landmark, type LucideIcon, Radio, Send, Trophy, Vote, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";
import Badge from "@/components/ui/Badge";
import SectionLabel from "@/components/ui/SectionLabel";
import type { Segment, TrendingClaim } from "@/lib/types";
import { TRENDING_DATA, type TrendingKey } from "@/lib/sandbox";
import { APIError } from "@/lib/api";

interface ClaimFeedProps {
  segment: Segment;
  onSelectClaim: (claim: string, data: TrendingClaim) => void;
  onBack?: () => void;
}

const iconMap: Record<string, LucideIcon> = {
  Cpu,
  Landmark,
  Vote,
  Trophy,
  Briefcase,
};

const colorMap: Record<string, string> = {
  gold: "#D4A74F",
  red: "#DC2626",
  blu: "#2563EB",
  grn: "#059669",
  amb: "#D97706",
};

function verdictVariant(verdict: string): "FALSE" | "MISLEADING" | "TRUE" | "UNVERIFIABLE" {
  const upper = verdict.toUpperCase();
  if (upper === "FALSE") return "FALSE";
  if (upper === "MISLEADING") return "MISLEADING";
  if (upper === "TRUE") return "TRUE";
  return "UNVERIFIABLE";
}

function severityColor(label: string): string {
  const upper = label.toUpperCase();
  if (upper === "CRITICAL") return "#DC2626";
  if (upper === "HIGH") return "#D97706";
  if (upper === "MEDIUM") return "#2563EB";
  return "#059669";
}

function mapSegmentToKey(segment: Segment): TrendingKey {
  const id = segment.id.toLowerCase();
  if (id.includes("tech")) return "tech";
  if (id.includes("finance")) return "finance";
  if (id.includes("politics")) return "politics";
  if (id.includes("sports")) return "sports";
  return "tech";
}

function normalizeTrending(raw: any): TrendingClaim {
  return {
    claim: String(raw.claim ?? raw.claim_text ?? ""),
    verdict: String(raw.verdict ?? "UNVERIFIABLE"),
    confidence: Number(raw.confidence ?? 0),
    risk_score: Number(raw.risk_score ?? 0),
    risk_label: String(raw.risk_label ?? raw.risk_level ?? "LOW"),
    reason: String(raw.reason ?? ""),
    sources: Array.isArray(raw.sources)
      ? raw.sources.map((x: unknown) => String(x))
      : Array.isArray(raw.sources_checked)
        ? raw.sources_checked.map((x: unknown) => String(x))
        : [],
    heat: Math.max(10, Math.min(100, Number(raw.heat ?? raw.virality ?? 0))),
    platform: String(raw.platform ?? "Unknown"),
  };
}

async function fetchTrendingClaims(segment: Segment): Promise<TrendingClaim[]> {
  const segmentKey = mapSegmentToKey(segment);
  const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

  try {
    const response = await fetch(
      `${apiUrl}/api/trending?segment=${encodeURIComponent(segmentKey)}`,
      { cache: "no-store" },
    );
    if (!response.ok) {
      throw new APIError(`Trending API unavailable (${response.status})`, "HTTP_ERROR", response.status);
    }
    const payload = await response.json();
    if (!Array.isArray(payload)) {
      throw new APIError("Trending API returned invalid format", "PARSE_ERROR");
    }
    const normalized = payload.map(normalizeTrending).filter((item) => item.claim.length > 0);
    if (normalized.length > 0) return normalized.slice(0, 6);
  } catch {
    // Fall back to sandbox data.
  }

  return TRENDING_DATA[segmentKey];
}

export default function ClaimFeed({ segment, onSelectClaim, onBack }: ClaimFeedProps) {
  const [customClaim, setCustomClaim] = useState("");
  const [claims, setClaims] = useState<TrendingClaim[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeCard, setActiveCard] = useState<string | null>(null);

  const segmentColor = useMemo(() => colorMap[segment.color] ?? "#D4A74F", [segment.color]);
  const Icon = iconMap[segment.icon] ?? Cpu;

  useEffect(() => {
    let cancelled = false;

    async function loadClaims() {
      setLoading(true);
      const data = await fetchTrendingClaims(segment);
      if (!cancelled) {
        setClaims(data.slice(0, 6));
        setLoading(false);
      }
    }

    loadClaims();
    return () => {
      cancelled = true;
    };
  }, [segment]);

  const pickClaim = (item: TrendingClaim, idx: number) => {
    const key = `${item.claim}-${idx}`;
    setActiveCard(key);
    window.setTimeout(() => {
      onSelectClaim(item.claim, item);
      setActiveCard(null);
    }, 170);
  };

  const submitCustom = () => {
    const claim = customClaim.trim();
    if (!claim) return;
    const synthetic: TrendingClaim = {
      claim,
      verdict: "UNVERIFIABLE",
      confidence: 0,
      risk_score: 0,
      risk_label: "LOW",
      reason: "Manual verification request from operator.",
      sources: [],
      heat: 0,
      platform: "Operator Input",
    };
    onSelectClaim(claim, synthetic);
    setCustomClaim("");
  };

  return (
    <section id="trending" className="space-y-4">
      <SectionLabel>Narrative Discovery</SectionLabel>

      <button
        type="button"
        onClick={onBack}
        className="inline-flex items-center gap-1.5 font-mono text-[10px] uppercase font-bold tracking-[0.18em] text-[#9CA3AF] transition hover:text-[#111827]"
      >
        ← Change Segment
      </button>

      <div className="flex flex-wrap items-center justify-between gap-3">
        <h2 className="flex items-center gap-3 font-heading text-2xl font-bold text-[#111827]">
          <span
            className="inline-flex h-10 w-10 items-center justify-center rounded-xl border border-[#F3F4F6] bg-[#FAFAFB]"
            style={{ color: segmentColor }}
          >
            <Icon className="h-5 w-5" />
          </span>
          Trending in {segment.label || "Segment"}
          <span className="font-mono text-[10px] font-bold uppercase tracking-[0.16em] text-[#9CA3AF]">
            · {claims.length} claims detected
          </span>
        </h2>
        <span className="inline-flex items-center gap-1.5 rounded-full border border-[#D1FAE5] bg-[#ECFDF5] px-3 py-1.5 font-mono text-[9px] font-bold uppercase tracking-[0.2em] text-[#059669]">
          <Radio className="h-3 w-3 animate-pulse" />
          Monitoring Active
        </span>
      </div>

      {loading ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, index) => (
            <div
              key={`skeleton-${index}`}
              className="h-[160px] animate-pulse rounded-2xl border border-[#F3F4F6] bg-white shadow-sm"
            />
          ))}
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {claims.map((item, idx) => {
            const cardKey = `${item.claim}-${idx}`;
            const sourceCredit = item.sources[0] ?? "Verified Network";
            return (
              <motion.button
                key={cardKey}
                type="button"
                whileHover={{ y: -4, borderColor: "#D4A74F" }}
                animate={activeCard === cardKey ? { scale: 1.02, borderColor: "#D4A74F" } : { scale: 1 }}
                transition={{ duration: 0.2 }}
                onClick={() => pickClaim(item, idx)}
                className="group relative overflow-hidden rounded-2xl border border-[#E5E7EB] bg-white p-5 text-left shadow-sm transition-all hover:shadow-md"
              >
                <div className="mb-4 flex items-start justify-between gap-4">
                   <Badge variant={verdictVariant(item.verdict)}>{item.verdict}</Badge>
                   <div className="text-right">
                      <div className="font-heading text-2xl font-extrabold" style={{ color: severityColor(item.risk_label) }}>
                        {Math.round(item.heat)}
                      </div>
                      <div className="font-mono text-[8px] font-bold uppercase tracking-widest text-[#9CA3AF]">Severity</div>
                   </div>
                </div>

                <p className="line-clamp-3 font-sans text-[14px] font-medium leading-relaxed text-[#374151] group-hover:text-[#111827]">
                  {item.claim}
                </p>

                <div className="mt-6 flex items-center justify-between border-t border-[#F3F4F6] pt-4">
                  <span className="font-mono text-[9px] font-bold uppercase tracking-[0.18em] text-[#9CA3AF]">
                    {sourceCredit}
                  </span>
                  <div className="h-6 w-6 rounded-full bg-[#FAFAFB] flex items-center justify-center text-[#9CA3AF] group-hover:text-[#D4A74F] group-hover:bg-[#D4A74F]/10 transition-colors">
                     <ArrowRight size={12} />
                  </div>
                </div>
              </motion.button>
            );
          })}
        </div>
      )}

      <div className="flex flex-col gap-3 sm:flex-row pt-4">
        <div className="relative flex-1">
          <input
            value={customClaim}
            onChange={(event) => setCustomClaim(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Enter") submitCustom();
            }}
            placeholder="Paste custom claim to verify..."
            className="h-12 w-full rounded-xl border border-[#E5E7EB] bg-white px-4 py-3 text-sm text-[#111827] outline-none transition focus:border-[#D4A74F] focus:ring-1 focus:ring-[#D4A74F]/20"
          />
        </div>
        <button 
          onClick={submitCustom}
          className="inline-flex h-12 items-center justify-center gap-2 rounded-xl bg-[#D4A74F] px-8 text-xs font-bold text-black transition hover:bg-[#c3963e]"
        >
          <Send className="h-4 w-4" />
          VERIFY CUSTOM CLAIM
        </button>
      </div>
    </section>
  );
}
