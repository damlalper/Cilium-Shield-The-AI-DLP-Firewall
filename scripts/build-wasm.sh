#!/bin/bash
set -e

echo "Building Cilium-Shield Wasm Filter..."

cd wasm-filter

# Check if tinygo is installed
if ! command -v tinygo &> /dev/null; then
    echo "Error: tinygo is not installed. Please install tinygo to compile the Wasm module."
    echo "Visit: https://tinygo.org/getting-started/install/"
    exit 1
fi

# Build the Wasm module
echo "Compiling Go to Wasm using TinyGo..."
tinygo build -o filter.wasm -target=wasi main.go

if [ -f filter.wasm ]; then
    echo "✓ Wasm module compiled successfully: filter.wasm"

    # Calculate SHA256 hash
    if command -v sha256sum &> /dev/null; then
        SHA256=$(sha256sum filter.wasm | awk '{print $1}')
        echo "SHA256: $SHA256"
        echo "$SHA256" > filter.wasm.sha256
    elif command -v shasum &> /dev/null; then
        SHA256=$(shasum -a 256 filter.wasm | awk '{print $1}')
        echo "SHA256: $SHA256"
        echo "$SHA256" > filter.wasm.sha256
    fi

    # Get file size
    SIZE=$(ls -lh filter.wasm | awk '{print $5}')
    echo "Size: $SIZE"
else
    echo "✗ Failed to compile Wasm module"
    exit 1
fi

echo ""
echo "Next steps:"
echo "1. Build the Docker image: docker build -t your-registry/cilium-shield-wasm:v1 ."
echo "2. Push to registry: docker push your-registry/cilium-shield-wasm:v1"
echo "3. Update k8s/cilium-envoy-config.yaml with the SHA256 hash"
