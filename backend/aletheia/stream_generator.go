package aletheia

import (
	"math/rand"
	"sync"
	"time"
)

type ClaimStreamGenerator struct {
	watcher *Watcher
	report  *IntelligenceReport

	minInterval time.Duration
	maxInterval time.Duration

	mu      sync.Mutex
	running bool
	stopCh  chan struct{}
	wg      sync.WaitGroup

	rng    *rand.Rand
	claims []IntelligenceMatch
	index  int
}

func NewClaimStreamGenerator(
	watcher *Watcher,
	report *IntelligenceReport,
	minInterval time.Duration,
	maxInterval time.Duration,
) *ClaimStreamGenerator {
	if minInterval <= 0 {
		minInterval = 2 * time.Second
	}
	if maxInterval < minInterval {
		maxInterval = minInterval
	}

	claims := report.AllClaims()
	rng := rand.New(rand.NewSource(20260313))
	rng.Shuffle(len(claims), func(i, j int) { claims[i], claims[j] = claims[j], claims[i] })

	return &ClaimStreamGenerator{
		watcher:     watcher,
		report:      report,
		minInterval: minInterval,
		maxInterval: maxInterval,
		rng:         rng,
		claims:      claims,
		stopCh:      make(chan struct{}),
	}
}

func (g *ClaimStreamGenerator) Start() {
	g.mu.Lock()
	if g.running {
		g.mu.Unlock()
		return
	}
	g.running = true
	g.stopCh = make(chan struct{})
	g.mu.Unlock()

	g.wg.Add(1)
	go g.loop()
}

func (g *ClaimStreamGenerator) Stop() {
	g.mu.Lock()
	if !g.running {
		g.mu.Unlock()
		return
	}
	g.running = false
	close(g.stopCh)
	g.mu.Unlock()
	g.wg.Wait()
}

func (g *ClaimStreamGenerator) loop() {
	defer g.wg.Done()
	if len(g.claims) == 0 {
		return
	}

	for {
		select {
		case <-g.stopCh:
			return
		default:
		}

		match := g.claims[g.index]
		event := MisinformationEvent{
			ID:             match.Claim.ClaimID,
			Timestamp:      time.Now().UTC(),
			Claim:          match.Claim.ClaimText,
			Platform:       match.Claim.Platform,
			Virality:       match.Claim.Virality,
			PlatformSpread: g.randomPlatformSpread(),
			Segment:        match.Event.Segment,
			RelatedEventID: match.Event.EventID,
		}
		g.watcher.Emit(event)

		g.index++
		if g.index >= len(g.claims) {
			g.index = 0
			g.rng.Shuffle(len(g.claims), func(i, j int) {
				g.claims[i], g.claims[j] = g.claims[j], g.claims[i]
			})
		}

		delay := g.randomDelay()
		timer := time.NewTimer(delay)
		select {
		case <-g.stopCh:
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (g *ClaimStreamGenerator) randomDelay() time.Duration {
	if g.maxInterval <= g.minInterval {
		return g.minInterval
	}
	span := g.maxInterval - g.minInterval
	return g.minInterval + time.Duration(g.rng.Int63n(int64(span)+1))
}

func (g *ClaimStreamGenerator) randomPlatformSpread() map[string]int {
	platforms := []string{"X", "Reddit", "WhatsApp", "Telegram", "YouTube"}
	weights := make([]int, len(platforms))
	totalWeight := 0
	for i := range platforms {
		weights[i] = g.rng.Intn(90) + 10
		totalWeight += weights[i]
	}

	totalSpread := 100
	out := make(map[string]int, len(platforms))
	assigned := 0
	for i := 0; i < len(platforms)-1; i++ {
		value := (weights[i] * totalSpread) / totalWeight
		out[platforms[i]] = value
		assigned += value
	}
	out[platforms[len(platforms)-1]] = totalSpread - assigned
	return out
}
