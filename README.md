# ALETHEIA — Truth Engine
HACKANOVA 5.0 · Agentic AI Track · Team Bytes_and_Bros

ALETHEIA is a real-time, tamper-proof misinformation intelligence dashboard and verification engine. It utilizes an autonomous pipeline (Go backend) powered by Gemini 2.5-flash to intercept, analyze, and cryptographically ledger digital claims across the web before they spread.

## Architecture Pipeline
\`\`\`text
[ Social Media / Web / Simulation Feed ]
           ↓
    ( Ingestion Engine )
           ↓
[ Gemini 2.5-flash Verification ] ← ( Web Searching & Ground Truth DB )
           ↓
  ( Agentic Reasoning )
           ↓
[ Cryptographic Hash Chain ] ← H(d) = SHA256(claim || verdict || timestamp || H(d-1))
           ↓
  ( SSE Streaming API )
           ↓
[ Next.js 14 Web Dashboard ] ← ( You are here )
\`\`\`

## Quick Start
1. Ensure the Go backend API is running correctly \`go run ./... --sandbox --web\`
2. Navigate into the UI dir: \`cd aletheia-ui\`
3. Install packages: \`npm install\`
4. Run Development Server: \`npm run dev\`
5. Open browser at: [http://localhost:3000](http://localhost:3000)

## Environment Variables
Create a \`.env.local\` file in the \`aletheia-ui\` root.

| Variable | Description | Default |
| :--- | :--- | :--- |
| \`NEXT_PUBLIC_API_URL\` | Target backend endpoint URL | \`http://localhost:8080\` |
| \`NEXT_PUBLIC_APP_NAME\` | Application Title | \`ALETHEIA\` |
| \`NEXT_PUBLIC_CHAIN_VERSION\` | Current Chain Build tag | \`v0.7.0-agent\` |

## Backend API Endpoints
All API calls from the frontend to \`/api/*\` are automatically proxied via 'rewrites' to \`http://localhost:8080\`.

| Endpoint | Method | Purpose |
| :--- | :--- | :--- |
| \`/api/stream\` | \`GET\` (SSE) | Server-Sent Events stream for live intercepts |
| \`/api/verify\` | \`POST\` | Submits an independent claim for analysis |
| \`/api/stats\` | \`GET\` | Returns global verification statistics |
| \`/api/chain\` | \`GET\` | Returns the complete cryptographic ledger |

## Tech Stack
* **Frontend**: Next.js 14 App Router, TypeScript, TailwindCSS, Framer Motion
* **Backend**: Go (Golang)
* **AI Model**: Gemini 2.5-flash
* **Infrastructure**: Immutable Cryptographic Chain State
