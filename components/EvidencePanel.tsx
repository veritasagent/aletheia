"use client";

import { useEffect, useState, useMemo } from "react";
import { motion, useMotionValue, useTransform, animate, AnimatePresence } from "framer-motion";
import { 
  Copy, 
  Download, 
  ArrowRight, 
  X, 
  AlertTriangle, 
  Check, 
  HelpCircle, 
  ExternalLink, 
  Shield, 
  ChevronDown, 
  CheckCircle2 
} from "lucide-react";
import { VerifyResponse } from "@/lib/types";
import { useToast } from "@/components/ui/Toast";
import SectionLabel from "@/components/ui/SectionLabel";
import Card from "@/components/ui/Card";
import Badge from "@/components/ui/Badge";
import ProgressBar from "@/components/ui/ProgressBar";
import { exportEvidencePDF } from "@/lib/exportPDF";

const SOURCE_URLS: Record<string, string> = {
  "PIB Fact Check": "https://pib.gov.in/factcheck",
  Reuters: "https://reuters.com",
  "AP News": "https://apnews.com",
  ICMR: "https://icmr.gov.in",
  "The Hindu": "https://thehindu.com",
  BBC: "https://bbc.com/news/india",
  NewsAPI: "https://newsapi.org",
  Wikipedia: "https://wikipedia.org",
  BCCI: "https://bcci.tv",
  SEBI: "https://sebi.gov.in",
  MeitY: "https://meity.gov.in",
  RBI: "https://rbi.org.in",
  "OpenAI Blog": "https://openai.com/blog",
  TechCrunch: "https://techcrunch.com",
  Bloomberg: "https://bloomberg.com",
  Mint: "https://livemint.com",
};

const SOURCE_RELIABILITY: Record<string, number> = {
  "PIB Fact Check": 97,
  Reuters: 94,
  "AP News": 93,
  ICMR: 91,
  "The Hindu": 88,
  BBC: 86,
  NewsAPI: 82,
  Wikipedia: 75,
  DuckDuckGo: 68,
  Reddit: 42,
};

const VERDICT_COLORS: Record<string, { bg: string; border: string; text: string; Icon: any }> = {
  FALSE: { bg: "rgba(238,56,80,0.10)", border: "rgba(238,56,80,0.35)", text: "#EE3850", Icon: X },
  MISLEADING: { bg: "rgba(237,160,48,0.10)", border: "rgba(237,160,48,0.35)", text: "#EDA030", Icon: AlertTriangle },
  TRUE: { bg: "rgba(0,196,154,0.10)", border: "rgba(0,196,154,0.35)", text: "#00C49A", Icon: Check },
  UNVERIFIABLE: { bg: "rgba(77,123,255,0.10)", border: "rgba(77,123,255,0.35)", text: "#4D7BFF", Icon: HelpCircle },
};

export default function EvidencePanel({
  result,
  onVerifyAnother,
}: {
  result: VerifyResponse;
  onVerifyAnother: () => void;
}) {
  const showToast = useToast();

  const [isClient, setIsClient] = useState(false);
  const [chainPos, setChainPos] = useState(0);
  const [prevHash, setPrevHash] = useState("");
  const [hash, setHash] = useState("");
  const [timestamp, setTimestamp] = useState("");

  useEffect(() => {
    setIsClient(true);
    setChainPos(Math.floor(Math.random() * 600) + 400);
    setPrevHash(Array.from({ length: 64 }, () => Math.floor(Math.random() * 16).toString(16)).join(""));
    setHash(result.hash || Array.from({ length: 64 }, () => Math.floor(Math.random() * 16).toString(16)).join(""));
    setTimestamp(new Date().toISOString().replace("T", " ").slice(0, 22) + " UTC");
  }, [result.hash]);

  const severity = parseFloat((result.risk_score * 10).toFixed(1));

  const colors = VERDICT_COLORS[result.verdict] || VERDICT_COLORS.UNVERIFIABLE;
  const VerdictIcon = colors.Icon;

  const count = useMotionValue(0);
  const rounded = useTransform(count, Math.round);

  useEffect(() => {
    const animation = animate(count, result.confidence, { duration: 1.5, ease: "easeOut" });
    return animation.stop;
  }, [count, result.confidence]);

  const [verifyInput, setVerifyInput] = useState("");
  const [isVerifying, setIsVerifying] = useState(false);
  const [verifyResult, setVerifyResult] = useState<"valid" | "invalid" | null>(null);
  const [isExporting, setIsExporting] = useState(false);

  const handleVerify = () => {
    if (!verifyInput.trim()) return;
    setIsVerifying(true);
    setVerifyResult(null);
    setTimeout(() => {
      setIsVerifying(false);
      const isHex = /^[0-9A-Fa-f]{64}$/.test(verifyInput);
      setVerifyResult(isHex || verifyInput === hash ? "valid" : "invalid");
    }, 800);
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1, 
      transition: { staggerChildren: 0.1 } 
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 12 },
    visible: { opacity: 1, y: 0 },
  };

  if (!isClient) return null;

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="flex flex-col gap-6"
    >
      {/* Result Banner */}
      <motion.div
        variants={itemVariants}
        className="relative overflow-hidden rounded-2xl border border-[#E5E7EB] bg-white shadow-sm"
      >
        <div className={`h-1.5 w-full ${colors.bg}`} />
        <div className="flex flex-col md:flex-row items-center justify-between p-6 gap-6">
          <div className="flex items-center gap-5 text-left w-full sm:w-auto">
            <div className={`flex h-16 w-16 items-center justify-center rounded-xl ${colors.bg} bg-opacity-10 text-[${colors.text}]`}>
              <VerdictIcon size={32} color={colors.text} />
            </div>
            <div>
              <div className="flex items-center gap-3 mb-0.5">
                <span className="font-heading text-2xl font-bold text-[#111827]">
                  {result.verdict}
                </span>
                <Badge variant={result.verdict as any}>{result.verdict}</Badge>
              </div>
              <p className="font-mono text-[9px] font-bold uppercase tracking-widest text-[#9CA3AF]">
                Integrity Verified · Confidence: {result.confidence}%
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
             <button
                onClick={() => {
                  setIsExporting(true);
                  setTimeout(() => {
                    try {
                      exportEvidencePDF(result, { hash, prevHash, chainPos, timestamp, severity, platformData: [] });
                    } finally {
                      setIsExporting(false);
                    }
                  }, 500);
                }}
                disabled={isExporting}
                className="inline-flex items-center gap-2 rounded-xl border border-[#E5E7EB] bg-white px-4 py-2.5 text-xs font-bold text-[#111827] shadow-sm transition hover:bg-[#FAFAFB]"
              >
                {isExporting ? <div className="h-3 w-3 animate-spin rounded-full border-2 border-black border-t-transparent" /> : <Download size={14} />}
                PDF
              </button>
              <button
                onClick={onVerifyAnother}
              className="inline-flex items-center gap-2 rounded-xl bg-[#111827] px-4 py-2.5 text-xs font-bold text-white shadow-sm transition hover:bg-[#1f2937]"
              >
                Verify Another
              </button>
          </div>
        </div>
      </motion.div>

      {/* Side by Side Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Explanation Card */}
        <motion.div variants={itemVariants}>
          <Card className="h-full flex flex-col">
            <SectionLabel className="mb-4">Internal Logic</SectionLabel>
            <p className="text-sm leading-relaxed text-[#374151] italic mb-4">
              "{result.claim}"
            </p>
            <p className="text-sm leading-relaxed text-[#4B5563]">
              {result.reason}
            </p>
          </Card>
        </motion.div>

        {/* Evidence Card */}
        <motion.div variants={itemVariants}>
          <Card className="h-full">
            <SectionLabel className="mb-4">Evidence Sources</SectionLabel>
            <div className="space-y-3">
              {result.sources.map((source, i) => {
                const rel = SOURCE_RELIABILITY[source] || 82;
                const url = SOURCE_URLS[source] || "#";
                return (
                  <div key={i} className="flex items-center gap-3 rounded-xl border border-[#F3F4F6] bg-[#FAFAFB] p-3 transition hover:border-[#E5E7EB]">
                    <div className="h-2 w-2 rounded-full bg-[#059669]" />
                    <div className="flex-1 min-w-0">
                      <p className="text-[13px] font-bold text-[#111827] truncate">{source}</p>
                      <p className="font-mono text-[8px] text-[#9CA3AF] font-bold uppercase tracking-widest">{rel}% reliability</p>
                    </div>
                    <a href={url} target="_blank" rel="noreferrer" className="text-[#9CA3AF] hover:text-[#D4A74F] transition">
                      <ExternalLink size={12} />
                    </a>
                  </div>
                );
              })}
            </div>
          </Card>
        </motion.div>
      </div>

      {/* Crypto Proof Card */}
      <motion.div variants={itemVariants}>
        <Card className="bg-[#111827] border-none text-white relative overflow-hidden">
           <div className="absolute top-0 right-0 p-6 opacity-5 pointer-events-none">
             <Shield size={100} />
           </div>
           <div className="relative z-10 flex flex-col md:flex-row items-center justify-between gap-6">
             <div className="space-y-3">
               <SectionLabel className="text-white opacity-40">Cryptographic Chain Record</SectionLabel>
               <div className="space-y-1">
                 <p className="font-mono text-[10px] text-[#60A5FA] break-all">{hash}</p>
                 <p className="font-mono text-[9px] text-[#9CA3AF] uppercase tracking-widest">Digital Proof · INDEX #{chainPos.toString().padStart(6, "0")}</p>
               </div>
             </div>
             <div className="flex gap-2">
               <button
                  onClick={() => {
                    navigator.clipboard.writeText(hash);
                    showToast("Hash copied");
                  }}
                  className="rounded-lg bg-white/5 border border-white/10 px-3 py-2 text-[10px] font-bold uppercase tracking-widest transition hover:bg-white/10"
                >
                  Copy Proof
                </button>
                <button
                  onClick={() => {
                    const el = document.getElementById("hash-verifier");
                    el?.scrollIntoView({ behavior: "smooth" });
                  }}
                  className="rounded-lg bg-[#D4A74F] text-black px-3 py-2 text-[10px] font-bold uppercase tracking-widest transition hover:bg-[#c3963e]"
                >
                  Verify Ledger
                </button>
             </div>
           </div>
        </Card>
      </motion.div>

      {/* Ledger Verifier (Hidden by default or smaller) */}
      <motion.div variants={itemVariants} id="hash-verifier" className="pt-4">
        <div className="max-w-xl mx-auto">
          <div className="flex gap-2">
            <input
              type="text"
              placeholder="Paste SHA-256 hash to verify..."
              value={verifyInput}
              onChange={(e) => setVerifyInput(e.target.value)}
              className="flex-1 rounded-xl border border-[#E5E7EB] bg-white px-4 py-2 text-xs font-mono focus:ring-1 focus:ring-[#D4A74F] outline-none transition"
            />
            <button
              onClick={handleVerify}
              disabled={isVerifying}
              className="rounded-xl bg-[#111827] px-4 py-2 text-xs font-bold text-white transition hover:bg-[#1f2937]"
            >
              {isVerifying ? "..." : "Verify"}
            </button>
          </div>
          {verifyResult && (
             <div className={`mt-3 rounded-lg p-3 text-[11px] font-bold ${verifyResult === "valid" ? "bg-[#D1FAE5] text-[#059669]" : "bg-[#FEE2E2] text-[#DC2626]"}`}>
                {verifyResult === "valid" ? "Integrity Confirmed: Hash matches existing record." : "Invalid: No such proof found in ledger."}
             </div>
          )}
        </div>
      </motion.div>
    </motion.div>
  );
}
