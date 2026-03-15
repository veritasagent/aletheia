"use client";

import { useMemo, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { X } from "lucide-react";
import SegmentSelector from "@/components/SegmentSelector";
import ClaimFeed from "@/components/ClaimFeed";
import AgentReasoning from "@/components/AgentReasoning";
import EvidencePanel from "@/components/EvidencePanel";
import LiveFeedSidebar from "@/components/LiveFeedSidebar";
import type { Segment, TrendingClaim, VerifyResponse } from "@/lib/types";

type MonitorState = "segments" | "trending" | "reasoning" | "result";

const transition = {
  initial: { opacity: 0, y: 10 },
  animate: { opacity: 1, y: 0, transition: { duration: 0.35, ease: [0.25, 1, 0.5, 1] } },
  exit: { opacity: 0, y: -10, transition: { duration: 0.2 } },
};

export default function MonitorPage() {
  const [state, setState] = useState<MonitorState>("segments");
  const [selectedSegment, setSelectedSegment] = useState<Segment | null>(null);
  const [selectedClaim, setSelectedClaim] = useState<string>("");
  const [result, setResult] = useState<VerifyResponse | null>(null);
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);

  const canRenderTrending = useMemo(
    () => state === "trending" && selectedSegment !== null,
    [state, selectedSegment],
  );

  const handleSelectSegment = (segment: Segment) => {
    setSelectedSegment(segment);
    setSelectedClaim("");
    setResult(null);
    setState("trending");
  };

  const handleSelectClaim = (claim: string, _data?: TrendingClaim) => {
    setSelectedClaim(claim);
    setState("reasoning");
    setMobileSidebarOpen(false);
  };

  const handleReasoningComplete = (verifyResult: VerifyResponse) => {
    setResult(verifyResult);
    setState("result");
  };

  const resetFlow = () => {
    setState("segments");
    setSelectedSegment(null);
    setSelectedClaim("");
    setResult(null);
  };

  const renderBreadcrumbs = () => {
    const steps = [
      { id: "segments", label: "Segments", active: true, onClick: resetFlow },
      { id: "trending", label: selectedSegment?.id.toUpperCase() || "Segment", active: state !== "segments", onClick: () => { if (state !== "segments") { setState("trending"); setResult(null); setSelectedClaim(""); } } },
      { id: "reasoning", label: "Verification", active: state === "reasoning" || state === "result", onClick: () => {} },
      { id: "result", label: "Analysis", active: state === "result", onClick: () => {} }
    ];

    return (
      <nav className="flex items-center gap-3 font-mono text-[10px] font-bold uppercase tracking-widest mb-10 text-[#9CA3AF]">
        {steps.map((step, idx) => (
          <div key={step.id} className="flex items-center gap-3">
            <button 
              onClick={step.onClick}
              disabled={!step.active || state === step.id}
              className={`transition-colors outline-none ${step.active ? (state === step.id ? "text-[#111827] cursor-default" : "text-[#D4A74F] hover:text-[#111827]") : "text-[#D1D5DB] cursor-not-allowed"}`}
            >
              {step.label}
            </button>
            {idx < steps.length - 1 && <span className="opacity-30">/</span>}
          </div>
        ))}
      </nav>
    );
  };

  return (
    <section className="mx-auto w-full max-w-[1440px] px-4 py-8 md:px-8 md:py-12">
      <div className="grid gap-12 lg:grid-cols-[1fr_320px]">
        {/* Main Workspace */}
        <div className="min-w-0">
          {renderBreadcrumbs()}

          <AnimatePresence mode="wait">
            {state === "segments" && (
              <motion.div key="segments" {...transition}>
                <header className="mb-10">
                  <h1 className="font-heading text-4xl font-extrabold tracking-tight text-[#111827]">Intelligence Dashboard</h1>
                  <p className="text-base text-[#6B7280] mt-2">Surface emerging narratives across global intelligence segments.</p>
                </header>
                <SegmentSelector
                  selectedSegmentId={selectedSegment?.id}
                  onSelect={handleSelectSegment}
                />
              </motion.div>
            )}

            {canRenderTrending && selectedSegment && (
              <motion.div key="trending" {...transition}>
                <header className="mb-10">
                  <h1 className="font-heading text-4xl font-extrabold tracking-tight text-[#111827]">
                    {selectedSegment.label} <span className="text-[#9CA3AF] opacity-50">/ Live Intercepts</span>
                  </h1>
                  <p className="text-base text-[#6B7280] mt-2">Active MISINFO signals detected within {selectedSegment.label}.</p>
                </header>
                <ClaimFeed
                  segment={selectedSegment}
                  onSelectClaim={handleSelectClaim}
                  onBack={resetFlow}
                />
              </motion.div>
            )}

            {(state === "reasoning" || state === "result") && (
              <motion.div key="analysis-workspace" {...transition} className="space-y-8">
                {/* Reasoning Header/Banner Area */}
                {state === "reasoning" && (
                   <div className="bg-white border border-[#E5E7EB] rounded-2xl shadow-sm p-8 text-center">
                     <div className="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-[#FAFAFB] text-[#D4A74F] mb-6">
                        <span className="relative flex h-3 w-3">
                          <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-[#D4A74F] opacity-75"></span>
                          <span className="relative inline-flex h-3 w-3 rounded-full bg-[#D4A74F]"></span>
                        </span>
                     </div>
                     <h2 className="font-heading text-2xl font-bold text-[#111827] mb-3">Verification Pipeline Active</h2>
                     <p className="text-[#4B5563] italic text-lg max-w-2xl mx-auto leading-relaxed">
                       "{selectedClaim}"
                     </p>
                   </div>
                )}

                {/* The analysis results are revealed here */}
                {state === "result" && result && (
                  <EvidencePanel result={result} onVerifyAnother={resetFlow} />
                )}
                
                {/* Reasoning logs are persistent below the active cards if analysis is running */}
                {(state === "reasoning" || (state === "result" && result)) && (
                   <div className="pt-4 border-t border-[#F3F4F6]">
                      <AgentReasoning claim={selectedClaim} onComplete={handleReasoningComplete} />
                   </div>
                )}
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        {/* Persistent Desktop Sidebar */}
        <aside className="hidden lg:block relative border-l border-[#F3F4F6] pl-8">
          <div className="sticky top-[100px]">
            <header className="mb-6">
               <h3 className="font-heading text-sm font-bold uppercase tracking-widest text-[#111827]">Integrity Feed</h3>
               <div className="h-1 w-8 bg-[#D4A74F] mt-2 rounded-full" />
            </header>
            <LiveFeedSidebar onPickClaim={(claim) => handleSelectClaim(claim)} />
          </div>
        </aside>
      </div>

      {/* Floating Mobile Sidebar Trigger */}
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => setMobileSidebarOpen(true)}
        className="lg:hidden fixed bottom-6 right-6 z-40 flex h-14 w-14 items-center justify-center rounded-full bg-[#D4A74F] text-black shadow-2xl"
      >
        <span className="relative flex h-3 w-3">
          <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-white opacity-40"></span>
          <span className="relative inline-flex h-3 w-3 rounded-full bg-white"></span>
        </span>
      </motion.button>

      {/* Mobile Slide-over Sidebar */}
      <AnimatePresence>
        {mobileSidebarOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setMobileSidebarOpen(false)}
              className="fixed inset-0 z-50 bg-black/10 backdrop-blur-sm lg:hidden"
            />
            <motion.div
              initial={{ x: "100%" }}
              animate={{ x: 0 }}
              exit={{ x: "100%" }}
              transition={{ type: "spring", damping: 25, stiffness: 200 }}
              className="fixed inset-y-0 right-0 z-50 w-full max-w-sm border-l border-[#E5E7EB] bg-white p-6 shadow-2xl lg:hidden overflow-y-auto"
            >
              <div className="mb-8 flex items-center justify-between">
                <div>
                  <h3 className="font-heading text-xl font-bold text-[#111827]">Live Signals</h3>
                  <p className="text-xs text-[#9CA3AF] mt-1 uppercase tracking-widest font-mono">Real-time Intercepts</p>
                </div>
                <button
                  onClick={() => setMobileSidebarOpen(false)}
                  className="rounded-full p-2 text-[#9CA3AF] hover:bg-[#F3F4F6] hover:text-[#111827] transition-colors"
                >
                  <X className="h-6 w-6" />
                </button>
              </div>
              <LiveFeedSidebar onPickClaim={(claim) => handleSelectClaim(claim)} />
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </section>
  );
}
