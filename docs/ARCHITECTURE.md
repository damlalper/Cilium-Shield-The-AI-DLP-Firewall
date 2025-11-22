### 2. Teknik Derinlik Dosyası: `ARCHITECTURE.md`
*Bu dosya senin "Senior Engineer" olduğunu kanıtlayan yerdir. Jürideki teknik kişiler (Liz Rice vb.) buraya bayılır.*

# 🏗️ System Architecture

## Overview
Cilium-Shield leverages the **Cilium EnvoyConfig CRD** to inject custom **WebAssembly (Wasm)** filters into the sidecar-free Service Mesh. This allows us to perform Layer 7 inspection on outbound traffic without modifying the application code.

## High-Level Diagram

```mermaid
graph LR
    subgraph Kubernetes Cluster
        A[App Pod] -->|HTTP Req| B(Cilium Node Agent)
        B -->|Redirect| C{Envoy Proxy}
        
        subgraph "Envoy Processing"
            C -->|Input Stream| D[Wasm Module (TinyGo)]
            D -->|Logic| E{Contains PII?}
            E -- Yes --> F[Redact Data & Log Event]
            E -- No --> G[Pass Through]
        end
        
        F -.->|gRPC/HTTP| H[Shield Observer Backend]
    end
    
    G -->|Sanitized Req| I((OpenAI API))
    F -->|Redacted Req| I
    H --> J[Dashboard UI]
```

## Core Components
1. The Wasm Filter (Data Plane)
Language: Written in TinyGo using the proxy-wasm-go-sdk.

Function: It hooks into the OnHttpRequestBody phase.

Performance: Because it runs inside the Envoy sandbox, it incurs negligible latency (<1ms) compared to external sidecars.

Detection Logic: Uses high-performance Regex to scan for:

Credit Card Numbers (Luhn Algorithm check)

OpenAI Secret Keys (sk-...)

Email Addresses

 2. Cilium Integration
We use the CiliumEnvoyConfig Custom Resource Definition (CRD). This tells Cilium to:

Locate the listener for egress traffic on port 80/443.

Download the Wasm module from the container registry.

Inject it into the filter chain for targeted pods.

3. The Observer (Control Plane)
Stack: Golang + Redis.

Role: Receives asynchronous logs from the Wasm filter via header metadata extraction. It stores "Leak Events" for auditing without blocking the main traffic flow.

## Design Decisions: Why Wasm + Cilium?
Why not a Sidecar? Running a full sidecar (like classic Istio) for every pod consumes too much RAM. Cilium's per-node proxy model is more efficient.

Why not pure eBPF? While eBPF is great for L3/L4, parsing complex JSON bodies in kernel space is restricted (verifier limits) and hard to maintain. Wasm in userspace (Envoy) is the perfect compromise for L7 tasks.