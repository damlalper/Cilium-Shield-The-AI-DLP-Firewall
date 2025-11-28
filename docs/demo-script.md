
# 3-Minute Demo Video Script: Cilium-Shield

**Objective:** Quickly and clearly demonstrate the problem, our unique solution, and the business value.

**Tone:** Professional, fast-paced, and technical.

---

### **[0:00-0:30] The Problem: "Shadow AI" is a Ticking Time Bomb**

**(Scene: A developer is shown in a code editor. They copy a large JSON object containing fake customer data, including a credit card number and an email address.)**

**Narrator:** "In today's enterprise, developers are racing to integrate AI. But this speed creates a massive risk: 'Shadow AI'. Sensitive data—customer PII, financial records—is silently leaking into public large language models."

**(Scene: The developer pastes the JSON into a web UI for an LLM and hits 'send'. We see a network traffic monitoring tool showing the raw, unencrypted data leaving the cluster.)**

**Narrator:** "Traditional firewalls are blind to this. They see encrypted traffic to a trusted endpoint like 'api.openai.com', but they have no idea that sensitive data is inside the payload. This is a compliance nightmare."

---

### **[0:30-1:30] The Solution: Cilium-Shield - eBPF & Wasm Powered DLP**

**(Scene: A terminal window appears. The user applies a single Kubernetes manifest with `kubectl apply -f k8s/cilium-envoy-config.yaml`.)**

**Narrator:** "This is where Cilium-Shield comes in. We don't use clunky, resource-heavy sidecars. With a single `CiliumEnvoyConfig` manifest, we dynamically inject a high-performance WebAssembly filter directly into Cilium's Envoy proxy."

**(Scene: An architecture diagram (from `ARCHITECTURE.md`) is shown on screen, highlighting the Wasm module running inside Envoy.)**

**Narrator:** "Our TinyGo-based Wasm module leverages eBPF's L7 visibility to inspect HTTP request bodies *before* they leave the pod. It's sandboxed, incredibly fast, and requires zero application code changes."

---

### **[1:30-2:30] Live Demonstration: Redaction in Action**

**(Scene: We're back to the developer's screen. They attempt to send the *exact same sensitive JSON* data again.)**

**Narrator:** "Let's replay that exact same request, but this time, with Cilium-Shield active."

**(Scene: We again see the network traffic monitor. This time, the JSON payload is shown, but the credit card number and API key are replaced with `[REDACTED]`.)**

**Narrator:** "Look. The data is intercepted in-flight. Our Wasm filter,. using a Luhn check for accuracy, has identified and redacted the credit card number and the API key. A custom header, `X-Cilium-Shield-Status: REDACTED`, is added for traceability."

**(Scene: The screen switches to the "CISO Command Center" Next.js dashboard. A new event immediately pops up in the table, showing the timestamp, source pod, redacted data type, and destination.)**

**Narrator:** "And critically, the security team has instant visibility. Our Go-based observer logs every single redaction event, providing a real-time audit trail of potential data leaks across the entire cluster."

---

### **[2:30-3:00] The "Why": Winning with Cilium**

**(Scene: A final slide with the project logo and key takeaways.)**

**Narrator:** "Cilium-Shield isn't just a tool; it's a new approach to cloud-native security. By combining the kernel-level power of eBPF with the flexibility of WebAssembly, we provide granular, payload-aware security that is:
* **Performant:** No sidecars means less overhead.
* **Secure:** Sandboxed and memory-safe.
* **Transparent:** No developer friction.

Cilium-Shield allows businesses to embrace AI, confident that their most valuable data is protected. Thank you."
