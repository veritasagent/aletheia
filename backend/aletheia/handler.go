package aletheia

import (
	"fmt"
	"strings"
)

type Handler struct {
	chain    *DetectionChain
	demoMode bool
	server   *AletheiaServer
}

func NewHandler(chain *DetectionChain, demo bool) *Handler {
	return &Handler{
		chain:    chain,
		demoMode: demo,
	}
}

func (h *Handler) SetServer(s *AletheiaServer) {
	h.server = s
}

func (h *Handler) Handle(event MisinformationEvent) {
	event.RiskScore = calculateRiskScore(event.Confidence, event.Virality, event.Severity)
	event.RiskLevel = riskLevel(event.RiskScore)
	event.RiskLabel = event.RiskLevel

	record := h.chain.Append(event)

	switch event.Verdict {
	case "FALSE", "MISLEADING":
		fmt.Printf("[FLAGGED] %s | %s | risk=%.2f (%s)\n", event.ID, event.Verdict, event.RiskScore, event.RiskLevel)
	case "TRUE":
		fmt.Printf("[COMPLIANT] %s | %s\n", event.ID, event.Verdict)
	case "UNVERIFIED":
		fmt.Printf("[REVIEW] %s | UNVERIFIED\n", event.ID)
	default:
		fmt.Printf("[UNKNOWN] %s | %s\n", event.ID, event.Verdict)
	}

	if h.demoMode {
		h.printEvidenceBox(record)
	}

	if h.server != nil {
		h.server.BroadcastDetection(record)
	}
}

func (h *Handler) printEvidenceBox(record DetectionRecord) {
	hashPrefix := record.Hash
	if len(hashPrefix) > 16 {
		hashPrefix = hashPrefix[:16]
	}

	sourcesStr := strings.Join(record.SourcesChecked, ", ")
	if len(sourcesStr) > 70 {
		sourcesStr = sourcesStr[:67] + "..."
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf("ALETHEIA EVIDENCE | ID: %s | Segment: %s\n", record.ID, record.Segment)
	fmt.Printf("Claim: %s\n", record.Claim)
	fmt.Printf("Verdict: %s | Confidence: %d%% | Severity: %s\n", record.Verdict, record.Confidence, record.Severity)
	level := record.RiskLevel
	if level == "" {
		level = record.RiskLabel
	}
	fmt.Printf("Risk: %.2f (%s) | Platform: %s | Virality: %d\n", record.RiskScore, level, record.Platform, record.Virality)
	fmt.Printf("Sources: %s\n", sourcesStr)
	fmt.Printf("Reasoning: %s\n", record.ReasoningSummary)
	fmt.Printf("Hash: sha256:%s\n", hashPrefix)
	fmt.Println("--------------------------------------------------------------------------")
}
