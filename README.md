# 🛡️ Cilium-Shield: The AI-DLP Firewall
> **Zero-Trust Data Leakage Prevention for the GenAI Era.**
> Powered by eBPF, Cilium Service Mesh & WebAssembly.

![GitHub License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)
![Tech Stack](https://img.shields.io/badge/Stack-Cilium%20%7C%20Wasm%20%7C%20Go-purple)
![Hackathon](https://img.shields.io/badge/Submission-eBPF%20Summit%202025-orange)

## 📺 [Watch the Demo Video (3 min)](LINK_TO_YOUR_YOUTUBE_VIDEO)

---

## 🚀 The Problem: Shadow AI
Enterprises are racing to adopt GenAI, but developers often inadvertently send sensitive customer data (PII, API Keys, Credit Cards) to public LLMs like OpenAI or Anthropic. Traditional firewalls work at Layer 3/4 and cannot inspect encrypted JSON payloads effectively without heavy sidecars.

## 💡 The Solution
**Cilium-Shield** is a high-performance, Kubernetes-native DLP system. It injects **WebAssembly (Wasm)** filters into the **Cilium Service Mesh (Envoy)** to inspect and sanitize outbound LLM traffic in real-time.

### ✨ Key Features
* **🔎 Deep Packet Inspection (L7):** Inspects HTTP/JSON bodies leaving the cluster.
* **⚡ Ultra-Low Latency:** Runs as a Wasm module inside Envoy (no context switching).
* **🔒 Real-time Redaction:** Automatically masks PII (e.g., `sk-proj-***`, `4444-****`) before it leaves the pod.
* **🚨 Security Dashboard:** A visual command center to audit blocked attempts.

---

## 🛠️ Tech Stack
* **Data Plane:** [Cilium Service Mesh](https://cilium.io/) + Envoy Proxy
* **Filter Logic:** WebAssembly (TinyGo)
* **Control Plane:** Golang
* **Frontend:** Next.js + Tailwind CSS

## 🏗️ Architecture Snapshot
*(See [ARCHITECTURE.md](ARCHITECTURE.md) for deep dive)*

The system sits between your Microservices and the Public Internet, intercepting traffic via Cilium's Envoy integration.

---

## ⚡ Quick Start

### Prerequisites
* Kubernetes Cluster (Kind/Minikube or Cloud)
* Cilium CLI installed & Hubble enabled

### Installation
1. **Deploy the Wasm Filter Policy:**
```bash
   kubectl apply -f k8s/cilium-envoy-config.yaml
   ```

2. Deploy the Observer Backend:

```bash
kubectl apply -f k8s/backend-deployment.yaml
```

3. Access the Dashboard:

```bash
kubectl port-forward svc/cilium-shield-ui 3000:3000
```
Open: http://localhost:3000


## 🏆 Why This Matters for the eBPF Ecosystem
This project demonstrates how Cilium can be extended beyond networking into the Application Security domain using Wasm, bridging the gap between platform engineering and SecOps.