"use client";

import { useState } from "react";
import { Loader2, SearchCheck } from "lucide-react";
import { verifyClaimAPI } from "@/lib/api";
import type { VerifyResponse } from "@/lib/types";

interface HashVerifierProps {
  onResult?: (result: VerifyResponse) => void;
}

export default function HashVerifier({ onResult }: HashVerifierProps) {
  const [claim, setClaim] = useState("");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState(false);

  async function handleVerify() {
    const trimmed = claim.trim();
    if (!trimmed) {
      setError("Enter a claim to verify.");
      return;
    }

    setLoading(true);
    setError("");
    try {
      const result = await verifyClaimAPI(trimmed);
      onResult?.(result);
      setClaim("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Verification failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="rounded-2xl border border-s3 bg-s1/80 p-4">
      <h2 className="font-heading text-lg font-bold text-t1">Hash Verifier</h2>
      <div className="mt-3 rounded-xl border border-s3 bg-s2/70 p-3">
        <label htmlFor="claim-input" className="mb-2 block text-xs uppercase tracking-wide text-t2">
          Claim Input
        </label>
        <textarea
          id="claim-input"
          value={claim}
          onChange={(event) => setClaim(event.target.value)}
          rows={4}
          className="w-full resize-none rounded-lg border border-s3 bg-s1 px-3 py-2 text-sm text-t1 outline-none transition focus:border-blu"
        />
        <button
          type="button"
          onClick={handleVerify}
          disabled={loading}
          className="mt-3 inline-flex items-center gap-2 rounded-lg border border-blu/40 bg-blu/20 px-4 py-2 text-sm font-semibold text-t1 transition hover:bg-blu/30 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {loading ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              Verifying...
            </>
          ) : (
            <>
              <SearchCheck className="h-4 w-4" />
              Verify Claim
            </>
          )}
        </button>
      </div>
      {error ? <p className="mt-3 text-sm text-red">{error}</p> : null}
    </section>
  );
}
