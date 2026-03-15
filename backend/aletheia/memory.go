package aletheia

import (
	"strings"
	"sync"
)

type ImmuneMemory struct {
	chain *DetectionChain
	mu    sync.RWMutex
}

func NewImmuneMemory(chain *DetectionChain) *ImmuneMemory {
	return &ImmuneMemory{chain: chain}
}

func (m *ImmuneMemory) CheckMemory(claim string) (DetectionRecord, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	records := m.chain.Export()
	incomingLower := strings.ToLower(claim)

	for _, record := range records {
		recordLower := strings.ToLower(record.Claim)

		if incomingLower == recordLower {
			return record, true
		}
		if strings.Contains(incomingLower, recordLower) && len(recordLower) > 20 {
			return record, true
		}
		if strings.Contains(recordLower, incomingLower) && len(incomingLower) > 20 {
			return record, true
		}
	}

	return DetectionRecord{}, false
}

func (m *ImmuneMemory) FormatMemoryHit(event MisinformationEvent, record DetectionRecord) MisinformationEvent {
	event.Verdict = record.Verdict
	event.Confidence = record.Confidence
	event.Reason = "IMMUNE MEMORY HIT - Previously verified claim. Original detection: " + record.Timestamp.Format("2006-01-02 15:04:05")
	event.ReasoningSummary = event.Reason
	event.ReasoningChain = []string{
		"Checking immune memory ledger",
		"Matching prior cryptographic record",
		"Reusing verified verdict",
		"Verdict determined: " + record.Verdict,
	}
	event.SourcesChecked = cloneStrings(record.SourcesChecked)
	event.Severity = record.Severity
	event.CounterNarrative = record.CounterNarrative
	event.RelatedEventID = record.RelatedEventID
	event.RelatedHeadline = record.RelatedHeadline
	event.Segment = record.Segment

	if event.Platform == "" {
		event.Platform = record.Platform
	}
	if event.Virality == 0 {
		event.Virality = record.Virality
	}

	return event
}
