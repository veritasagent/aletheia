import type { ChainRecord, Segment, TrendingClaim } from "@/lib/types";

export type TrendingKey = "tech" | "finance" | "politics" | "sports";

export const fallbackSegments: Segment[] = [
  { id: "tech", label: "Technology", icon: "Cpu", color: "blu" },
  { id: "finance", label: "Finance", icon: "Landmark", color: "gold" },
  { id: "politics", label: "Politics", icon: "Vote", color: "red" },
  { id: "sports", label: "Sports", icon: "Trophy", color: "grn" },
];

export const TRENDING_DATA: Record<TrendingKey, TrendingClaim[]> = {
  tech: [
    {
      claim: "OpenAI GPT-5 secretly deployed to 10M Indian users",
      verdict: "MISLEADING",
      confidence: 82,
      risk_score: 0.78,
      risk_label: "HIGH",
      heat: 92,
      platform: "LinkedIn",
      reason:
        "OpenAI has not announced any India-specific GPT-5 deployment. Claim misrepresents an infrastructure update blog post.",
      sources: ["Reuters", "OpenAI Blog", "TechCrunch"],
    },
    {
      claim: "Reliance Jio shutting down 5G due to spectrum dispute with DoT",
      verdict: "FALSE",
      confidence: 91,
      risk_score: 0.88,
      risk_label: "CRITICAL",
      heat: 87,
      platform: "WhatsApp",
      reason:
        "DoT confirmed no dispute. Jio 5G operational in 406 cities. PIB Fact Check issued formal debunking.",
      sources: ["PIB Fact Check", "Reuters", "Economic Times"],
    },
    {
      claim: "India AI regulation bill passed banning generative AI",
      verdict: "FALSE",
      confidence: 88,
      risk_score: 0.82,
      risk_label: "HIGH",
      heat: 74,
      platform: "Twitter/X",
      reason:
        "No AI regulation bill has been tabled in Lok Sabha. MeitY discussion paper only, no legislation introduced.",
      sources: ["MeitY Official", "The Hindu", "Mint"],
    },
    {
      claim: "Apple to manufacture 40% of iPhones in India by end of 2026",
      verdict: "MISLEADING",
      confidence: 70,
      risk_score: 0.45,
      risk_label: "MEDIUM",
      heat: 55,
      platform: "Bloomberg",
      reason:
        "Apple announced India expansion but 40% figure unconfirmed. Analyst range is 15-25%.",
      sources: ["Bloomberg", "Financial Times", "Reuters"],
    },
    {
      claim: "ChatGPT launched in Hindi with real-time voice support in India",
      verdict: "TRUE",
      confidence: 94,
      risk_score: 0.12,
      risk_label: "LOW",
      heat: 41,
      platform: "Product Hunt",
      reason:
        "OpenAI confirmed Hindi Voice Mode in December 2025, available on iOS and Android India.",
      sources: ["OpenAI Official", "NDTV", "Economic Times"],
    },
    {
      claim: "Samsung building ₹40,000 crore semiconductor fab in Bengaluru",
      verdict: "UNVERIFIABLE",
      confidence: 58,
      risk_score: 0.52,
      risk_label: "MEDIUM",
      heat: 63,
      platform: "Business Standard",
      reason:
        "Samsung held talks with Karnataka govt but no MOU signed. Based on anonymous sources only.",
      sources: ["Business Standard", "Reuters", "Karnataka Govt"],
    },
  ],
  finance: [
    {
      claim: "RBI has frozen all UPI transactions above ₹2,000 effective tonight",
      verdict: "FALSE",
      confidence: 93,
      risk_score: 0.89,
      risk_label: "CRITICAL",
      heat: 95,
      platform: "WhatsApp",
      reason:
        "RBI issued no such order. NPCI confirmed normal UPI operations; message is fabricated and reused from old scam forwards.",
      sources: ["RBI Bulletin", "NPCI", "PIB Fact Check"],
    },
    {
      claim: "SBI declared bankruptcy and branches will remain closed for 30 days",
      verdict: "FALSE",
      confidence: 96,
      risk_score: 0.9,
      risk_label: "CRITICAL",
      heat: 89,
      platform: "Telegram",
      reason:
        "SBI quarterly filing shows positive capital ratios and normal operations. Viral screenshot is edited from a satire page.",
      sources: ["SBI Investor Relations", "SEBI Filings", "Reuters"],
    },
    {
      claim: "Government approved a one-time 35% tax on all savings account balances",
      verdict: "MISLEADING",
      confidence: 81,
      risk_score: 0.73,
      risk_label: "HIGH",
      heat: 78,
      platform: "Twitter/X",
      reason:
        "An old policy discussion note was misquoted as enacted law. No Finance Ministry notification exists.",
      sources: ["Ministry of Finance", "The Hindu BusinessLine", "PIB Fact Check"],
    },
    {
      claim: "India to launch retail CBDC cashbacks nationwide from April 2026",
      verdict: "UNVERIFIABLE",
      confidence: 61,
      risk_score: 0.48,
      risk_label: "MEDIUM",
      heat: 57,
      platform: "LinkedIn",
      reason:
        "RBI has pilot updates but no final nationwide cashback circular has been published.",
      sources: ["RBI Digital Rupee Updates", "Livemint", "Reuters"],
    },
    {
      claim: "NSE to shift equity trading settlement entirely to T+0 by July 2026",
      verdict: "MISLEADING",
      confidence: 74,
      risk_score: 0.44,
      risk_label: "MEDIUM",
      heat: 49,
      platform: "Economic Times",
      reason:
        "SEBI approved limited T+0 optional framework; full market-wide transition date is not announced.",
      sources: ["SEBI Circulars", "NSE Announcements", "Economic Times"],
    },
    {
      claim: "Income tax return filing deadline extended to August 31, 2026",
      verdict: "TRUE",
      confidence: 88,
      risk_score: 0.18,
      risk_label: "LOW",
      heat: 35,
      platform: "Press Release",
      reason:
        "CBDT issued official extension notice on its portal and social channels.",
      sources: ["CBDT", "Income Tax Portal", "PIB"],
    },
  ],
  politics: [
    {
      claim: "Election Commission cancelled all bypolls citing EVM vulnerabilities",
      verdict: "FALSE",
      confidence: 92,
      risk_score: 0.86,
      risk_label: "CRITICAL",
      heat: 91,
      platform: "Twitter/X",
      reason:
        "ECI election calendar remains active. Viral memo is forged and does not match ECI document format.",
      sources: ["Election Commission of India", "PIB Fact Check", "Reuters"],
    },
    {
      claim: "Parliament passed emergency law restricting independent digital news outlets",
      verdict: "MISLEADING",
      confidence: 83,
      risk_score: 0.75,
      risk_label: "HIGH",
      heat: 84,
      platform: "YouTube",
      reason:
        "No such law has passed both houses. Clip references a standing committee discussion draft.",
      sources: ["PRS Legislative Research", "Lok Sabha Records", "The Hindu"],
    },
    {
      claim: "Cabinet approved nationwide internet shutdown protocol for protest coverage",
      verdict: "UNVERIFIABLE",
      confidence: 57,
      risk_score: 0.58,
      risk_label: "MEDIUM",
      heat: 68,
      platform: "Telegram",
      reason:
        "Policy consultation rumors exist but no gazetted notification confirms this decision.",
      sources: ["MHA Notifications", "Reuters", "Indian Express"],
    },
    {
      claim: "Supreme Court struck down all political funding disclosures",
      verdict: "FALSE",
      confidence: 90,
      risk_score: 0.79,
      risk_label: "HIGH",
      heat: 73,
      platform: "WhatsApp",
      reason:
        "Court judgment wording is being reversed in viral posts; disclosures remain enforceable in current order.",
      sources: ["Supreme Court Orders", "Bar and Bench", "Reuters"],
    },
    {
      claim: "New voter ID update deadline moved to June 30 nationwide",
      verdict: "TRUE",
      confidence: 87,
      risk_score: 0.21,
      risk_label: "LOW",
      heat: 39,
      platform: "ECI Portal",
      reason:
        "ECI published revised timeline and state officers circulated matching advisories.",
      sources: ["Election Commission of India", "State CEO Portals", "ANI"],
    },
    {
      claim: "Government to dissolve state assemblies before monsoon session",
      verdict: "UNVERIFIABLE",
      confidence: 54,
      risk_score: 0.5,
      risk_label: "MEDIUM",
      heat: 61,
      platform: "News Blog",
      reason:
        "No constitutional notification exists; claim is based on speculative editorials and unnamed sources.",
      sources: ["Gazette of India", "Press Information Bureau", "Reuters"],
    },
  ],
  sports: [
    {
      claim: "BCCI suspended IPL 2026 after anti-corruption raids in three franchises",
      verdict: "FALSE",
      confidence: 94,
      risk_score: 0.87,
      risk_label: "CRITICAL",
      heat: 93,
      platform: "WhatsApp",
      reason:
        "BCCI fixtures and broadcasts are ongoing. Viral statement carries fake BCCI letterhead and incorrect signatory.",
      sources: ["BCCI Official", "ESPNcricinfo", "Reuters"],
    },
    {
      claim: "ICC has banned India from Champions Trophy due to pitch tampering case",
      verdict: "FALSE",
      confidence: 91,
      risk_score: 0.82,
      risk_label: "HIGH",
      heat: 81,
      platform: "Twitter/X",
      reason:
        "ICC disciplinary bulletin contains no such order. Story originated from parody account reposts.",
      sources: ["ICC Media Releases", "Reuters", "Cricbuzz"],
    },
    {
      claim: "Neeraj Chopra tested positive before Asian Athletics meet",
      verdict: "MISLEADING",
      confidence: 79,
      risk_score: 0.69,
      risk_label: "HIGH",
      heat: 76,
      platform: "YouTube",
      reason:
        "A leaked document references routine sample collection, not adverse findings. Report falsely labels it as positive.",
      sources: ["World Athletics", "SAI", "The Indian Express"],
    },
    {
      claim: "Indian football federation dissolved after FIFA governance warning",
      verdict: "UNVERIFIABLE",
      confidence: 59,
      risk_score: 0.53,
      risk_label: "MEDIUM",
      heat: 62,
      platform: "Reddit",
      reason:
        "FIFA issued compliance observations, but no dissolution order has been published by AIFF.",
      sources: ["FIFA", "AIFF", "Sportstar"],
    },
    {
      claim: "IOA confirms India will bid for 2036 Olympics this summer",
      verdict: "TRUE",
      confidence: 85,
      risk_score: 0.2,
      risk_label: "LOW",
      heat: 43,
      platform: "Press Conference",
      reason:
        "IOA officials publicly confirmed bid-preparation timeline and intent documents.",
      sources: ["IOA", "Reuters", "Hindustan Times"],
    },
    {
      claim: "WPL salary cap removed immediately; franchises can sign unlimited overseas stars",
      verdict: "MISLEADING",
      confidence: 72,
      risk_score: 0.41,
      risk_label: "MEDIUM",
      heat: 51,
      platform: "Sports Blog",
      reason:
        "WPL committee discussed cap revisions for next cycle; no immediate removal has been ratified.",
      sources: ["BCCI WPL Regulations", "ESPNcricinfo", "Cricbuzz"],
    },
  ],
};

export const fallbackClaims: TrendingClaim[] = [
  ...TRENDING_DATA.tech.slice(0, 3),
  ...TRENDING_DATA.finance.slice(0, 3),
];

export const fallbackChain: ChainRecord[] = [
  {
    hash: "ab9a1e89f8a0622d5f6d6e7a9b0f67a9a2c3a3123411a8f8d6e2b1f50ab44319",
    prev_hash: "41f92d38b85e0a1d0fd1df6e4a0b11d8cc2e01e92f9a44d2d3bf9f64b7d1302a",
    chain_pos: 173,
    timestamp: "2026-03-13T14:29:55Z",
  },
  {
    hash: "d9f8f1e1a44a3130cd5740f91a7c73626b0207ea8e9b2f8b03e7d6d4f5d3af91",
    prev_hash: "ab9a1e89f8a0622d5f6d6e7a9b0f67a9a2c3a3123411a8f8d6e2b1f50ab44319",
    chain_pos: 174,
    timestamp: "2026-03-13T14:30:10Z",
  },
];
