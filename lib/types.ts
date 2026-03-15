export interface Segment {
  id: string;
  label: string;
  icon: string;
  color: string;
}

export interface TrendingClaim {
  claim: string;
  verdict: string;
  confidence: number;
  risk_score: number;
  risk_label: string;
  reason: string;
  sources: string[];
  heat: number;
  platform: string;
}

export interface VerifyResponse {
  id?: string;
  claim: string;
  verdict: string;
  confidence: number;
  risk_score: number;
  risk_label: string;
  reason: string;
  sources: string[];
  severity?: string;
  virality?: number;
  platform?: string;
  platform_spread?: Record<string, number>;
  chain_position?: number;
  timestamp?: string;
  previous_hash?: string;
  current_hash?: string;
  evidence_hash?: string;
  hash?: string;
}

export interface ChainRecord {
  hash: string;
  prev_hash: string;
  chain_pos: number;
  timestamp: string;
}
