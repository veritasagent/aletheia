package aletheia

import "time"

type MisinformationEvent struct {
	ID               string         `json:"id"`
	Timestamp        time.Time      `json:"timestamp"`
	Claim            string         `json:"claim_text"`
	Platform         string         `json:"platform"`
	Virality         int            `json:"virality"`
	PlatformSpread   map[string]int `json:"platform_spread,omitempty"`
	Segment          string         `json:"segment,omitempty"`
	Verdict          string         `json:"verdict"`
	Confidence       int            `json:"confidence"`
	Reason           string         `json:"reason,omitempty"`
	ReasoningSteps   []string       `json:"reasoning_steps,omitempty"`
	ReasoningChain   []string       `json:"reasoning_chain,omitempty"`
	ReasoningSummary string         `json:"reasoning_summary,omitempty"`
	SourcesChecked   []string       `json:"sources,omitempty"`
	Severity         string         `json:"severity"`
	RiskScore        float64        `json:"risk_score"`
	RiskLevel        string         `json:"risk_level"`
	RiskLabel        string         `json:"risk_label,omitempty"`
	Hash             string         `json:"-"`
	EvidenceHash     string         `json:"evidence_hash,omitempty"`
	PreviousHash     string         `json:"previous_hash,omitempty"`
	CurrentHash      string         `json:"current_hash,omitempty"`
	ChainPosition    int            `json:"chain_position,omitempty"`
	RelatedEventID   string         `json:"related_event_id,omitempty"`
	RelatedHeadline  string         `json:"related_headline,omitempty"`
	CounterNarrative string         `json:"counter_narrative,omitempty"`
}
