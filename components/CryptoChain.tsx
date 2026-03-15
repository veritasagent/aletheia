"use client";

import { Copy, Download, ShieldCheck } from "lucide-react";
import Card from "@/components/ui/Card";
import GoldButton from "@/components/ui/GoldButton";
import SectionLabel from "@/components/ui/SectionLabel";
import type { VerifyResponse } from "@/lib/types";

interface CryptoChainProps {
  result: VerifyResponse;
  onCopyHash: () => void;
  onDownloadEvidence: () => void;
  onVerifyIntegrity: () => void;
}

function resolvedHash(result: VerifyResponse) {
  return (
    result.evidence_hash ??
    result.current_hash ??
    result.hash ??
    "0000000000000000000000000000000000000000000000000000000000000000"
  );
}

export default function CryptoChain({
  result,
  onCopyHash,
  onDownloadEvidence,
  onVerifyIntegrity,
}: CryptoChainProps) {
  const hash = resolvedHash(result);
  const previousHash =
    result.previous_hash ??
    "0000000000000000000000000000000000000000000000000000000000000000";
  const chainPosition = result.chain_position ?? 0;
  const ts = result.timestamp ?? new Date().toISOString();

  return (
    <Card className="space-y-4">
      <SectionLabel>Cryptographic Evidence Chain</SectionLabel>

      <div className="grid gap-2 sm:grid-cols-2">
        <div className="rounded-lg border border-white/[0.08] bg-s3 p-3">
          <p className="font-mono text-[8px] uppercase tracking-[0.2em] text-t3">SHA-256 Hash</p>
          <p className="mt-1 break-all font-mono text-xs text-t1">{hash}</p>
        </div>
        <div className="rounded-lg border border-white/[0.08] bg-s3 p-3">
          <p className="font-mono text-[8px] uppercase tracking-[0.2em] text-t3">Timestamp</p>
          <p className="mt-1 text-xs text-t1">{new Date(ts).toLocaleString()}</p>
        </div>
        <div className="rounded-lg border border-white/[0.08] bg-s3 p-3">
          <p className="font-mono text-[8px] uppercase tracking-[0.2em] text-t3">Chain Position</p>
          <p className="mt-1 text-xs text-t1">#{chainPosition}</p>
        </div>
        <div className="rounded-lg border border-white/[0.08] bg-s3 p-3">
          <p className="font-mono text-[8px] uppercase tracking-[0.2em] text-t3">Previous Hash</p>
          <p className="mt-1 break-all font-mono text-xs text-t1">{previousHash}</p>
        </div>
      </div>

      <div className="rounded-lg border border-white/[0.08] bg-s3 p-3">
        <p className="font-mono text-[8px] uppercase tracking-[0.2em] text-t3">Chain Formula</p>
        <p className="mt-1 break-all font-mono text-xs text-t2">
          H(d) = SHA256( claim ∥ verdict ∥ timestamp ∥ H(d−1) )
        </p>
      </div>

      <div className="grid gap-2 sm:grid-cols-3">
        <GoldButton onClick={onDownloadEvidence} className="w-full justify-center">
          <Download className="h-3.5 w-3.5" />
          ⬇ DOWNLOAD EVIDENCE PDF
        </GoldButton>
        <button
          type="button"
          onClick={onCopyHash}
          className="inline-flex w-full items-center justify-center gap-2 rounded-[8px] border border-white/[0.12] bg-s3 px-3 py-2 font-mono text-[10px] uppercase tracking-[0.16em] text-t1 transition hover:border-white/[0.25]"
        >
          <Copy className="h-3.5 w-3.5" />
          COPY HASH
        </button>
        <button
          type="button"
          onClick={onVerifyIntegrity}
          className="inline-flex w-full items-center justify-center gap-2 rounded-[8px] border border-grn/35 bg-grn/15 px-3 py-2 font-mono text-[10px] uppercase tracking-[0.16em] text-grn transition hover:brightness-110"
        >
          <ShieldCheck className="h-3.5 w-3.5" />
          VERIFY INTEGRITY
        </button>
      </div>
    </Card>
  );
}
