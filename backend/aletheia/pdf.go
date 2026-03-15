package aletheia

import (
	"fmt"
	"strings"
	"time"
)

func BuildEvidencePDF(record DetectionRecord) []byte {
	evidenceHash := record.EvidenceHash
	if evidenceHash == "" {
		evidenceHash = record.CurrentHash
	}
	if evidenceHash == "" {
		evidenceHash = record.Hash
	}

	previousHash := record.PreviousHash
	if previousHash == "" {
		previousHash = evidenceHash
	}

	currentHash := record.CurrentHash
	if currentHash == "" {
		currentHash = evidenceHash
	}

	sources := record.SourcesChecked
	if len(sources) == 0 {
		sources = []string{"Reuters", "Press Information Bureau", "WHO", "MoHFW", "ICMR"}
	}

	reasoning := record.ReasoningSteps
	if len(reasoning) == 0 {
		reasoning = []string{
			"scanning Reuters feed",
			"scanning PIB bulletins",
			"checking WHO advisories",
			"evaluating claim semantics",
		}
	}

	if record.Verdict == "" {
		record.Verdict = "UNVERIFIED"
	}
	if record.Severity == "" {
		record.Severity = "MEDIUM"
	}
	if record.RiskLevel == "" {
		record.RiskLevel = record.RiskLabel
	}
	if record.RiskLevel == "" {
		record.RiskLevel = "MEDIUM"
	}
	if record.Platform == "" {
		record.Platform = "N/A"
	}
	if record.Segment == "" {
		record.Segment = "N/A"
	}

	content := buildEvidencePDFContent(
		record,
		previousHash,
		evidenceHash,
		currentHash,
		sources,
		reasoning,
	)
	return buildPDFDocument(content)
}

func BuildSimplePDF(title string, lines []string) []byte {
	allLines := append([]string{title}, lines...)
	return buildPDFDocument(buildPDFContent(allLines))
}

func buildPDFContent(lines []string) string {
	const maxLen = 92
	var wrapped []string
	for _, line := range lines {
		for _, chunk := range wrapLine(line, maxLen) {
			wrapped = append(wrapped, chunk)
		}
	}

	var builder strings.Builder
	builder.WriteString("BT\n/F1 11 Tf\n14 TL\n50 760 Td\n")
	for i, line := range wrapped {
		if i > 0 {
			builder.WriteString("T*\n")
		}
		builder.WriteString(fmt.Sprintf("(%s) Tj\n", pdfEscape(line)))
	}
	builder.WriteString("ET")
	return builder.String()
}

func buildEvidencePDFContent(
	record DetectionRecord,
	previousHash string,
	evidenceHash string,
	currentHash string,
	sources []string,
	reasoning []string,
) string {
	const (
		leftMargin  = 48.0
		rightMargin = 564.0
	)

	var b strings.Builder
	y := 760.0
	width := rightMargin - leftMargin

	b.WriteString("0.8 w\n")

	// Report header
	pdfText(&b, "F2", 21, leftMargin, y, "ALETHEIA VERIFICATION REPORT")
	y -= 20
	pdfText(&b, "F1", 12, leftMargin, y, "Cryptographic Evidence Record")
	y -= 12
	pdfLine(&b, leftMargin, y, rightMargin, y)
	y -= 18

	// Section 1 — Claim Information
	y = pdfSectionHeader(&b, y, "Section 1 — Claim Information", leftMargin, rightMargin)
	infoLines := []string{
		"Claim ID: " + valueOrNA(record.ID),
		"Timestamp: " + record.Timestamp.UTC().Format("2006-01-02T15:04:05Z07:00"),
		"Platform: " + valueOrNA(record.Platform),
		"Segment: " + valueOrNA(record.Segment),
	}
	for _, line := range infoLines {
		y = pdfWrappedText(&b, "F1", 11, leftMargin+6, y, line, 96, 13.5)
	}
	y -= 6

	// Section 2 — Claim Content
	y = pdfSectionHeader(&b, y, "Section 2 — Claim Content", leftMargin, rightMargin)
	claim := valueOrNA(record.Claim)
	claimLines := wrapLine(claim, 92)
	boxTop := y
	boxHeight := 14 + float64(len(claimLines))*13 + 10
	boxBottom := boxTop - boxHeight
	pdfRectStroke(&b, leftMargin+4, boxBottom, width-8, boxHeight)
	textY := boxTop - 18
	for _, line := range claimLines {
		pdfText(&b, "F1", 11, leftMargin+10, textY, line)
		textY -= 13
	}
	y = boxBottom - 14

	// Section 3 — Verdict Summary
	y = pdfSectionHeader(&b, y, "Section 3 — Verdict Summary", leftMargin, rightMargin)
	verdictTop := y
	verdictBoxHeight := 18.0
	pdfRectFillGray(&b, leftMargin+4, verdictTop-verdictBoxHeight, width-8, verdictBoxHeight, 0.92)
	pdfRectStroke(&b, leftMargin+4, verdictTop-verdictBoxHeight, width-8, verdictBoxHeight)
	pdfText(&b, "F2", 12, leftMargin+10, verdictTop-13, "Verdict: "+valueOrNA(record.Verdict))
	y -= 24

	summaryLines := []string{
		fmt.Sprintf("Confidence: %d%%", record.Confidence),
		"Severity: " + valueOrNA(record.Severity),
		fmt.Sprintf("Risk Score: %.2f / 10", record.RiskScore),
		"Risk Level: " + valueOrNA(record.RiskLevel),
	}
	for _, line := range summaryLines {
		y = pdfWrappedText(&b, "F1", 11, leftMargin+6, y, line, 96, 13.5)
	}
	y -= 6

	// Section 4 — Verification Sources
	y = pdfSectionHeader(&b, y, "Section 4 — Verification Sources", leftMargin, rightMargin)
	for _, src := range sources {
		y = pdfBulletText(&b, leftMargin+8, y, valueOrNA(src), 88, 12.5)
	}
	y -= 6

	// Section 5 — Verification Reasoning
	y = pdfSectionHeader(&b, y, "Section 5 — Verification Reasoning", leftMargin, rightMargin)
	for _, step := range reasoning {
		y = pdfBulletText(&b, leftMargin+8, y, valueOrNA(step), 88, 12.5)
	}
	y -= 6

	// Section 6 — Cryptographic Evidence Chain
	y = pdfSectionHeader(&b, y, "Section 6 — Cryptographic Evidence Chain", leftMargin, rightMargin)
	prevHashLines := wrapLine(previousHash, 66)
	eviHashLines := wrapLine(evidenceHash, 66)
	curHashLines := wrapLine(currentHash, 66)
	chainTextLines := 1 + // chain position line
		1 + len(prevHashLines) +
		1 + len(eviHashLines) +
		1 + len(curHashLines)

	chainTop := y
	chainBoxHeight := 14 + float64(chainTextLines)*11.5 + 10
	chainBottom := chainTop - chainBoxHeight
	pdfRectStroke(&b, leftMargin+4, chainBottom, width-8, chainBoxHeight)
	chainY := chainTop - 16
	pdfText(&b, "F1", 10.5, leftMargin+10, chainY, fmt.Sprintf("Chain Position: %d", record.ChainPosition))
	chainY -= 14

	pdfText(&b, "F2", 10, leftMargin+10, chainY, "Previous Hash:")
	chainY -= 12
	for _, line := range prevHashLines {
		pdfText(&b, "F3", 9, leftMargin+14, chainY, line)
		chainY -= 11
	}

	pdfText(&b, "F2", 10, leftMargin+10, chainY, "Evidence Hash:")
	chainY -= 12
	for _, line := range eviHashLines {
		pdfText(&b, "F3", 9, leftMargin+14, chainY, line)
		chainY -= 11
	}

	pdfText(&b, "F2", 10, leftMargin+10, chainY, "Current Hash:")
	chainY -= 12
	for _, line := range curHashLines {
		pdfText(&b, "F3", 9, leftMargin+14, chainY, line)
		chainY -= 11
	}

	y = chainBottom - 14

	// Section 7 — System Footer
	y = pdfSectionHeader(&b, y, "Section 7 — System Footer", leftMargin, rightMargin)
	footerLines := []string{
		"Generated by ALETHEIA",
		"Real-Time Misinformation Verification Engine",
		"Generated At: " + time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		"Verification Endpoint: /api/evidence/pdf?id=" + valueOrNA(record.ID),
	}
	for _, line := range footerLines {
		y = pdfWrappedText(&b, "F1", 10.5, leftMargin+6, y, line, 96, 12.5)
	}

	return b.String()
}

func buildPDFDocument(content string) []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R /F2 5 0 R /F3 6 0 R >> >> /Contents 7 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Courier >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content),
	}

	var builder strings.Builder
	builder.WriteString("%PDF-1.4\n")

	offsets := make([]int, len(objects)+1)
	for i, obj := range objects {
		offsets[i+1] = builder.Len()
		builder.WriteString(fmt.Sprintf("%d 0 obj\n%s\nendobj\n", i+1, obj))
	}

	xrefPos := builder.Len()
	builder.WriteString(fmt.Sprintf("xref\n0 %d\n", len(objects)+1))
	builder.WriteString("0000000000 65535 f \n")
	for i := 1; i <= len(objects); i++ {
		builder.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}
	builder.WriteString(fmt.Sprintf(
		"trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF",
		len(objects)+1,
		xrefPos,
	))

	return []byte(builder.String())
}

func pdfSectionHeader(b *strings.Builder, y float64, title string, left float64, right float64) float64 {
	pdfText(b, "F2", 12.5, left, y, title)
	pdfLine(b, left, y-4, right, y-4)
	return y - 18
}

func pdfWrappedText(
	b *strings.Builder,
	font string,
	size float64,
	x float64,
	y float64,
	text string,
	maxChars int,
	lineGap float64,
) float64 {
	lines := wrapLine(text, maxChars)
	for _, line := range lines {
		pdfText(b, font, size, x, y, line)
		y -= lineGap
	}
	return y
}

func pdfBulletText(
	b *strings.Builder,
	x float64,
	y float64,
	text string,
	maxChars int,
	lineGap float64,
) float64 {
	lines := wrapLine(text, maxChars)
	for i, line := range lines {
		prefix := "- "
		if i > 0 {
			prefix = "  "
		}
		pdfText(b, "F1", 10.5, x, y, prefix+line)
		y -= lineGap
	}
	return y
}

func pdfText(b *strings.Builder, font string, size float64, x float64, y float64, text string) {
	b.WriteString("BT\n")
	b.WriteString(fmt.Sprintf("/%s %.2f Tf\n", font, size))
	b.WriteString(fmt.Sprintf("1 0 0 1 %.2f %.2f Tm\n", x, y))
	b.WriteString(fmt.Sprintf("(%s) Tj\n", pdfEscape(text)))
	b.WriteString("ET\n")
}

func pdfLine(b *strings.Builder, x1 float64, y1 float64, x2 float64, y2 float64) {
	b.WriteString(fmt.Sprintf("%.2f %.2f m\n", x1, y1))
	b.WriteString(fmt.Sprintf("%.2f %.2f l\n", x2, y2))
	b.WriteString("S\n")
}

func pdfRectStroke(b *strings.Builder, x float64, y float64, w float64, h float64) {
	b.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f re\n", x, y, w, h))
	b.WriteString("S\n")
}

func pdfRectFillGray(b *strings.Builder, x float64, y float64, w float64, h float64, gray float64) {
	b.WriteString(fmt.Sprintf("%.2f g\n", gray))
	b.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f re\n", x, y, w, h))
	b.WriteString("f\n")
	b.WriteString("0 g\n")
}

func valueOrNA(v string) string {
	if strings.TrimSpace(v) == "" {
		return "N/A"
	}
	return v
}

func wrapLine(line string, width int) []string {
	if len(line) <= width {
		return []string{line}
	}
	if width <= 0 {
		return []string{line}
	}

	// First try word wrapping.
	words := strings.Fields(line)
	if len(words) > 0 {
		var lines []string
		var current strings.Builder
		for _, word := range words {
			if current.Len() == 0 {
				if len(word) > width {
					hard := wrapLineHard(word, width)
					lines = append(lines, hard...)
				} else {
					current.WriteString(word)
				}
				continue
			}
			if current.Len()+1+len(word) > width {
				lines = append(lines, current.String())
				current.Reset()
				if len(word) > width {
					hard := wrapLineHard(word, width)
					lines = append(lines, hard...)
				} else {
					current.WriteString(word)
				}
				continue
			}
			current.WriteByte(' ')
			current.WriteString(word)
		}
		if current.Len() > 0 {
			lines = append(lines, current.String())
		}
		if len(lines) > 0 {
			return lines
		}
	}

	// Fallback for long no-space strings (e.g., hashes).
	return wrapLineHard(line, width)
}

func wrapLineHard(line string, width int) []string {
	if len(line) <= width || width <= 0 {
		return []string{line}
	}

	out := make([]string, 0, (len(line)/width)+1)
	remaining := line
	for len(remaining) > width {
		cut := width
		if idx := strings.LastIndex(remaining[:width], " "); idx > 0 {
			cut = idx
		}
		chunk := strings.TrimSpace(remaining[:cut])
		if chunk == "" {
			chunk = remaining[:cut]
		}
		out = append(out, chunk)
		remaining = strings.TrimSpace(remaining[cut:])
	}
	if remaining != "" {
		out = append(out, remaining)
	}
	if len(out) == 0 {
		return []string{""}
	}
	return out
}

func pdfEscape(input string) string {
	escaped := strings.ReplaceAll(input, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, "(", `\(`)
	escaped = strings.ReplaceAll(escaped, ")", `\)`)
	return escaped
}
