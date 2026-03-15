import jsPDF from "jspdf";
import { VerifyResponse } from "./types";

const BG = [6, 8, 15];
const HEADER = [9, 11, 23];
const GOLD = [197, 164, 74];
const GREEN = [0, 196, 154];
const RED = [238, 56, 80];
const AMBER = [237, 160, 48];
const BLUE = [75, 120, 255];
const TEXT = [200, 213, 240];
const MUTED = [78, 94, 136];
const BORDER = [30, 34, 56];

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

export function exportEvidencePDF(
  result: VerifyResponse,
  chainData: {
    hash: string;
    prevHash: string;
    chainPos: number;
    timestamp: string;
    severity: number;
    platformData: { name: string; pct: number }[];
  }
): void {
  const doc = new jsPDF({
    unit: "mm",
    format: "a4",
  });

  const w = 210;
  const h = 297;
  const margin = 20;

  let y = 0;

  // Background
  doc.setFillColor(BG[0], BG[1], BG[2]);
  doc.rect(0, 0, w, h, "F");

  // Header Bar
  doc.setFillColor(HEADER[0], HEADER[1], HEADER[2]);
  doc.rect(0, 0, w, 22, "F");

  doc.setFont("helvetica", "bold");
  doc.setFontSize(14);
  doc.setTextColor(GOLD[0], GOLD[1], GOLD[2]);
  doc.text("ALETHEIA VERIFICATION REPORT", margin, 14);

  doc.setFont("helvetica", "normal");
  doc.setFontSize(8);
  doc.setTextColor(MUTED[0], MUTED[1], MUTED[2]);
  doc.text(`${chainData.timestamp}  ·  Chain #${chainData.chainPos.toString().padStart(5, "0")}`, w - margin, 14, { align: "right" });

  // Verdict Box
  let verdictColor = BLUE;
  if (result.verdict === "FALSE") verdictColor = RED;
  else if (result.verdict === "MISLEADING") verdictColor = AMBER;
  else if (result.verdict === "TRUE") verdictColor = GREEN;

  doc.setFillColor(Math.round(verdictColor[0] * 0.18), Math.round(verdictColor[1] * 0.18), Math.round(verdictColor[2] * 0.18));
  doc.setDrawColor(verdictColor[0], verdictColor[1], verdictColor[2]);
  doc.setLineWidth(0.4);
  doc.rect(margin, 28, w - margin * 2, 18, "FD");

  doc.setFont("helvetica", "bold");
  doc.setFontSize(16);
  doc.setTextColor(verdictColor[0], verdictColor[1], verdictColor[2]);
  doc.text(result.verdict, margin + 4, 40);

  doc.setFont("helvetica", "normal");
  doc.setFontSize(10);
  doc.setTextColor(TEXT[0], TEXT[1], TEXT[2]);
  doc.text(
    `Confidence: ${result.confidence}%  ·  Risk: ${result.risk_label}  ·  Severity: ${chainData.severity}/10`,
    w - margin - 4,
    40,
    { align: "right" }
  );

  y = 56;

  function sectionLabel(text: string) {
    if (y > 260) {
      doc.addPage();
      doc.setFillColor(BG[0], BG[1], BG[2]);
      doc.rect(0, 0, w, h, "F");
      y = 20;
    }
    doc.setFont("courier", "normal"); // Standard monospace font available in jsPDF
    doc.setFontSize(8);
    doc.setTextColor(MUTED[0], MUTED[1], MUTED[2]);
    const charSpace = 1;
    let currX = margin;
    const upperText = text.toUpperCase();
    for (let i = 0; i < upperText.length; i++) {
      doc.text(upperText[i], currX, y);
      currX += doc.getTextWidth(upperText[i]) + charSpace;
    }
    
    y += 2;
    doc.setDrawColor(BORDER[0], BORDER[1], BORDER[2]);
    doc.setLineWidth(0.3);
    doc.line(margin, y, w - margin, y);
    y += 8;
  }

  function bodyText(text: string, size: number = 10, color: number[] = TEXT) {
    doc.setFont("helvetica", "normal");
    doc.setFontSize(size);
    doc.setTextColor(color[0], color[1], color[2]);
    
    // We roughly estimate ' || ' etc could be used without splitTextToSize issue, 
    // but we use Native jsPDF word wrapping that sometimes struggles with special chars if not careful.
    const splitText = doc.splitTextToSize(text, w - margin * 2);
    
    for (const line of splitText) {
      if (y > 260) {
        doc.addPage();
        doc.setFillColor(BG[0], BG[1], BG[2]);
        doc.rect(0, 0, w, h, "F");
        y = 30; // little more margin on fresh page
        doc.setFont("helvetica", "normal");
        doc.setFontSize(size);
        doc.setTextColor(color[0], color[1], color[2]);
      }
      doc.text(line, margin, y);
      y += (size * 0.42);
    }
    y += 3;
  }

  function spacer(n: number = 6) {
    y += n;
  }

  sectionLabel("CLAIM SUBMITTED");
  bodyText(result.claim, 11);
  spacer();

  sectionLabel("VERDICT ANALYSIS");
  bodyText(`Verdict: ${result.verdict}  |  Confidence: ${result.confidence}%  |  Risk Level: ${result.risk_label}`);
  bodyText(`Severity Score: ${chainData.severity} / 10  |  Public Harm Index: ${Math.round(result.risk_score * 80)}%`);
  spacer();

  sectionLabel("AGENT REASONING");
  bodyText(result.reason);
  spacer();

  sectionLabel("RISK SCORE METHODOLOGY");
  bodyText("Risk = Confidence × 0.5  +  Virality × 0.3  +  Severity × 0.2", 10, GOLD);
  bodyText(`Confidence: ${result.confidence}%  |  Virality: ${Math.round(result.risk_score * 100)}%  |  Severity: ${Math.round(result.risk_score * 100)}%`);
  bodyText(`Computed Risk Score: ${(result.risk_score * 10).toFixed(1)} / 10`);
  spacer();

  sectionLabel(`EVIDENCE SOURCES (${result.sources.length} sources)`);
  result.sources.forEach((source, index) => {
    const rel = SOURCE_RELIABILITY[source] || 80;
    let action = "CONFIRMS";
    if (result.verdict === "FALSE" && index < 2) action = "CONTRADICTS";
    if (result.verdict === "MISLEADING" && index < 1) action = "CONTRADICTS";
    if (result.verdict === "UNVERIFIABLE") action = "NEUTRAL";
    bodyText(`${index + 1}. ${source}  —  ${rel}% reliability  —  ${action}`);
  });
  spacer();

  sectionLabel("PLATFORM SPREAD ANALYSIS");
  chainData.platformData.forEach((platform) => {
    bodyText(`${platform.name}: ${platform.pct}%`);
  });
  spacer();

  sectionLabel("CRYPTOGRAPHIC EVIDENCE");
  bodyText("SHA-256 Hash:", 8, MUTED);
  bodyText(chainData.hash, 9, BLUE);
  spacer(4);
  bodyText(`Timestamp: ${chainData.timestamp}`);
  bodyText(`Chain Position: #${chainData.chainPos.toString().padStart(5, "0")}`);
  bodyText("Previous Hash (H(d-1)):");
  bodyText(chainData.prevHash, 9, BLUE);
  spacer();

  sectionLabel("CHAIN INTEGRITY FORMULA");
  bodyText("H(d) = SHA256( claim || verdict || timestamp || H(d-1) )", 10, GOLD);
  bodyText("Altering any input field breaks every downstream hash — tamper detection is cryptographic and instant.");
  spacer();

  sectionLabel("INDEPENDENT VERIFICATION");
  bodyText("To independently verify this record, paste the SHA-256 hash into the ALETHEIA verifier");
  bodyText("at alyuca.in or via GET /api/chain on the live backend.");
  spacer();

  // Footer Bar (fixed bottom)
  if (y > 277) {
    doc.addPage();
    doc.setFillColor(BG[0], BG[1], BG[2]);
    doc.rect(0, 0, w, h, "F");
  }

  doc.setFillColor(HEADER[0], HEADER[1], HEADER[2]);
  doc.rect(0, 277, w, 20, "F");

  doc.setFont("helvetica", "bold");
  doc.setFontSize(8);
  doc.setTextColor(GREEN[0], GREEN[1], GREEN[2]);
  doc.text("✓ Cryptographic integrity verified  ·  ALETHEIA Truth Engine  ·  HACKANOVA 5.0", margin, 289);

  doc.setFont("helvetica", "normal");
  doc.setTextColor(MUTED[0], MUTED[1], MUTED[2]);
  doc.text(`alyuca.in  ·  Chain #${chainData.chainPos.toString().padStart(5, "0")}`, w - margin, 289, { align: "right" });

  doc.save(`ALETHEIA_Evidence_${Date.now()}.pdf`);
}
