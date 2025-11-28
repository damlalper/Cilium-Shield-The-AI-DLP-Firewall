# Cilium-Shield Quick Start Guide

This guide will help you get Cilium-Shield up and running in under 10 minutes.

## Prerequisites

- Kubernetes cluster with Cilium installed
- `kubectl` configured
- `docker` installed
- `tinygo` installed (for Wasm compilation)
- `node` and `npm` installed (for backend/frontend)
- `make` installed (optional, for convenience)

## Option 1: Local Development (Fastest)

### 1. Install Dependencies

```bash
make install
# or manually:
cd backend && npm install
cd ../frontend && npm install
```

### 2. Start Dev Environment

```bash
make dev
# or manually:
bash scripts/start-dev.sh
```

This will start:
- Backend API on http://localhost:3001
- Frontend Dashboard on http://localhost:3000

### 3. Test the API

```bash
# In a new terminal
make test
# or manually:
bash scripts/test-redaction.sh
```

### 4. View the Dashboard

Open http://localhost:3000 in your browser to see the CISO Command Center.

---

## Option 2: Kubernetes Deployment (Full Demo)

### 1. Build the Wasm Filter

```bash
make build-wasm
# or manually:
cd wasm-filter
tinygo build -o filter.wasm -target=wasi main.go
```

### 2. Build the Control Plane

```bash
make build-control-plane
# or manually:
cd control-plane
CGO_ENABLED=0 GOOS=linux go build -o server .
```

### 3. Build Docker Images

```bash
# Build Wasm image
cd wasm-filter
docker build -t your-registry/cilium-shield-wasm:v1 .
docker push your-registry/cilium-shield-wasm:v1

# Build control plane image
cd ../control-plane
docker build -t your-registry/cilium-shield-observer:v1 .
docker push your-registry/cilium-shield-observer:v1
```

### 4. Update Kubernetes Manifests

Update the following files with your registry addresses:
- `k8s/control-plane.yaml` - Change `registry.example.com` to your registry
- `k8s/cilium-envoy-config.yaml` - Update the Wasm module URI and SHA256 hash

Get the SHA256 hash:
```bash
sha256sum wasm-filter/filter.wasm
```

### 5. Deploy to Kubernetes

```bash
make k8s-deploy
# or manually:
kubectl apply -f k8s/control-plane.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/cilium-envoy-config.yaml
```

### 6. Test the Deployment

```bash
make demo
# This will:
# 1. Create a test pod
# 2. Send a request with sensitive data
# 3. Show the redaction in action
```

### 7. Access the Dashboard

```bash
# Get the frontend service
kubectl get svc cilium-shield-frontend

# Port forward if needed
kubectl port-forward svc/cilium-shield-frontend 3000:80
```

Then open http://localhost:3000

---

## Troubleshooting

### Backend not starting

```bash
# Check backend logs
cd backend
npm install
npm start
```

### Frontend not showing data

1. Check if backend is running: `curl http://localhost:3001`
2. Check browser console for errors
3. Verify `REACT_APP_API_URL` environment variable

### Wasm compilation fails

```bash
# Install tinygo
# macOS: brew install tinygo
# Linux: see https://tinygo.org/getting-started/install/

# Verify installation
tinygo version
```

### Kubernetes pods not starting

```bash
# Check pod status
kubectl get pods -A

# Check pod logs
kubectl logs -n kube-system -l app=cilium-shield-observer
kubectl logs -l app=cilium-shield-backend
```

---

## Next Steps

1. Read [ARCHITECTURE.md](docs/ARCHITECTURE.md) to understand how it works
2. Check [TESTING.md](docs/TESTING.md) for detailed testing instructions
3. Watch the [demo video](#) to see it in action
4. Customize the redaction rules in `wasm-filter/main.go`

## Support

For issues and questions:
- Check [HATA_GIDERME.md](docs/HATA_GIDERME.md) for common problems
- Open an issue on GitHub
- Read the [PRD.md](docs/PRD.md) for product details
