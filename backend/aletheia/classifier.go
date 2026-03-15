package aletheia

import (
	"fmt"
	"strings"
)

var globalMemory *ImmuneMemory
var globalSandboxVerifier *SandboxVerifier

func SetImmuneMemory(m *ImmuneMemory) {
	globalMemory = m
}

func SetSandboxVerifier(v *SandboxVerifier) {
	globalSandboxVerifier = v
}

type SandboxVerifier struct {
	report *IntelligenceReport
}

func NewSandboxVerifier(report *IntelligenceReport) *SandboxVerifier {
	return &SandboxVerifier{report: report}
}

func Classify(event MisinformationEvent) MisinformationEvent {
	if globalMemory != nil {
		if record, found := globalMemory.CheckMemory(event.Claim); found {
			return finalizeClassification(globalMemory.FormatMemoryHit(event, record))
		}
	}

	if globalSandboxVerifier != nil {
		return finalizeClassification(globalSandboxVerifier.Verify(event))
	}

	return finalizeClassification(runRuleBasedFallback(event))
}

func (v *SandboxVerifier) Verify(event MisinformationEvent) MisinformationEvent {
	if v == nil || v.report == nil {
		return runRuleBasedFallback(event)
	}

	match, ok := v.report.MatchClaim(event.ID, event.Claim)
	if !ok {
		event.Verdict = "UNVERIFIED"
		event.Confidence = 55
		event.Severity = "MEDIUM"
		event.SourcesChecked = []string{"Sandbox Intelligence Index"}
		event.ReasoningSteps = defaultReasoningSteps()
		event.ReasoningChain = []string{
			"Scanning Reuters",
			"Checking PIB",
			"Verifying via WHO",
			"Evaluating claim",
			"Verdict determined: UNVERIFIED",
		}
		event.ReasoningSummary = "No high-confidence event match was found in the sandbox intelligence report."
		event.Reason = event.ReasoningSummary
		event.CounterNarrative = "No verified event mapping exists for this claim yet. Keep it unverified until trusted sources confirm."
		return event
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = match.Event.Timestamp
	}
	if event.Platform == "" {
		event.Platform = match.Claim.Platform
	}
	if event.Virality == 0 {
		event.Virality = match.Claim.Virality
	}
	if event.Segment == "" {
		event.Segment = match.Event.Segment
	}

	event.RelatedEventID = match.Event.EventID
	event.RelatedHeadline = match.Event.Headline
	event.Verdict = match.Claim.Verdict
	event.Confidence = match.Claim.Confidence
	event.Severity = match.Claim.Severity
	event.CounterNarrative = match.Claim.CounterNarrative
	event.SourcesChecked = uniqueStrings(append(cloneStrings(match.Event.CredibleSources), match.Event.VerificationSource))
	event.ReasoningSteps = defaultReasoningSteps()
	event.ReasoningChain = buildReasoningChain(match.Event, event.Verdict)
	event.ReasoningSummary = fmt.Sprintf(
		"Matched to %s (%s context) and verified with %d trusted sources.",
		match.Event.EventID,
		match.Event.Context,
		len(event.SourcesChecked),
	)
	event.Reason = event.ReasoningSummary

	return event
}

func runRuleBasedFallback(event MisinformationEvent) MisinformationEvent {
	lowerClaim := strings.ToLower(event.Claim)
	event.SourcesChecked = []string{"Local Rule Engine"}

	switch {
	case strings.Contains(lowerClaim, "ban"), strings.Contains(lowerClaim, "suspend"), strings.Contains(lowerClaim, "freeze"), strings.Contains(lowerClaim, "shutdown"):
		event.Verdict = "FALSE"
		event.Confidence = 84
		event.Severity = "HIGH"
		event.CounterNarrative = "No official bulletin confirms this claim. Verify from trusted institutional sources before sharing."
	case strings.Contains(lowerClaim, "cure"), strings.Contains(lowerClaim, "vaccine"), strings.Contains(lowerClaim, "secret"), strings.Contains(lowerClaim, "viral"):
		event.Verdict = "MISLEADING"
		event.Confidence = 76
		event.Severity = "CRITICAL"
		event.CounterNarrative = "This claim is likely context-distorted. Refer to trusted medical and public-information bulletins."
	case strings.Contains(lowerClaim, "official"), strings.Contains(lowerClaim, "announced"), strings.Contains(lowerClaim, "confirmed"):
		event.Verdict = "TRUE"
		event.Confidence = 72
		event.Severity = "LOW"
		event.CounterNarrative = "Claim appears plausible, but should still be cross-checked in official channels."
	default:
		event.Verdict = "UNVERIFIED"
		event.Confidence = 58
		event.Severity = "MEDIUM"
		event.CounterNarrative = "Insufficient corroboration. Keep this claim in unverified state until trusted evidence appears."
	}

	event.ReasoningSteps = defaultReasoningSteps()
	event.ReasoningChain = []string{
		"Scanning Reuters",
		"Checking PIB",
		"Verifying via WHO",
		"Evaluating claim",
		"Verdict determined: " + event.Verdict,
	}
	event.ReasoningSummary = "Rule-based fallback applied because no sandbox intelligence match was available."
	event.Reason = event.ReasoningSummary
	return event
}

func buildReasoningChain(event IntelligenceEvent, verdict string) []string {
	steps := []string{
		"Scanning Reuters feed",
		"Scanning PIB bulletins",
		"Scanning WHO/public advisories",
	}
	for _, src := range event.CredibleSources {
		if len(steps) >= 5 {
			break
		}
		steps = append(steps, "Scanning "+src)
	}
	steps = append(steps, "Evaluating claim")
	steps = append(steps, "Verdict determined: "+verdict)
	return uniqueStringsPreserveOrder(steps)
}

func finalizeClassification(event MisinformationEvent) MisinformationEvent {
	if len(event.ReasoningSteps) == 0 {
		event.ReasoningSteps = defaultReasoningSteps()
	}
	if len(event.ReasoningChain) == 0 {
		event.ReasoningChain = append([]string{}, event.ReasoningSteps...)
		event.ReasoningChain = append(event.ReasoningChain, "Verdict determined: "+event.Verdict)
	}
	event.RiskScore = calculateRiskScore(event.Confidence, event.Virality, event.Severity)
	event.RiskLevel = riskLevel(event.RiskScore)
	event.RiskLabel = event.RiskLevel
	return event
}

func defaultReasoningSteps() []string {
	return []string{
		"Scanning Reuters",
		"Checking PIB",
		"Verifying via WHO",
		"Evaluating claim",
	}
}

func calculateRiskScore(confidence int, virality int, severity string) float64 {
	confidenceScore := float64(confidence) / 10.0
	if confidenceScore < 0 {
		confidenceScore = 0
	}
	if confidenceScore > 10 {
		confidenceScore = 10
	}

	viralityScore := float64(virality) / 100000.0
	if viralityScore < 0 {
		viralityScore = 0
	}
	if viralityScore > 10 {
		viralityScore = 10
	}

	severityScore := severityNumeric(severity)
	score := (confidenceScore * 0.5) + (viralityScore * 0.3) + (severityScore * 0.2)
	if score < 0 {
		return 0
	}
	if score > 10 {
		return 10
	}
	return score
}

func severityNumeric(severity string) float64 {
	switch strings.ToUpper(strings.TrimSpace(severity)) {
	case "LOW":
		return 2
	case "MEDIUM":
		return 5
	case "HIGH":
		return 8
	case "CRITICAL":
		return 10
	default:
		return 3
	}
}

func riskLevel(score float64) string {
	switch {
	case score <= 3:
		return "LOW"
	case score <= 6:
		return "MEDIUM"
	case score <= 8:
		return "HIGH"
	default:
		return "CRITICAL"
	}
}

func uniqueStrings(items []string) []string {
	return uniqueStringsPreserveOrder(items)
}

func uniqueStringsPreserveOrder(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
