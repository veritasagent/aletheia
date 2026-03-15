package aletheia

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type DetectionRecord struct {
	ID               string         `json:"id"`
	Claim            string         `json:"claim_text"`
	Verdict          string         `json:"verdict"`
	Confidence       int            `json:"confidence"`
	RiskScore        float64        `json:"risk_score"`
	RiskLevel        string         `json:"risk_level"`
	RiskLabel        string         `json:"risk_label,omitempty"`
	Severity         string         `json:"severity"`
	Platform         string         `json:"platform"`
	Virality         int            `json:"virality"`
	PlatformSpread   map[string]int `json:"platform_spread,omitempty"`
	Segment          string         `json:"segment,omitempty"`
	SourcesChecked   []string       `json:"sources"`
	Reason           string         `json:"reason,omitempty"`
	ReasoningSteps   []string       `json:"reasoning_steps"`
	ReasoningChain   []string       `json:"reasoning_chain,omitempty"`
	ReasoningSummary string         `json:"reasoning_summary,omitempty"`
	RelatedEventID   string         `json:"related_event_id,omitempty"`
	RelatedHeadline  string         `json:"related_headline,omitempty"`
	CounterNarrative string         `json:"counter_narrative,omitempty"`
	Timestamp        time.Time      `json:"timestamp"`
	Hash             string         `json:"hash,omitempty"`
	EvidenceHash     string         `json:"evidence_hash"`
	PreviousHash     string         `json:"previous_hash"`
	CurrentHash      string         `json:"current_hash"`
	ChainPosition    int            `json:"chain_position"`
}

type DetectionChain struct {
	records []DetectionRecord
	mutex   sync.RWMutex
}

func NewDetectionChain() *DetectionChain {
	return &DetectionChain{
		records: make([]DetectionRecord, 0),
	}
}

func (c *DetectionChain) Append(event MisinformationEvent) DetectionRecord {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	timestamp := event.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}

	var previousHash string
	if len(c.records) > 0 {
		previousHash = recordHashValue(c.records[len(c.records)-1])
	}
	if previousHash == "" {
		previousHash = computeChainHash(event.Claim, event.Verdict, timestamp, "")
	}

	hashString := computeChainHash(event.Claim, event.Verdict, timestamp, previousHash)
	if hashString == "" {
		hashString = computeChainHash(event.Claim, event.Verdict, timestamp, "")
	}
	evidenceHash := hashString
	currentHash := hashString
	if evidenceHash == "" {
		evidenceHash = hashString
	}
	if currentHash == "" {
		currentHash = hashString
	}
	if previousHash == "" {
		previousHash = hashString
	}

	record := DetectionRecord{
		ID:               event.ID,
		Claim:            event.Claim,
		Verdict:          event.Verdict,
		Confidence:       event.Confidence,
		RiskScore:        event.RiskScore,
		RiskLevel:        event.RiskLevel,
		RiskLabel:        event.RiskLabel,
		Severity:         event.Severity,
		Platform:         event.Platform,
		Virality:         event.Virality,
		PlatformSpread:   cloneIntMap(event.PlatformSpread),
		Segment:          event.Segment,
		SourcesChecked:   cloneStrings(event.SourcesChecked),
		Reason:           event.Reason,
		ReasoningSteps:   cloneStrings(event.ReasoningSteps),
		ReasoningChain:   cloneStrings(event.ReasoningChain),
		ReasoningSummary: event.ReasoningSummary,
		RelatedEventID:   event.RelatedEventID,
		RelatedHeadline:  event.RelatedHeadline,
		CounterNarrative: event.CounterNarrative,
		Timestamp:        timestamp,
		Hash:             hashString,
		EvidenceHash:     evidenceHash,
		PreviousHash:     previousHash,
		CurrentHash:      currentHash,
		ChainPosition:    len(c.records) + 1,
	}

	c.records = append(c.records, record)
	return record
}

func (c *DetectionChain) Verify() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for i, record := range c.records {
		var previousHash string
		if i > 0 {
			previousHash = recordHashValue(c.records[i-1])
		}
		if previousHash == "" {
			previousHash = computeChainHash(record.Claim, record.Verdict, record.Timestamp, "")
		}
		if record.PreviousHash != previousHash {
			return false
		}
		expectedHash := computeChainHash(record.Claim, record.Verdict, record.Timestamp, record.PreviousHash)
		if recordHashValue(record) != expectedHash {
			return false
		}
	}
	return true
}

func (c *DetectionChain) Export() []DetectionRecord {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	exported := make([]DetectionRecord, len(c.records))
	for i := range c.records {
		exported[i] = c.records[i]
		exported[i].PlatformSpread = cloneIntMap(c.records[i].PlatformSpread)
		exported[i].SourcesChecked = cloneStrings(c.records[i].SourcesChecked)
		exported[i].ReasoningSteps = cloneStrings(c.records[i].ReasoningSteps)
		exported[i].ReasoningChain = cloneStrings(c.records[i].ReasoningChain)
	}
	return exported
}

func (c *DetectionChain) FindByID(id string) (DetectionRecord, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for i := len(c.records) - 1; i >= 0; i-- {
		if c.records[i].ID == id {
			record := c.records[i]
			record.PlatformSpread = cloneIntMap(record.PlatformSpread)
			record.SourcesChecked = cloneStrings(record.SourcesChecked)
			record.ReasoningSteps = cloneStrings(record.ReasoningSteps)
			record.ReasoningChain = cloneStrings(record.ReasoningChain)
			return record, true
		}
	}
	return DetectionRecord{}, false
}

func computeChainHash(claim string, verdict string, timestamp time.Time, previousHash string) string {
	claimBytes := []byte(claim)
	verdictBytes := []byte(verdict)
	timestampBytes := []byte(timestamp.UTC().Format(time.RFC3339Nano))
	previousHashBytes := []byte(previousHash)

	hashInput := append(claimBytes, verdictBytes...)
	hashInput = append(hashInput, timestampBytes...)
	hashInput = append(hashInput, previousHashBytes...)

	hashBytes := sha256.Sum256(hashInput)
	return hex.EncodeToString(hashBytes[:])
}

func cloneStrings(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	out := make([]string, len(items))
	copy(out, items)
	return out
}

func cloneIntMap(items map[string]int) map[string]int {
	if len(items) == 0 {
		return nil
	}
	out := make(map[string]int, len(items))
	for k, v := range items {
		out[k] = v
	}
	return out
}

func recordHashValue(record DetectionRecord) string {
	if record.CurrentHash != "" {
		return record.CurrentHash
	}
	if record.EvidenceHash != "" {
		return record.EvidenceHash
	}
	return record.Hash
}
