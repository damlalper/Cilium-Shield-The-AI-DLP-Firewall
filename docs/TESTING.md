
# Testing Cilium-Shield

This document outlines the steps to test the Cilium-Shield MVP.

## 1. Prerequisites

*   A Kubernetes cluster with Cilium installed and enabled.
*   `tinygo` installed to compile the Wasm module.
*   `docker` installed to build and push the Wasm module container image.
*   `kubectl` installed and configured to access your Kubernetes cluster.

## 2. Compile the Wasm Module

1.  Navigate to the `wasm-filter` directory:
    ```sh
    cd wasm-filter
    ```
2.  Compile the Go code to a Wasm module:
    ```sh
    tinygo build -o filter.wasm -target wasm ./main.go
    ```
    This will create a `filter.wasm` file in the `wasm-filter` directory.

## 3. Build and Push the Wasm Module Container Image

1.  Create a `Dockerfile` in the `wasm-filter` directory:
    ```dockerfile
    FROM scratch
    COPY filter.wasm /
    ```
2.  Build the container image:
    ```sh
    docker build -t registry.example.com/cilium-shield-wasm:v1 .
    ```
3.  Push the container image to your container registry:
    ```sh
    docker push registry.example.com/cilium-shield-wasm:v1
    ```
    **Note:** Replace `registry.example.com` with your actual container registry address.

## 4. Deploy the Control Plane

1.  Navigate to the `control-plane` directory:
    ```sh
    cd control-plane
    ```
2.  Build the Go application:
    ```sh
    go build -o server .
    ```
3.  Create a `Dockerfile` for the control plane:
    ```dockerfile
    FROM scratch
    COPY server /
    CMD ["/server"]
    ```
4.  Build and push the container image:
    ```sh
    docker build -t registry.example.com/cilium-shield-observer:v1 .
    docker push registry.example.com/cilium-shield-observer:v1
    ```
5.  Create a Kubernetes deployment for the control plane. Create a `k8s/control-plane.yaml` file with the following content:
    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: cilium-shield-observer
      namespace: kube-system
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: cilium-shield-observer
      template:
        metadata:
          labels:
            app: cilium-shield-observer
        spec:
          containers:
            - name: observer
              image: registry.example.com/cilium-shield-observer:v1
              ports:
                - containerPort: 8080
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: cilium-shield-observer
      namespace: kube-system
    spec:
      selector:
        app: cilium-shield-observer
      ports:
        - protocol: TCP
          port: 80
          targetPort: 8080
    ```
6.  Apply the deployment:
    ```sh
    kubectl apply -f k8s/control-plane.yaml
    ```

## 5. Apply the Cilium CRD

1.  Update the `k8s/cilium-envoy-config.yaml` file with the correct SHA256 hash of your Wasm module. You can get the hash by running:
    ```sh
    sha256sum wasm-filter/filter.wasm
    ```
2.  Apply the Cilium CRD:
    ```sh
    kubectl apply -f k8s/cilium-envoy-config.yaml
    ```

## 6. Test the Redaction

1.  Create a test pod:
    ```sh
    kubectl run -it --rm --image=curlimages/curl:latest test-pod -- sh
    ```
2.  From the test pod's shell, send a request to `api.openai.com` with sensitive data:
    ```sh
    curl -X POST -H "Content-Type: application/json" -d '{"credit_card": "49927398716"}' http://api.openai.com/v1/completions
    ```
3.  Check the logs of the `cilium-shield-observer` pod to see the redaction event:
    ```sh
    kubectl logs -n kube-system -l app=cilium-shield-observer
    ```
    You should see a log message similar to this:
    ```
    Received event: Source: <test-pod-ip>, Destination: http://api.openai.com/v1/completions, Type: REDACTED_CREDIT_CARD
    ```
4.  You can also check the response headers of the request to see the `X-Cilium-Shield-Status` header.
