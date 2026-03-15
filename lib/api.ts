import type { VerifyResponse } from "@/lib/types";
import { TRENDING_DATA } from "@/lib/sandbox";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

type StatsResponse = { total: number; false_count: number; misleading_count: number; memory_size: number };

let connectionState = { connected: false, failCount: 0, lastCheck: 0 };

export class APIError extends Error {
  readonly code: string;
  readonly status?: number;

  constructor(message: string, code: string, status?: number) {
    super(message);
    this.name = "APIError";
    this.code = code;
    this.status = status;
  }
}

async function fetchWithTimeout(url: string, options: RequestInit = {}, ms: number = 8000) {
  const controller = new AbortController();
  const id = setTimeout(() => controller.abort(), ms);
  try {
    const response = await fetch(url, { ...options, signal: controller.signal });
    return response;
  } finally {
    clearTimeout(id);
  }
}

async function withFallback<T>(apiFn: () => Promise<T>, fallbackFn: () => T): Promise<T> {
  try {
    const minTimeGap = 10000; // don't hammer the API if it's down permanently
    if (!connectionState.connected && connectionState.failCount >= 3 && Date.now() - connectionState.lastCheck < minTimeGap) {
      console.warn("Backend marked as down. Using fallback immediately.");
      return fallbackFn();
    }

    connectionState.lastCheck = Date.now();
    const result = await apiFn();
    
    connectionState.connected = true;
    connectionState.failCount = 0;
    return result;
  } catch (err) {
    connectionState.failCount++;
    console.warn(`API request failed (${connectionState.failCount} failures). Using fallback...`, err);
    if (connectionState.failCount >= 3) {
      connectionState.connected = false;
    }
    return fallbackFn();
  }
}

export async function verifyClaimAPI(claim: string): Promise<VerifyResponse> {
  return withFallback(
    async () => {
      const res = await fetchWithTimeout(`${API_URL}/api/verify`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ claim }),
        cache: "no-store",
      }, 8000);
      
      if (!res.ok) throw new Error(`Verify API failed: ${res.status}`);
      return res.json() as Promise<VerifyResponse>;
    },
    () => {
      const allClaims = Object.values(TRENDING_DATA).flat();
      const match = allClaims.find(c => c.claim.toLowerCase().includes(claim.toLowerCase().substring(0, 20)));
      return {
        claim: claim,
        verdict: match?.verdict || "UNVERIFIABLE",
        confidence: match?.confidence || 0,
        risk_score: match?.risk_score || 0.5,
        risk_label: match?.risk_label || "MEDIUM",
        reason: match?.reason || "Sandbox fallback: Verification server unreachable.",
        sources: match?.sources || []
      };
    }
  );
}

export async function getStats(): Promise<StatsResponse> {
  return withFallback(
    async () => {
      const res = await fetchWithTimeout(`${API_URL}/api/stats`, { cache: "no-store" }, 3000);
      if (!res.ok) throw new Error("Stats API failed");
      return res.json() as Promise<StatsResponse>;
    },
    () => ({ total: 0, false_count: 0, misleading_count: 0, memory_size: 0 })
  );
}

export async function getChain() {
  const res = await fetch(`${API_URL}/api/chain`);
  if (!res.ok) throw new Error("Failed to fetch chain");
  return res.json();
}

export function createSSEConnection(
  onMessage: (data: any) => void,
  onError?: () => void
) {
  const source = new EventSource(`${API_URL}/api/stream`);

  source.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.type === "detection" || data.claim) {
        onMessage(data.data || data);
      }
    } catch (err) {
      console.error("SSE parse error", err);
    }
  };

  source.onerror = (err) => {
    console.error("SSE connection error", err);
    if (onError) onError();
  };

  return source;
}
