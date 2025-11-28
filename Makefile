.PHONY: help build-wasm build-control-plane build-all test dev clean install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	cd backend && npm install
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "All dependencies installed!"

build-wasm: ## Build the Wasm filter
	@bash scripts/build-wasm.sh

build-control-plane: ## Build the Go control plane
	@bash scripts/build-control-plane.sh

build-all: build-wasm build-control-plane ## Build all components

dev: ## Start development environment (backend + frontend)
	@bash scripts/start-dev.sh

test: ## Run redaction tests
	@bash scripts/test-redaction.sh

test-go: ## Run Go tests
	@echo "Running Wasm filter tests..."
	cd wasm-filter && go test -v ./...
	@echo "Running control plane tests..."
	cd control-plane && go test -v ./...

docker-wasm: ## Build Wasm Docker image
	cd wasm-filter && docker build -t cilium-shield-wasm:latest .

docker-control-plane: ## Build control plane Docker image
	cd control-plane && docker build -t cilium-shield-observer:latest .

docker-all: docker-wasm docker-control-plane ## Build all Docker images

clean: ## Clean build artifacts
	rm -f wasm-filter/filter.wasm
	rm -f wasm-filter/filter.wasm.sha256
	rm -f control-plane/server
	@echo "Clean complete!"

k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -f k8s/control-plane.yaml
	kubectl apply -f k8s/backend.yaml
	kubectl apply -f k8s/frontend.yaml
	@echo "Waiting for Cilium-Shield components to be ready..."
	kubectl wait --for=condition=ready pod -l app=cilium-shield-observer -n kube-system --timeout=60s
	kubectl wait --for=condition=ready pod -l app=cilium-shield-backend --timeout=60s
	kubectl wait --for=condition=ready pod -l app=cilium-shield-frontend --timeout=60s
	@echo "Deployment complete!"

k8s-delete: ## Delete from Kubernetes
	kubectl delete -f k8s/frontend.yaml || true
	kubectl delete -f k8s/backend.yaml || true
	kubectl delete -f k8s/control-plane.yaml || true
	kubectl delete -f k8s/test-pod.yaml || true

demo: ## Run the demo test
	@echo "Creating test pod..."
	kubectl apply -f k8s/test-pod.yaml
	@echo "Waiting for test pod to be ready..."
	kubectl wait --for=condition=ready pod/test-pod --timeout=60s
	@echo ""
	@echo "Sending test request with sensitive data..."
	kubectl exec test-pod -- curl -X POST -H "Content-Type: application/json" \
		-d '{"credit_card": "4532015112830366", "email": "admin@company.com", "api_key": "sk-proj-abcdefgh12345678"}' \
		http://api.openai.com/v1/completions
	@echo ""
	@echo "Check the dashboard at http://localhost:3000 to see the redaction event!"
