
# Cilium-Shield: The AI-DLP Firewall

[![eBPF Day Hackathon 2025](https://img.shields.io/badge/eBPF%20Day%20Hackathon-2025-blue)](https://ebpfday.com/)
[![Built with Cilium](https://img.shields.io/badge/Built%20with-Cilium-9900ef)](https://cilium.io/)
[![Wasm Powered](https://img.shields.io/badge/Wasm-Powered-6140ef)](https://webassembly.org/)

**Cilium-Shield is a Kubernetes-native Data Leakage Prevention (DLP) system designed to secure Large Language Model (LLM) adoption by preventing sensitive data from leaving your cluster.**

It leverages the power of **Cilium and eBPF** to intercept outbound traffic to AI APIs and uses a high-performance **WebAssembly (Wasm)** filter to redact PII, financial data, and credentials in real-time.

---

## 🚀 The Problem: "Shadow AI" is a Critical Security Blindspot

As organizations race to adopt Generative AI, developers frequently use public APIs from providers like OpenAI, Anthropic, and Google. This creates a significant risk:

*   **Sensitive Data Leakage:** Developers might accidentally or intentionally paste internal source code, customer data (PII), or database schemas into prompts.
*   **Compliance Violations:** This "Shadow AI" usage can lead to severe GDPR, CCPA, and PCI-DSS compliance violations.
*   **Traditional Tools are Ineffective:** Network firewalls are blind to the contents of encrypted (TLS) traffic, and sidecar-based service meshes add significant latency and resource overhead.

Cilium-Shield solves this by providing payload-aware security at the source, without compromising performance.

## ✨ How It Works: eBPF + Wasm for Ultimate L7 Security

Cilium-Shield extends Cilium's powerful eBPF-based networking with a lightweight, sandboxed Wasm module for Layer 7 inspection.

```mermaid
graph LR
    subgraph Kubernetes Cluster
        A[App Pod] -->|HTTP Req with PII| B(Cilium Node Agent)
        B -->|L7 Redirect| C{Envoy Proxy}
        
        subgraph "Envoy Filter Chain (Wasm)"
            C -->|Request Body| D[Shield Wasm Filter]
            D -- Scans for PII/Credentials --> E{Sensitive Data Found?}
            E -- Yes --> F[Redact Data & Add Header]
            E -- No --> G[Allow Unmodified]
        end
        
        F -->|Redacted Request| I((api.openai.com))
        G -->|Clean Request| I
        F -.->|Log Event (gRPC/HTTP)| H[Shield Observer Backend]
        H --> J[CISO Dashboard]
    end
```

1.  **eBPF Interception:** Cilium's eBPF programs intercept all network traffic at the kernel level. When it sees an outbound request to a configured AI API endpoint, it redirects it to the Envoy proxy.
2.  **Wasm Injection:** We use a `CiliumClusterwideEnvoyConfig` CRD to dynamically inject our TinyGo-based Wasm filter into the Envoy proxy's filter chain. This is a native Cilium feature, requiring no sidecars.
3.  **Real-time Redaction:** The Wasm filter inspects the HTTP request body. It uses optimized regex and the Luhn algorithm to find and replace sensitive data with `[REDACTED]`. The operation is extremely fast (<1ms latency) and memory-safe due to the Wasm sandbox.
4.  **Audit Logging:** A `X-Cilium-Shield-Status: REDACTED` header is added to the request, and a log event is sent asynchronously to our Go-based `Shield Observer` for real-time visibility.
5.  **CISO Dashboard:** A Next.js-based dashboard provides a "CISO Command Center" view of all redaction events, showing which pods are attempting to leak data and to where.

## 💻 Development Commands

```bash
# Install all dependencies
make install

# Start local development (backend + frontend)
make dev

# Run tests
make test

# Build Wasm filter
make build-wasm

# Build control plane
make build-control-plane

# Build all Docker images
make docker-all

# Deploy to Kubernetes
make k8s-deploy

# Run demo
make demo

# Clean build artifacts
make clean
```

## 🏆 Why Cilium-Shield Wins: Hackathon Judging Criteria

### 1. Relevance to eBPF & Cilium (Native Extension)
*   **No Sidecars:** We directly leverage Cilium's per-node proxy architecture, making our solution far more efficient than traditional service meshes.
*   **CRD-Native:** Configuration is managed entirely through the `CiliumClusterwideEnvoyConfig` CRD, proving this is a true extension of the Cilium ecosystem.

### 2. Technical Depth (Engineering Excellence)
*   **Wasm & TinyGo:** Our filter is written in TinyGo, compiling to a tiny (~1MB) and highly performant Wasm module. We include detailed comments on memory management and performance considerations.
*   **High-Performance Redaction:** We use efficient regex and the Luhn algorithm for accurate credit card validation, going beyond simple pattern matching.
*   **Concurrent Go Backend:** The Go observer uses goroutines to ingest a high volume of log events without blocking or dropping data, ensuring a reliable audit trail.

### 3. Creativity (Solving a Modern Problem)
*   **Securing the AI Revolution:** We apply the battle-tested power of eBPF and Cilium to a new, critical problem: **GenAI Data Security**.
*   **Proactive Prevention:** Instead of just detecting leaks, we prevent them at the source, turning Cilium into a proactive AI-DLP Firewall.

### 4. Clarity & Presentation (Judge-Friendly)
*   **Clean Code & Architecture:** The codebase is well-documented, and the architecture is clearly explained.
*   **Compelling Demo:** Our demo script tells a clear story, showing the problem and the immediate, powerful solution.

## 🎬 Demo Video

*[Link to your 3-minute demo video or GIF will go here after you record it]*

## 🚀 Quick Start

### Local Development (2 minutes)

```bash
# Install dependencies
make install

# Start development environment
make dev
```

This starts:
- Backend API: http://localhost:3001
- Dashboard: http://localhost:3000

Test it:
```bash
make test
```

### Full Kubernetes Deployment

See [QUICKSTART.md](QUICKSTART.md) for detailed instructions.

## 📁 Project Structure

```
Cilium-Shield/
├── wasm-filter/          # TinyGo Wasm L7 filter
│   ├── main.go           # Wasm filter logic (Luhn algorithm, regex)
│   ├── main_test.go      # Unit tests
│   └── Dockerfile        # Multi-stage build for Wasm
├── control-plane/        # Go Observer backend
│   ├── server.go         # Event ingestion API
│   ├── server_test.go    # Tests
│   └── Dockerfile        # Go binary container
├── backend/              # Node.js Express API
│   ├── index.js          # REST API for dashboard
│   └── package.json
├── frontend/             # React Dashboard
│   ├── src/
│   │   └── components/
│   │       └── RedactionDashboard.jsx
│   └── package.json
├── ui/                   # UI component reference
│   └── components/
│       └── RedactionDashboard.jsx
├── k8s/                  # Kubernetes manifests
│   ├── cilium-envoy-config.yaml
│   ├── control-plane.yaml
│   ├── backend.yaml
│   ├── frontend.yaml
│   └── test-pod.yaml
├── scripts/              # Build and test scripts
│   ├── build-wasm.sh
│   ├── build-control-plane.sh
│   ├── test-redaction.sh
│   └── start-dev.sh
├── docs/                 # Documentation
│   ├── ARCHITECTURE.md   # Technical deep dive
│   ├── PRD.md            # Product requirements
│   ├── TESTING.md        # Testing guide
│   ├── HATA_GIDERME.md   # Troubleshooting
│   └── demo-script.md    # Video demo script
├── Makefile              # Convenience commands
├── QUICKSTART.md         # Quick start guide
└── README.md             # This file
```

## 🛠️ Getting Started

**Quick Start:** See [QUICKSTART.md](QUICKSTART.md) for a 10-minute setup guide.

**Detailed Testing:** See [TESTING.md](docs/TESTING.md) for comprehensive build, deploy, and test instructions.

**Architecture:** See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for technical deep dive.
