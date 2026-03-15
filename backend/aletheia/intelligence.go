package aletheia

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unicode"
)

type IntelligenceMatch struct {
	Event IntelligenceEvent
	Claim IntelligenceClaim
}

type intelligenceRef struct {
	EventIndex int
	ClaimIndex int
}

type IntelligenceReport struct {
	Segments    []string
	Events      []IntelligenceEvent
	byClaimID   map[string]intelligenceRef
	byClaimText map[string]intelligenceRef
}

type segmentBlueprint struct {
	Name         string
	Context      string
	Sources      []string
	Subjects     []string
	Developments []string
}

func NewSandboxIntelligenceReport() *IntelligenceReport {
	blueprints := []segmentBlueprint{
		{
			Name:    "Finance",
			Context: "India",
			Sources: []string{"RBI Bulletin", "Ministry of Finance", "Reuters India", "Press Information Bureau"},
			Subjects: []string{
				"UPI settlement window", "digital rupee pilot", "banking fraud helpline", "cooperative bank compliance cycle",
			},
			Developments: []string{
				"receives a regulatory clarification", "enters a monitored expansion phase", "is updated in a public circular", "starts a quarterly audit milestone", "is confirmed as operational nationwide",
			},
		},
		{
			Name:    "Technology",
			Context: "India",
			Sources: []string{"MeitY", "NASSCOM", "Reuters Technology", "Press Information Bureau"},
			Subjects: []string{
				"national AI compute cluster", "responsible AI policy draft", "public sector chatbot rollout", "AI skilling mission",
			},
			Developments: []string{
				"moves to the implementation stage", "gets a standards revision", "adds new operational guardrails", "is validated through pilot feedback", "is published in an official update",
			},
		},
		{
			Name:    "Business",
			Context: "India",
			Sources: []string{"Ministry of Commerce", "NITI Aayog", "Reuters Markets", "SEBI Bulletin"},
			Subjects: []string{
				"manufacturing incentive corridor", "SME credit guarantee", "exports performance tracker", "logistics cost reduction plan",
			},
			Developments: []string{
				"shows steady policy execution", "records an official progress update", "is announced for multi-state rollout", "enters the compliance review cycle", "is reaffirmed in the budget briefing",
			},
		},
		{
			Name:    "Science",
			Context: "Global",
			Sources: []string{"Nature", "ScienceDaily", "ISRO", "Reuters Science"},
			Subjects: []string{
				"lunar resource mapping mission", "fusion research milestone", "climate modeling research release", "deep-space telescope calibration",
			},
			Developments: []string{
				"publishes a verified status report", "adds new peer-reviewed findings", "issues a technical clarification", "extends research collaboration coverage", "confirms reproducibility benchmarks",
			},
		},
		{
			Name:    "Politics",
			Context: "India",
			Sources: []string{"Election Commission of India", "Press Information Bureau", "Parliament Bulletin", "Reuters India Politics"},
			Subjects: []string{
				"state election schedule notice", "parliament session agenda", "cabinet policy communication", "public grievance platform update",
			},
			Developments: []string{
				"is formally published in official channels", "is clarified through a government briefing", "receives a procedural correction", "is confirmed through parliamentary records", "is updated with verified implementation notes",
			},
		},
		{
			Name:    "Health",
			Context: "Global",
			Sources: []string{"WHO Bulletin", "MoHFW", "ICMR", "Reuters Health"},
			Subjects: []string{
				"seasonal fever surveillance program", "vaccine safety dashboard", "public hospital diagnostics network", "genomic monitoring pipeline",
			},
			Developments: []string{
				"publishes a verified status report", "adds new clinical guidance", "issues a preventative advisory", "extends monitoring coverage", "confirms no emergency restrictions",
			},
		},
		{
			Name:    "Geopolitics",
			Context: "Global",
			Sources: []string{"Ministry of External Affairs", "UN News", "Reuters World", "Official Embassy Briefing"},
			Subjects: []string{
				"bilateral trade dialogue", "consular advisory bulletin", "regional security consultation", "multilateral summit agenda",
			},
			Developments: []string{
				"is released as a joint statement", "is clarified by diplomatic channels", "is tracked with verified timelines", "is updated after formal consultations", "is recorded in multilateral communiques",
			},
		},
		{
			Name:    "Education",
			Context: "India",
			Sources: []string{"Ministry of Education", "UGC", "NCERT", "Press Information Bureau"},
			Subjects: []string{
				"national scholarship disbursal cycle", "board exam calendar note", "digital classrooms deployment", "teacher training certification window",
			},
			Developments: []string{
				"is shared through official education portals", "is synchronized with state implementation plans", "receives a formal clarification notice", "is verified in ministerial briefings", "is released for institutional compliance",
			},
		},
		{
			Name:    "Security",
			Context: "India",
			Sources: []string{"Ministry of Home Affairs", "National Disaster Management Authority", "PIB Fact Check", "Reuters India"},
			Subjects: []string{
				"national cyber fraud response cell", "critical infrastructure drill plan", "public emergency alert protocol", "state-level disaster preparedness grid",
			},
			Developments: []string{
				"is communicated through official control rooms", "is released with district-level instructions", "is confirmed as active and valid", "is updated to counter viral misinformation", "is reinforced through public safety outreach",
			},
		},
	}

	rng := rand.New(rand.NewSource(20260313))
	baseTime := time.Date(2026, time.March, 1, 8, 0, 0, 0, time.UTC)
	segments := make([]string, 0, len(blueprints))
	events := make([]IntelligenceEvent, 0, len(blueprints)*20)

	eventCounter := 1
	claimCounter := 1

	for _, blueprint := range blueprints {
		segments = append(segments, blueprint.Name)
		for _, subject := range blueprint.Subjects {
			for _, development := range blueprint.Developments {
				eventID := fmt.Sprintf("EVT-%03d", eventCounter)
				headline := fmt.Sprintf("%s %s", subject, development)
				summary := fmt.Sprintf(
					"Verified intelligence log: %s %s. This update is part of the March 1-13, 2026 monitoring window.",
					subject,
					development,
				)
				credibleSources := rotateSources(blueprint.Sources, eventCounter, 3)
				verificationSource := credibleSources[0]
				explanation := fmt.Sprintf(
					"Cross-checked by ALETHEIA against %s and corroborating institutional releases.",
					strings.Join(credibleSources, ", "),
				)
				timestamp := baseTime.Add(time.Duration(eventCounter-1) * 100 * time.Minute)

				claims, nextClaimCounter := generateClaimsForEvent(
					eventID,
					headline,
					timestamp,
					blueprint.Name,
					verificationSource,
					rng,
					claimCounter,
				)
				claimCounter = nextClaimCounter

				events = append(events, IntelligenceEvent{
					EventID:    eventID,
					Headline:   headline,
					Segment:    blueprint.Name,
					Topics:     topicsForSegment(blueprint.Name),
					Context:    blueprint.Context,
					Date:       timestamp.Format("2006-01-02"),
					TrendScore: eventTrendScore(blueprint.Name, rng),
					Sources:    uniqueStrings(append(cloneStrings(credibleSources), verificationSource)),
					Claims:     claims,

					Summary:            summary,
					Timestamp:          timestamp,
					CredibleSources:    credibleSources,
					VerificationSource: verificationSource,
					Explanation:        explanation,
				})
				eventCounter++
			}
		}
	}

	report := &IntelligenceReport{
		Segments:    segments,
		Events:      events,
		byClaimID:   make(map[string]intelligenceRef, len(events)*5),
		byClaimText: make(map[string]intelligenceRef, len(events)*5),
	}
	report.rebuildIndex()
	return report
}

func (r *IntelligenceReport) rebuildIndex() {
	for i, ev := range r.Events {
		for j, claim := range ev.Claims {
			ref := intelligenceRef{EventIndex: i, ClaimIndex: j}
			r.byClaimID[claim.ClaimID] = ref
			r.byClaimText[normalizeClaimText(claim.ClaimText)] = ref
		}
	}
}

func (r *IntelligenceReport) MatchClaim(claimID, claimText string) (IntelligenceMatch, bool) {
	if claimID != "" {
		if ref, ok := r.byClaimID[claimID]; ok {
			return r.matchFromRef(ref), true
		}
	}

	norm := normalizeClaimText(claimText)
	if norm == "" {
		return IntelligenceMatch{}, false
	}
	if ref, ok := r.byClaimText[norm]; ok {
		return r.matchFromRef(ref), true
	}

	bestScore := 0.0
	var bestRef intelligenceRef
	found := false
	for i, ev := range r.Events {
		for j, claim := range ev.Claims {
			score := overlapScore(norm, normalizeClaimText(claim.ClaimText))
			if score > bestScore {
				bestScore = score
				bestRef = intelligenceRef{EventIndex: i, ClaimIndex: j}
				found = true
			}
		}
	}
	if found && bestScore >= 0.62 {
		return r.matchFromRef(bestRef), true
	}
	return IntelligenceMatch{}, false
}

func (r *IntelligenceReport) Suggestions(limit int) []string {
	if limit <= 0 {
		return nil
	}

	candidates := r.AllClaims()
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Claim.Virality > candidates[j].Claim.Virality
	})

	out := make([]string, 0, limit)
	seen := make(map[string]struct{}, limit)
	for _, c := range candidates {
		if _, exists := seen[c.Claim.ClaimText]; exists {
			continue
		}
		seen[c.Claim.ClaimText] = struct{}{}
		out = append(out, c.Claim.ClaimText)
		if len(out) == limit {
			break
		}
	}
	return out
}

func (r *IntelligenceReport) AllClaims() []IntelligenceMatch {
	all := make([]IntelligenceMatch, 0, len(r.Events)*5)
	for i, ev := range r.Events {
		for j := range ev.Claims {
			all = append(all, r.matchFromRef(intelligenceRef{EventIndex: i, ClaimIndex: j}))
		}
	}
	return all
}

func (r *IntelligenceReport) TotalClaims() int {
	total := 0
	for _, ev := range r.Events {
		total += len(ev.Claims)
	}
	return total
}

func (r *IntelligenceReport) matchFromRef(ref intelligenceRef) IntelligenceMatch {
	ev := r.Events[ref.EventIndex]
	claim := ev.Claims[ref.ClaimIndex]
	return IntelligenceMatch{
		Event: ev,
		Claim: claim,
	}
}

func generateClaimsForEvent(
	eventID string,
	headline string,
	eventTimestamp time.Time,
	segment string,
	verificationSource string,
	rng *rand.Rand,
	startClaimCounter int,
) ([]IntelligenceClaim, int) {
	claims := make([]IntelligenceClaim, 0, 5)
	seen := map[string]bool{}

	type claimTemplate struct {
		TextTemplate string
		Verdict      string
		Confidence   int
		Severity     string
	}

	templates := []claimTemplate{
		{
			TextTemplate: "Breaking: %s has been suspended nationwide from midnight according to insiders.",
			Verdict:      "FALSE",
			Confidence:   89,
			Severity:     "HIGH",
		},
		{
			TextTemplate: "Viral forwards say %s proves a hidden emergency that media is suppressing.",
			Verdict:      "MISLEADING",
			Confidence:   78,
			Severity:     "MEDIUM",
		},
		{
			TextTemplate: "Leaked memo claims %s will trigger penalties for everyone within 24 hours.",
			Verdict:      "FALSE",
			Confidence:   92,
			Severity:     "CRITICAL",
		},
		{
			TextTemplate: "Unverified posts claim %s is being reversed tonight, but no bulletin is visible yet.",
			Verdict:      "UNVERIFIED",
			Confidence:   57,
			Severity:     "MEDIUM",
		},
		{
			TextTemplate: "Official update: %s, as reported by verified public sources.",
			Verdict:      "TRUE",
			Confidence:   84,
			Severity:     "LOW",
		},
		{
			TextTemplate: "Video circulating online claims %s has already caused nationwide disruption.",
			Verdict:      "MISLEADING",
			Confidence:   81,
			Severity:     "HIGH",
		},
		{
			TextTemplate: "Several viral posts allege %s will be enforced immediately without notice.",
			Verdict:      "FALSE",
			Confidence:   87,
			Severity:     "HIGH",
		},
		{
			TextTemplate: "A trending thread suggests %s proves a hidden government policy.",
			Verdict:      "MISLEADING",
			Confidence:   75,
			Severity:     "MEDIUM",
		},
		{
			TextTemplate: "Online rumors claim insiders confirmed %s late last night.",
			Verdict:      "UNVERIFIED",
			Confidence:   59,
			Severity:     "MEDIUM",
		},
		{
			TextTemplate: "Posts claim %s will impact over 70%% of citizens immediately.",
			Verdict:      "MISLEADING",
			Confidence:   77,
			Severity:     "MEDIUM",
		},
	}

	nextCounter := startClaimCounter
	for len(claims) < 5 {
		idx := rng.Intn(len(templates))
		tpl := templates[idx]
		claimText := fmt.Sprintf(tpl.TextTemplate, headline)
		if seen[claimText] {
			continue
		}
		seen[claimText] = true

		claimID := fmt.Sprintf("CLM-%04d", nextCounter)
		nextCounter++

		platform := weightedPlatform(rng)
		virality := platformViralityScore(platform, rng)
		trendScore := 40 + rng.Intn(60)
		timestamp := eventTimestamp.Add(time.Duration(rng.Intn(30)) * time.Minute)
		risk := float64(tpl.Confidence)/100*0.5 +
			float64(virality)/200000*0.3 +
			severityWeight(tpl.Severity)*0.2

		counterNarrative := fmt.Sprintf(
			"%s verification indicates the viral claim about \"%s\" is not supported by the intelligence report.",
			verificationSource,
			headline,
		)
		if tpl.Verdict == "TRUE" {
			counterNarrative = fmt.Sprintf(
				"Claim aligns with verified event logs for %s (%s).",
				segment,
				verificationSource,
			)
		}
		if tpl.Verdict == "UNVERIFIED" {
			counterNarrative = fmt.Sprintf(
				"No conclusive confirmation found in verified logs for %s. Treat as unverified until %s issues a bulletin.",
				eventID,
				verificationSource,
			)
		}

		explanation := fmt.Sprintf(
			"Claim checked against event %s. %s",
			eventID,
			counterNarrative,
		)

		claims = append(claims, IntelligenceClaim{
			ClaimID:          claimID,
			ClaimText:        claimText,
			Platform:         platform,
			PlatformIcon:     platformIcon(platform),
			Virality:         virality,
			TrendScore:       trendScore,
			Timestamp:        timestamp,
			SourceType:       "dataset",
			SourceCategory:   sourceCategory(platform),
			Verdict:          tpl.Verdict,
			Confidence:       tpl.Confidence,
			Severity:         tpl.Severity,
			RiskScore:        risk,
			VerifiedBy:       verificationSource,
			Explanation:      explanation,
			CounterNarrative: counterNarrative,
		})
	}

	return claims, nextCounter
}

func rotateSources(pool []string, offset int, count int) []string {
	if len(pool) == 0 {
		return nil
	}
	if count > len(pool) {
		count = len(pool)
	}
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, pool[(offset+i)%len(pool)])
	}
	return out
}

func platformViralityScore(platform string, rng *rand.Rand) int {
	switch platform {
	case "WhatsApp":
		return 50000 + rng.Intn(200000)
	case "X":
		return 20000 + rng.Intn(150000)
	case "YouTube":
		return 10000 + rng.Intn(80000)
	case "Reddit":
		return 5000 + rng.Intn(40000)
	case "Telegram":
		return 15000 + rng.Intn(120000)
	default:
		return 10000 + rng.Intn(50000)
	}
}

func platformIcon(platform string) string {
	switch platform {
	case "WhatsApp":
		return "\U0001F7E2"
	case "X":
		return "\u26AB"
	case "Telegram":
		return "\U0001F535"
	case "YouTube":
		return "\U0001F534"
	case "Reddit":
		return "\U0001F7E0"
	default:
		return "\u26AA"
	}
}

func sourceCategory(platform string) string {
	switch platform {
	case "WhatsApp", "Telegram":
		return "Messaging"
	case "X", "Reddit":
		return "Social"
	case "YouTube":
		return "Video"
	default:
		return "Web"
	}
}

func severityWeight(sev string) float64 {
	switch sev {
	case "CRITICAL":
		return 1.0
	case "HIGH":
		return 0.8
	case "MEDIUM":
		return 0.5
	default:
		return 0.2
	}
}

func weightedPlatform(rng *rand.Rand) string {
	r := rng.Intn(100)

	switch {
	case r < 35:
		return "WhatsApp"
	case r < 60:
		return "X"
	case r < 75:
		return "Telegram"
	case r < 90:
		return "YouTube"
	default:
		return "Reddit"
	}
}

func topicsForSegment(segment string) []string {
	switch segment {
	case "Finance":
		return cloneStrings([]string{"banking", "economy"})
	case "Technology":
		return cloneStrings([]string{"ai", "infrastructure"})
	case "Politics":
		return cloneStrings([]string{"election", "policy"})
	case "Health":
		return cloneStrings([]string{"public-health", "medicine"})
	case "Security":
		return cloneStrings([]string{"cybersecurity", "emergency"})
	case "Education":
		return cloneStrings([]string{"education", "policy"})
	case "Science":
		return cloneStrings([]string{"research", "innovation"})
	case "Geopolitics":
		return cloneStrings([]string{"diplomacy", "conflict"})
	case "Business":
		return cloneStrings([]string{"markets", "industry"})
	default:
		return cloneStrings([]string{"general"})
	}
}

func eventTrendScore(segment string, rng *rand.Rand) int {
	switch segment {
	case "Security", "Politics":
		// Security and politics narratives are more likely to trend at higher levels.
		if rng.Intn(100) < 70 {
			return 75 + rng.Intn(21) // 75-95
		}
		return 60 + rng.Intn(15) // 60-74
	default:
		if rng.Intn(100) < 20 {
			return 65 + rng.Intn(16) // 65-80
		}
		return 30 + rng.Intn(35) // 30-64
	}
}

func normalizeClaimText(input string) string {
	var builder strings.Builder
	space := false
	for _, r := range strings.ToLower(input) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			space = false
			continue
		}
		if !space {
			builder.WriteRune(' ')
			space = true
		}
	}
	return strings.TrimSpace(builder.String())
}

func overlapScore(a, b string) float64 {
	if a == "" || b == "" {
		return 0
	}
	aSet := tokenSet(a)
	bSet := tokenSet(b)
	if len(aSet) == 0 || len(bSet) == 0 {
		return 0
	}

	common := 0
	for token := range aSet {
		if _, ok := bSet[token]; ok {
			common++
		}
	}
	if common == 0 {
		return 0
	}
	den := len(aSet)
	if len(bSet) > den {
		den = len(bSet)
	}
	return float64(common) / float64(den)
}

func tokenSet(s string) map[string]struct{} {
	parts := strings.Fields(s)
	set := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		if len(p) <= 2 {
			continue
		}
		set[p] = struct{}{}
	}
	return set
}
