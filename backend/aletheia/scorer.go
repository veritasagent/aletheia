package aletheia

import "math"

func Score(event MisinformationEvent) (float64, string) {
	score := (float64(event.Confidence) / 100.0 * 0.5) +
		(math.Min(float64(event.Virality)/500000.0, 1.0) * 0.3) +
		(severityScore(event.Severity) * 0.2)

	var label string
	if share := score; share <= 0.3 {
		label = "LOW"
	} else if score <= 0.6 {
		label = "MEDIUM"
	} else if score <= 0.8 {
		label = "HIGH"
	} else {
		label = "CRITICAL"
	}

	return score, label
}

func severityScore(severity string) float64 {
	switch severity {
	case "LOW":
		return 0.2
	case "MEDIUM":
		return 0.5
	case "HIGH":
		return 0.8
	case "CRITICAL":
		return 1.0
	default:
		return 0.2
	}
}
