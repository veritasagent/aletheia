package aletheia

import "time"

type IntelligenceEvent struct {
	EventID    string   `json:"event_id"`
	Headline   string   `json:"headline"`
	Segment    string   `json:"segment"`
	Topics     []string `json:"topics"`
	Context    string   `json:"context"`
	Date       string   `json:"date"`
	TrendScore int      `json:"trend_score"`
	Sources    []string `json:"sources"`
	Claims     []Claim  `json:"claims"`

	Summary            string    `json:"summary"`
	Timestamp          time.Time `json:"timestamp"`
	CredibleSources    []string  `json:"credible_sources"`
	VerificationSource string    `json:"verification_source"`
	Explanation        string    `json:"explanation"`
}

type Claim struct {
	ClaimID string `json:"claim_id"`

	ClaimText    string    `json:"claim_text"`
	Platform     string    `json:"platform"`
	PlatformIcon string    `json:"platform_icon"`
	Virality     int       `json:"virality"`
	TrendScore   int       `json:"trend_score"`
	Timestamp    time.Time `json:"timestamp"`

	SourceType     string `json:"source_type"`
	SourceCategory string `json:"source_category"`

	Verdict    string  `json:"verdict"`
	Confidence int     `json:"confidence"`
	Severity   string  `json:"severity"`
	RiskScore  float64 `json:"risk_score"`

	VerifiedBy       string `json:"verified_by"`
	Explanation      string `json:"explanation"`
	CounterNarrative string `json:"counter_narrative"`
}

type IntelligenceClaim = Claim

func LoadSandboxEvents() []IntelligenceEvent {
	report := NewSandboxIntelligenceReport()
	events := make([]IntelligenceEvent, len(report.Events))
	copy(events, report.Events)
	return events
}
