#!/bin/bash
set -e

echo "Building Cilium-Shield Control Plane..."

cd control-plane

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go to compile the control plane."
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

# Build the Go application
echo "Building Go server..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

if [ -f server ]; then
    echo "✓ Control plane server compiled successfully"

    # Get file size
    SIZE=$(ls -lh server | awk '{print $5}')
    echo "Size: $SIZE"
else
    echo "✗ Failed to compile control plane server"
    exit 1
fi

echo ""
echo "Next steps:"
echo "1. Build the Docker image: docker build -t your-registry/cilium-shield-observer:v1 ."
echo "2. Push to registry: docker push your-registry/cilium-shield-observer:v1"
echo "3. Update k8s/control-plane.yaml with your registry address"
