# Kubernetes Manifests

This directory contains all Kubernetes manifests needed to deploy Cilium-Shield.

## Files

- `cilium-envoy-config.yaml` - CiliumClusterwideEnvoyConfig that injects the Wasm filter
- `control-plane.yaml` - Go-based Shield Observer control plane deployment
- `backend.yaml` - Node.js Express backend for the dashboard API
- `frontend.yaml` - React frontend dashboard deployment
- `test-pod.yaml` - Test pod for sending requests with sensitive data

## Deployment Order

1. **Deploy Control Plane** (Go Observer):
   ```bash
   kubectl apply -f control-plane.yaml
   ```

2. **Deploy Backend**:
   ```bash
   kubectl apply -f backend.yaml
   ```

3. **Deploy Frontend**:
   ```bash
   kubectl apply -f frontend.yaml
   ```

4. **Deploy Cilium Envoy Config** (after building and pushing Wasm module):
   ```bash
   # First update the SHA256 hash in cilium-envoy-config.yaml
   kubectl apply -f cilium-envoy-config.yaml
   ```

5. **Deploy Test Pod**:
   ```bash
   kubectl apply -f test-pod.yaml
   ```

## Testing

After deployment, test the redaction:

```bash
# Enter the test pod
kubectl exec -it test-pod -- sh

# Send a request with sensitive data
curl -X POST -H "Content-Type: application/json" \
  -d '{"credit_card": "4532015112830366", "email": "test@example.com"}' \
  http://api.openai.com/v1/completions

# Check the control plane logs
kubectl logs -n kube-system -l app=cilium-shield-observer

# Check the backend logs
kubectl logs -l app=cilium-shield-backend
```

## Accessing the Dashboard

```bash
# Get the frontend service URL
kubectl get svc cilium-shield-frontend

# Port forward if using ClusterIP
kubectl port-forward svc/cilium-shield-frontend 3000:80
```

Then open http://localhost:3000 in your browser.

## Architecture

```
┌─────────────────────────────────────────────────────┐
│ Kubernetes Cluster                                   │
│                                                       │
│  ┌─────────────┐    ┌──────────────────┐            │
│  │  Test Pod   │───>│ Cilium + Envoy   │            │
│  │             │    │  + Wasm Filter   │            │
│  └─────────────┘    └────────┬─────────┘            │
│                              │                       │
│                              │ Redacted Traffic      │
│                              v                       │
│                     ┌─────────────────┐              │
│                     │  External API   │              │
│                     │ (api.openai.com)│              │
│                     └─────────────────┘              │
│                                                       │
│  ┌──────────────┐   ┌──────────────┐                │
│  │   Control    │<──│   Backend    │<───┐           │
│  │    Plane     │   │   (Node.js)  │    │           │
│  │   (Go)       │   └──────────────┘    │           │
│  └──────────────┘                       │           │
│                                         │           │
│                     ┌──────────────┐    │           │
│                     │   Frontend   │────┘           │
│                     │   (React)    │                │
│                     └──────────────┘                │
└─────────────────────────────────────────────────────┘
```
