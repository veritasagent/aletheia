package aletheia

import (
	"fmt"
	"math/rand"
	"time"
)

type SandboxIngester struct {
	watcher *Watcher
	running bool
}

func NewSandboxIngester(watcher *Watcher) *SandboxIngester {
	return &SandboxIngester{
		watcher: watcher,
	}
}

func (s *SandboxIngester) Start() {
	s.running = true
	go func() {
		for s.running {
			// Select a random event from the dataset
			eventIdx := rand.Intn(len(SandboxDataset))
			event := SandboxDataset[eventIdx]

			// Select a random claim from that event
			claimIdx := rand.Intn(len(event.Claims))
			claim := event.Claims[claimIdx]

			// Create a MisinformationEvent
			mEvent := MisinformationEvent{
				ID:               fmt.Sprintf("SBX-%s", fmt.Sprintf("%d", time.Now().UnixNano())[:10]),
				Timestamp:        time.Now(),
				Claim:            claim.ClaimText,
				Platform:         claim.Platform,
				Virality:         claim.Virality,
				Verdict:          claim.Verdict,
				Confidence:       claim.Confidence,
				Severity:         claim.Severity,
				Reason:           fmt.Sprintf("SANDBOX SIMULATION — Source: %s", event.Headline),
				SourcesChecked:   claim.Sources,
				CounterNarrative: claim.CounterNarrative,
			}

			// Emit to watcher
			s.watcher.Emit(mEvent)

			// Wait 2-5 seconds
			delay := time.Duration(rand.Intn(4)+2) * time.Second
			time.Sleep(delay)
		}
	}()
}

func (s *SandboxIngester) Stop() {
	s.running = false
}
