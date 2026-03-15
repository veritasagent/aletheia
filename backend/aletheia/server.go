package aletheia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AletheiaServer struct {
	chain   *DetectionChain
	handler *Handler
	watcher *Watcher
	report  *IntelligenceReport
	port    int

	clients map[chan string]bool
	mu      sync.Mutex
}

func NewServer(
	chain *DetectionChain,
	handler *Handler,
	watcher *Watcher,
	report *IntelligenceReport,
	port int,
) *AletheiaServer {
	return &AletheiaServer{
		chain:   chain,
		handler: handler,
		watcher: watcher,
		report:  report,
		port:    port,
		clients: make(map[chan string]bool),
	}
}

func (s *AletheiaServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleDashboard)
	mux.HandleFunc("/api/stream", s.handleStream)
	mux.HandleFunc("/api/chain", s.handleChain)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/verify", s.handleVerify)
	mux.HandleFunc("/api/suggestions", s.handleSuggestions)
	mux.HandleFunc("/api/evidence/pdf", s.handleEvidencePDF)
	mux.HandleFunc("/api/independent-verify", s.handleIndependentVerify)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("ALETHEIA sandbox dashboard running on http://localhost:%d\n", s.port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server failed: %v\n", err)
	}
}

func (s *AletheiaServer) BroadcastDetection(record DetectionRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record = normalizeRecordForAPI(record)
	bytes, err := json.Marshal(record)
	if err != nil {
		return
	}
	message := fmt.Sprintf("data: %s\n\n", string(bytes))
	for client := range s.clients {
		select {
		case client <- message:
		default:
		}
	}
}

func (s *AletheiaServer) handleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	clientChan := make(chan string, 16)
	s.mu.Lock()
	s.clients[clientChan] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, clientChan)
		s.mu.Unlock()
		close(clientChan)
	}()

	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg, ok := <-clientChan:
			if !ok {
				return
			}
			fmt.Fprint(w, msg)
			flusher.Flush()
		}
	}
}

func (s *AletheiaServer) handleChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	records := s.chain.Export()
	for i := range records {
		records[i] = normalizeRecordForAPI(records[i])
	}
	_ = json.NewEncoder(w).Encode(records)
}

func (s *AletheiaServer) handleStats(w http.ResponseWriter, r *http.Request) {
	records := s.chain.Export()
	var falseCount, misleadingCount, trueCount, unverifiedCount, criticalCount int
	for _, record := range records {
		switch record.Verdict {
		case "FALSE":
			falseCount++
		case "MISLEADING":
			misleadingCount++
		case "TRUE":
			trueCount++
		case "UNVERIFIED":
			unverifiedCount++
		}
		if record.Severity == "CRITICAL" {
			criticalCount++
		}
	}

	stats := map[string]interface{}{
		"total":            len(records),
		"false_count":      falseCount,
		"misleading_count": misleadingCount,
		"true_count":       trueCount,
		"unverified_count": unverifiedCount,
		"critical_count":   criticalCount,
		"chain_valid":      s.chain.Verify(),
	}
	if s.report != nil {
		stats["dataset_events"] = len(s.report.Events)
		stats["dataset_claims"] = s.report.TotalClaims()
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}

func (s *AletheiaServer) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody struct {
		Claim string `json:"claim"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if reqBody.Claim == "" {
		http.Error(w, "claim is required", http.StatusBadRequest)
		return
	}

	event := MisinformationEvent{
		ID:        fmt.Sprintf("USR-%d", time.Now().UnixNano()),
		Claim:     reqBody.Claim,
		Platform:  "Operator Console",
		Virality:  5000,
		Timestamp: time.Now().UTC(),
	}

	classified := Classify(event)
	s.handler.Handle(classified)

	record, ok := s.chain.FindByID(classified.ID)
	if !ok {
		http.Error(w, "record not found after verification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(normalizeRecordForAPI(record))
}

func (s *AletheiaServer) handleSuggestions(w http.ResponseWriter, r *http.Request) {
	if s.report == nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]string{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.report.Suggestions(10))
}

func (s *AletheiaServer) handleEvidencePDF(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	record, ok := s.chain.FindByID(id)
	if !ok {
		http.Error(w, "record not found", http.StatusNotFound)
		return
	}
	record = normalizeRecordForAPI(record)
	pdfBytes := BuildEvidencePDF(record)
	filename := fmt.Sprintf("%s-evidence.pdf", record.ID)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	_, _ = w.Write(pdfBytes)
}

func (s *AletheiaServer) handleIndependentVerify(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	records := s.chain.Export()
	var record DetectionRecord
	index := -1
	for i := range records {
		if records[i].ID == id {
			record = normalizeRecordForAPI(records[i])
			index = i
			break
		}
	}
	if index < 0 {
		http.Error(w, "record not found", http.StatusNotFound)
		return
	}

	previousValid := true
	if index > 0 {
		previousValid = chooseEvidenceHash(normalizeRecordForAPI(records[index-1])) == record.PreviousHash
	}
	expectedHash := computeChainHash(record.Claim, record.Verdict, record.Timestamp, record.PreviousHash)
	valid := previousValid && (chooseEvidenceHash(record) == expectedHash)

	response := map[string]interface{}{
		"id":             id,
		"valid":          valid,
		"stored_hash":    chooseEvidenceHash(record),
		"expected_hash":  expectedHash,
		"previous_valid": previousValid,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (s *AletheiaServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}

func joinOrDefault(items []string, fallback string) string {
	if len(items) == 0 {
		return fallback
	}
	return strings.Join(items, ", ")
}

func chooseEvidenceHash(record DetectionRecord) string {
	if record.EvidenceHash != "" {
		return record.EvidenceHash
	}
	if record.CurrentHash != "" {
		return record.CurrentHash
	}
	return record.Hash
}

func normalizeRecordForAPI(record DetectionRecord) DetectionRecord {
	if record.RiskLevel == "" {
		record.RiskLevel = record.RiskLabel
	}
	if record.SourcesChecked == nil {
		record.SourcesChecked = []string{}
	}
	if len(record.ReasoningSteps) == 0 {
		record.ReasoningSteps = append([]string{}, defaultReasoningSteps()...)
	}
	if record.EvidenceHash == "" {
		record.EvidenceHash = chooseEvidenceHash(record)
	}
	if record.CurrentHash == "" {
		record.CurrentHash = record.EvidenceHash
		if record.CurrentHash == "" {
			record.CurrentHash = record.Hash
		}
	}
	if record.PreviousHash == "" {
		record.PreviousHash = record.EvidenceHash
		if record.PreviousHash == "" {
			record.PreviousHash = record.CurrentHash
		}
	}
	if record.ChainPosition <= 0 {
		record.ChainPosition = 1
	}
	return record
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>ALETHEIA Sandbox</title>
  <style>
    :root {
      --bg: #0b121f;
      --panel: #11203a;
      --panel-alt: #0f1a2f;
      --line: #294166;
      --accent: #12d3a6;
      --danger: #ff5f57;
      --warn: #ffc145;
      --ok: #32d583;
      --text: #e8f0ff;
      --muted: #9ab1cc;
      --font: "IBM Plex Sans", "Segoe UI", sans-serif;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: var(--font);
      background: radial-gradient(1200px 700px at 10% -10%, #18335f 0%, var(--bg) 60%);
      color: var(--text);
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }
    header {
      padding: 14px 18px;
      border-bottom: 1px solid var(--line);
      background: rgba(11, 18, 31, 0.9);
      display: flex;
      justify-content: space-between;
      gap: 12px;
      flex-wrap: wrap;
    }
    .title { font-weight: 700; letter-spacing: 0.08em; }
    .subtitle { color: var(--muted); font-size: 12px; }
    .stats { display: flex; gap: 8px; flex-wrap: wrap; }
    .stat {
      font-size: 12px;
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 999px;
      padding: 6px 10px;
    }
    main {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 12px;
      padding: 12px;
      flex: 1;
      min-height: 0;
    }
    .column {
      border: 1px solid var(--line);
      border-radius: 12px;
      overflow: hidden;
      background: linear-gradient(180deg, var(--panel), var(--panel-alt));
      display: flex;
      flex-direction: column;
      min-height: 0;
    }
    .column h2 {
      margin: 0;
      padding: 12px;
      border-bottom: 1px solid var(--line);
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: 0.12em;
      color: var(--muted);
    }
    .list {
      overflow: auto;
      padding: 10px;
      display: flex;
      flex-direction: column;
      gap: 8px;
    }
    .card {
      border: 1px solid var(--line);
      border-left: 3px solid var(--accent);
      border-radius: 8px;
      background: rgba(11, 18, 31, 0.55);
      padding: 10px;
      font-size: 13px;
      line-height: 1.35;
    }
    .meta {
      color: var(--muted);
      font-size: 11px;
      margin-bottom: 6px;
      display: flex;
      justify-content: space-between;
      gap: 8px;
      flex-wrap: wrap;
    }
    .badge {
      display: inline-block;
      font-size: 11px;
      font-weight: 700;
      padding: 2px 8px;
      border-radius: 999px;
      margin-bottom: 6px;
    }
    .badge.FALSE { background: rgba(255,95,87,0.2); color: var(--danger); }
    .badge.MISLEADING { background: rgba(255,193,69,0.2); color: var(--warn); }
    .badge.TRUE { background: rgba(50,213,131,0.2); color: var(--ok); }
    .badge.UNVERIFIED { background: rgba(154,177,204,0.2); color: var(--muted); }
    ul.steps {
      margin: 0;
      padding-left: 16px;
      color: var(--muted);
      font-size: 12px;
    }
    .actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      margin-top: 8px;
    }
    .btn, .link-btn {
      font: inherit;
      font-size: 12px;
      border-radius: 6px;
      padding: 6px 9px;
      border: 1px solid var(--line);
      background: rgba(18, 211, 166, 0.15);
      color: var(--text);
      text-decoration: none;
      cursor: pointer;
    }
    .btn:hover, .link-btn:hover {
      border-color: var(--accent);
      background: rgba(18, 211, 166, 0.22);
    }
    .hash { font-family: ui-monospace, "Cascadia Code", monospace; font-size: 11px; color: var(--muted); margin-top: 6px; }
    footer {
      border-top: 1px solid var(--line);
      padding: 10px 12px;
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      align-items: center;
      background: rgba(11, 18, 31, 0.92);
    }
    #claimInput {
      flex: 1;
      min-width: 220px;
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 10px;
      background: #0a1528;
      color: var(--text);
    }
    #verifyBtn {
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 10px 14px;
      background: var(--accent);
      color: #041417;
      font-weight: 700;
      cursor: pointer;
    }
    #verifyBtn:disabled { opacity: 0.5; cursor: not-allowed; }
    @media (max-width: 1080px) {
      main { grid-template-columns: 1fr; }
      .list { max-height: 40vh; }
    }
  </style>
</head>
<body>
  <header>
    <div>
      <div class="title">ALETHEIA SANDBOX MONITOR</div>
      <div class="subtitle">Simulated social feed driven by intelligence report (March 1-13, 2026)</div>
    </div>
    <div class="stats">
      <div class="stat">TOTAL <span id="s-total">0</span></div>
      <div class="stat">FALSE <span id="s-false">0</span></div>
      <div class="stat">MISLEADING <span id="s-misleading">0</span></div>
      <div class="stat">UNVERIFIED <span id="s-unverified">0</span></div>
      <div class="stat">CRITICAL <span id="s-critical">0</span></div>
      <div class="stat">CHAIN <span id="s-chain">OK</span></div>
    </div>
  </header>

  <main>
    <section class="column">
      <h2>Claim Stream</h2>
      <div id="stream" class="list"></div>
    </section>
    <section class="column">
      <h2>Agent Analysis</h2>
      <div id="analysis" class="list"></div>
    </section>
    <section class="column">
      <h2>Evidence Vault</h2>
      <div id="evidence" class="list"></div>
    </section>
  </main>

  <footer>
    <input id="claimInput" list="claimHints" placeholder="Submit claim for verification..." />
    <datalist id="claimHints"></datalist>
    <button id="verifyBtn">Verify Claim</button>
  </footer>

  <script>
    const seen = new Set();
    const stream = document.getElementById('stream');
    const analysis = document.getElementById('analysis');
    const evidence = document.getElementById('evidence');
    const verifyBtn = document.getElementById('verifyBtn');
    const claimInput = document.getElementById('claimInput');
    const claimHints = document.getElementById('claimHints');

    function shortHash(hash) {
      if (!hash) return 'n/a';
      return hash.length > 20 ? hash.slice(0, 20) + '...' : hash;
    }

    function safeSteps(record) {
      if (Array.isArray(record.ReasoningChain) && record.ReasoningChain.length > 0) return record.ReasoningChain;
      return ['Scanning Reuters', 'Scanning PIB', 'Scanning WHO', 'Evaluating claim', 'Verdict determined'];
    }

    function renderRecord(record) {
      record = normalizeRecord(record);
      if (!record || !record.ID || seen.has(record.ID)) return;
      seen.add(record.ID);

      const streamCard = document.createElement('article');
      streamCard.className = 'card';
      streamCard.innerHTML =
        '<div class="meta"><span>' + (record.Platform || 'Unknown') + '</span><span>' + (record.Segment || 'General') + '</span></div>' +
        '<div>' + record.Claim + '</div>' +
        '<div class="meta"><span>ID: ' + record.ID + '</span><span>Virality: ' + Number(record.Virality || 0).toLocaleString() + '</span></div>';
      stream.prepend(streamCard);

      const analysisCard = document.createElement('article');
      analysisCard.className = 'card';
      const steps = safeSteps(record).map(function(step) { return '<li>' + step + '</li>'; }).join('');
      analysisCard.innerHTML =
        '<span class="badge ' + record.Verdict + '">' + record.Verdict + '</span>' +
        '<div><strong>Confidence:</strong> ' + record.Confidence + '%</div>' +
        '<div><strong>Reasoning:</strong> ' + (record.ReasoningSummary || record.Reason || 'No reasoning summary') + '</div>' +
        '<ul class="steps">' + steps + '</ul>';
      analysis.prepend(analysisCard);

      const evidenceCard = document.createElement('article');
      evidenceCard.className = 'card';
      evidenceCard.innerHTML =
        '<div><strong>Record #</strong> ' + Number(record.ChainPosition || 0) + '</div>' +
        '<span class="badge ' + record.Verdict + '">' + record.Verdict + '</span>' +
        '<div><strong>Severity:</strong> ' + (record.Severity || 'N/A') + '</div>' +
        '<div><strong>Confidence:</strong> ' + record.Confidence + '%</div>' +
        '<div><strong>Sources:</strong> ' + ((record.SourcesChecked || []).join(', ') || 'N/A') + '</div>' +
        '<div><strong>Platform Spread:</strong> ' + (record.Platform || 'N/A') + ' / ' + Number(record.Virality || 0).toLocaleString() + '</div>' +
        '<div class="hash">sha256: ' + shortHash(record.EvidenceHash || record.CurrentHash || record.Hash) + '</div>' +
        '<div class="actions">' +
          '<a class="link-btn" target="_blank" href="/api/evidence/pdf?id=' + encodeURIComponent(record.ID) + '">Download Evidence PDF</a>' +
          '<button class="btn" data-id="' + record.ID + '">Independent Verification Tool</button>' +
        '</div>';
      const verifyButton = evidenceCard.querySelector('button');
      verifyButton.addEventListener('click', async function() {
        const id = verifyButton.getAttribute('data-id');
        const res = await fetch('/api/independent-verify?id=' + encodeURIComponent(id));
        const payload = await res.json();
        alert(payload.valid ? 'Hash verified for ' + id : 'Hash mismatch for ' + id);
      });
      evidence.prepend(evidenceCard);
    }

    function normalizeRecord(raw) {
      const rec = raw || {};
      return {
        ID: rec.ID || rec.id || '',
        Claim: rec.Claim || rec.claim_text || '',
        Platform: rec.Platform || rec.platform || '',
        Segment: rec.Segment || rec.segment || '',
        Virality: rec.Virality || rec.virality || 0,
        Verdict: rec.Verdict || rec.verdict || '',
        Confidence: rec.Confidence || rec.confidence || 0,
        Severity: rec.Severity || rec.severity || '',
        SourcesChecked: rec.SourcesChecked || rec.sources || [],
        ReasoningSummary: rec.ReasoningSummary || rec.reasoning_summary || rec.reason || '',
        Reason: rec.Reason || rec.reason || '',
        ReasoningChain: rec.ReasoningChain || rec.reasoning_chain || rec.reasoning_steps || [],
        CounterNarrative: rec.CounterNarrative || rec.counter_narrative || '',
        Hash: rec.Hash || rec.hash || '',
        EvidenceHash: rec.EvidenceHash || rec.evidence_hash || '',
        CurrentHash: rec.CurrentHash || rec.current_hash || rec.evidence_hash || rec.hash || '',
        PreviousHash: rec.PreviousHash || rec.previous_hash || '',
        ChainPosition: rec.ChainPosition || rec.chain_position || 0
      };
    }

    async function updateStats() {
      try {
        const res = await fetch('/api/stats');
        const stats = await res.json();
        document.getElementById('s-total').innerText = stats.total;
        document.getElementById('s-false').innerText = stats.false_count;
        document.getElementById('s-misleading').innerText = stats.misleading_count;
        document.getElementById('s-unverified').innerText = stats.unverified_count;
        document.getElementById('s-critical').innerText = stats.critical_count;
        document.getElementById('s-chain').innerText = stats.chain_valid ? 'OK' : 'FAILED';
      } catch (err) {
        console.error(err);
      }
    }

    async function loadInitial() {
      const chainRes = await fetch('/api/chain');
      const records = await chainRes.json();
      records.forEach(renderRecord);
      updateStats();
    }

    async function loadSuggestions() {
      try {
        const res = await fetch('/api/suggestions');
        const suggestions = await res.json();
        claimHints.innerHTML = '';
        suggestions.forEach(function(item) {
          const option = document.createElement('option');
          option.value = item;
          claimHints.appendChild(option);
        });
      } catch (err) {
        console.error(err);
      }
    }

    verifyBtn.addEventListener('click', async function() {
      const claim = claimInput.value.trim();
      if (!claim) return;
      verifyBtn.disabled = true;
      try {
        await fetch('/api/verify', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ claim: claim })
        });
        claimInput.value = '';
      } finally {
        verifyBtn.disabled = false;
      }
    });

    const source = new EventSource('/api/stream');
    source.onmessage = function(evt) {
      try {
        const record = JSON.parse(evt.data);
        renderRecord(record);
        updateStats();
      } catch (err) {
        console.error(err);
      }
    };

    loadInitial();
    loadSuggestions();
    setInterval(updateStats, 5000);
    setInterval(loadSuggestions, 60000);
  </script>
</body>
</html>`
