"use client";

import { useEffect } from "react";
import { AlertTriangle, RefreshCcw } from "lucide-react";
import Card from "@/components/ui/Card";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-[400px] flex-col items-center justify-center p-6">
      <Card className="max-w-md text-center">
        <div className="mb-4 flex justify-center">
          <div className="rounded-full bg-red-100 p-3 text-red-600">
            <AlertTriangle size={32} />
          </div>
        </div>
        <h2 className="font-heading text-2xl font-bold text-[#111827]">Something went wrong</h2>
        <p className="mt-2 text-sm text-[#4B5563]">
          A runtime error occurred while rendering the dashboard. This might be due to a hydration mismatch or API failure.
        </p>
        <div className="mt-6 font-mono text-[10px] text-[#9CA3AF] bg-[#FAFAFB] p-3 rounded-lg break-all">
          {error.message || "Unknown Error"}
        </div>
        <button
          onClick={() => reset()}
          className="mt-8 inline-flex items-center gap-2 rounded-xl bg-[#111827] px-6 py-3 text-xs font-bold text-white transition hover:bg-[#1f2937]"
        >
          <RefreshCcw size={14} />
          Try again
        </button>
      </Card>
    </div>
  );
}
