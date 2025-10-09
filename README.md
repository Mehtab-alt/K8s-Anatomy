# K8s-Anatomy
K8s-Anatomy is a project that models a microservices system as a living organism. Kubernetes forms the skeleton, while the Prometheus stack acts as a central nervous system that enables observation, pain responses (alerting), and autonomic reflexes (autoscaling) for a resilient, self-regulating system.

### Project Philosophy

This project is designed as a clear, step-by-step tutorial. To make the core concepts of monitoring as clear as possible, our organs are initially **decoupled**. This allows you to observe the vital signs of each microservice in isolation, which is the foundational skill for observability.

-   **Organs (`/organs`):** Our microservices are the vital organs. The `heart-service` provides a measurable "pulse" (request rate), while the `brain-service` performs a "cognitive task" with measurable latency. They operate independently so you can learn to monitor each one's unique health signals without interference.
-   **Skeletal System (`/skeletal-system`):** Kubernetes provides the structure. Its manifests define the organism's anatomy, giving the organs a place to live, connect, and scale. This includes an **autonomic response system** (Horizontal Pod Autoscaler) to react to stress.
-   **Central Nervous System (`/nervous-system`):** Prometheus and Grafana form the brain and nerves. This system listens to the signals from the organs, processes them, and provides a window into the organism's overall consciousness and well-being. It can also trigger a **pain response** (Alerting) when vitals are abnormal.

---

### Prerequisites: The Spark of Life

To begin, you will need the following tools installed and configured:
*   `kubectl` connected to a running Kubernetes cluster (e.g., Minikube, Kind, Docker Desktop, or a cloud provider).
*   `helm` for installing the nervous system.
*   `go` (v1.21+) for building the organs locally.
*   `docker` for synthesizing the organs (if not relying on pre-built images).

---

## Structure

organism-k8s/
├── .github/
│   └── workflows/
│       └── ci-build-push.yaml
├── docs/
│   └── images/
│       └── dashboard-screenshot.png  # Placeholder for the visual payoff
├── nervous-system/
│   ├── prometheus/
│   │   ├── organism-alerts.yaml
│   │   └── service-monitor.yaml
│   └── grafana/
│       └── dashboards/
│           └── organism-health.json
├── organs/
│   ├── heart-service/
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── main.go
│   └── brain-service/
│       ├── Dockerfile
│       ├── go.mod
│       ├── go.sum
│       └── main.go
├── skeletal-system/
│   ├── base/
│   │   ├── brain-service-deployment.yaml
│   │   ├── brain-service-service.yaml
│   │   ├── heart-service-deployment.yaml
│   │   ├── heart-service-hpa.yaml
│   │   ├── heart-service-service.yaml
│   │   └── prometheus-rbac.yaml
│   └── overlays/
│       └── production/
│           ├── kustomization.yaml
│           └── patch-replicas.yaml
└── README.md

---

## A Step-by-Step Guide to Life

Follow these steps to assemble and observe your organism.

### Step 1: Synthesizing the Organs

The organs are container images that run our services. This repository is configured with a GitHub Action (`.github/workflows/ci-build-push.yaml`) that automatically builds and pushes these images to a container registry when you push to `main`.

**Important:** Before you begin, you must have a container registry (like Docker Hub or GHCR) and configure the `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets in your GitHub repository settings for the Action to work.

**Alternative: Manual Synthesis**
If you wish to build the organs yourself, run the following commands from the root of the repository:
```bash
# Build the heart-service
docker build -t your-repo/organism-heart-service:latest ./organs/heart-service

# Build the brain-service
docker build -t your-repo/organism-brain-service:latest ./organs/brain-service

# Don't forget to push them to your registry!
# docker push your-repo/organism-heart-service:latest
# docker push your-repo/organism-brain-service:latest
```
***NOTE: You MUST update the image paths in `skeletal-system/base/heart-service-deployment.yaml` and `skeletal-system/base/brain-service-deployment.yaml` to point to your image repository.***

### Step 2: Activating the Central Nervous System

Now, we implant the nervous system using the `kube-prometheus-stack` Helm chart. This will install Prometheus (the brain), Grafana (the visual cortex), and Alertmanager (the pain center).

> **Why this chart?** The `kube-prometheus-stack` is the de-facto community standard for Kubernetes monitoring. It's a "batteries-included" package that not only installs the core tools but also pre-configures them to work together seamlessly. Crucially, it creates the Kubernetes Custom Resource Definitions (CRDs) like `ServiceMonitor` and `PrometheusRule` that allow us to define our monitoring and alerting declaratively, just like any other Kubernetes object. This saves immense manual configuration time.

```bash
# 1. Add the Prometheus community Helm repository
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# 2. Create a dedicated namespace for the nervous system
kubectl create namespace monitoring

# 3. Install the stack. We give it a release name 'prometheus-stack' for consistency.
helm install prometheus-stack prometheus-community/kube-prometheus-stack --namespace monitoring
```
It may take a few minutes for all components to become fully operational.

### Step 3: Assembling the Organism

With the nervous system in place, we can now build the organism's body using our Kubernetes manifests. The `kustomization` overlay will configure it for a production-like environment with scaled-up organs.

```bash
# This command assembles the deployments and services, giving our organs a home.
kubectl apply -k skeletal-system/overlays/production
```

### Step 4: Granting the Nervous System Vision

By default, the nervous system in the `monitoring` namespace cannot "see" the organs in the `default` namespace. We must apply RBAC rules to grant it permission to observe other parts of the body.

```bash
# This command applies the Role and RoleBinding needed for cross-namespace communication.
kubectl apply -f skeletal-system/base/prometheus-rbac.yaml
```

### Step 5: Connecting the Nerves

The brain is active and the organs are in place, but they aren't connected. The `ServiceMonitor` manifest acts as the nerve endings, telling Prometheus exactly how to listen to the vital signs from our organs.

```bash
# This command connects the brain to the organs.
kubectl apply -f nervous-system/prometheus/service-monitor.yaml
```

### Step 6: Teaching the Organism to Feel Pain (Alerting)

A resilient system must react to danger. We will now apply rules that teach the nervous system to fire alerts when vital signs are critical.

```bash
# This command creates PrometheusRule objects that Alertmanager will use.
kubectl apply -f nervous-system/prometheus/organism-alerts.yaml
```

### Step 7: Enabling the Fight-or-Flight Response (Autoscaling)

Finally, we give the organism an adrenal gland. The Horizontal Pod Autoscaler will automatically add more "heart" capacity when the system is under heavy load, allowing it to adapt to stress.

```bash
# This command creates the HPA for the heart-service.
kubectl apply -f skeletal-system/base/heart-service-hpa.yaml
```

### Step 8: Observing Consciousness

The organism is now alive and fully wired. Let's access its consciousness through the Grafana dashboard.

**1. Access Grafana:**
Open a new terminal and forward the Grafana port to your local machine.
```bash
kubectl port-forward svc/prometheus-stack-grafana 3000:80 -n monitoring
```
Now, open your web browser and navigate to `http://localhost:3000`.

**2. Log In:**
*   **Username:** `admin`
*   **Password:** Run this command to retrieve the auto-generated password:
    ```bash
    kubectl get secret prometheus-stack-grafana -n monitoring -o jsonpath="{.data.admin-password}" | base64 --decode; echo
    ```

**3. Import the Dashboard:**
*   In the Grafana UI, navigate to the **Dashboards** section.
*   Click **"New"** and then **"Import"**.
*   Use the "Upload JSON file" button to upload `nervous-system/grafana/dashboards/organism-health.json`.
*   Click **"Load"**, select your Prometheus data source from the dropdown, and then **"Import"**.

You should now see the **"Organism Health Status"** dashboard!

![Organism Health Dashboard](docs/images/dashboard-screenshot.png "A screenshot of the Grafana dashboard showing the organism's vital signs.")
*(Note: Replace the placeholder image at `docs/images/dashboard-screenshot.png` with a real screenshot of your dashboard in action!)*

**4. Generate a Pulse and Observe (In Isolation):**
The organism is idle. Let's stimulate each organ independently to see its unique response.

**A. Stimulate the Heart:**
First, port-forward the `heart-service`:
```bash
# In a new terminal:
kubectl port-forward svc/heart-service 8080:8080
```
Now, send it traffic using the `hey` load testing tool:
```bash
# If you don't have hey: go install go.uber.org/hey
# This sends 5 requests per second for 5 minutes.
hey -z 5m -q 5 http://localhost:8080/beat
```
**Observe on the dashboard:**
*   The **"System Pulse"** panel will rise dramatically.
*   The **"CPU Respiration"** for the `heart-service` pods will increase.
*   **The Adrenaline Response:** After a minute, run `kubectl get hpa` and `kubectl get pods -w`. You'll see Kubernetes automatically scale up new `heart-service` pods to handle the load!
*   The **"Cognitive Processing Time"** panel will remain flat.

**B. Stimulate the Brain:**
Stop the previous test. Now, port-forward the `brain-service`:
```bash
# In a new terminal:
kubectl port-forward svc/brain-service 8081:8080
```
Now, send it traffic:
```bash
# This sends 2 requests per second for 2 minutes.
hey -z 2m -q 2 http://localhost:8081/think
```
**Observe on the dashboard:**
*   The **"Cognitive Processing Time"** heatmap will now populate, showing the latency of the brain's "thoughts".
*   The **"CPU Respiration"** for the `brain-service` pods will increase.
*   The **"System Pulse"** panel will remain flat.

You have successfully monitored the vital signs of individual components within a larger system.

---

### Evolution: From Independent Organs to a True System

This project intentionally keeps the organs decoupled to teach the fundamentals of monitoring in isolation. The next step in your journey is to connect them.

A simple approach is to have the `brain-service` make a direct, synchronous HTTP call to the `heart-service`. However, this introduces significant risks:

1.  **Cascading Failures:** If the `heart-service` slows down or fails, the `brain-service` will also fail, potentially causing a system-wide outage. This is a fragile architectural pattern.
2.  **Ambiguous Metrics:** The `brain-service`'s latency metric would no longer measure just its own work; it would measure `(brain work + network time + heart work)`. If latency spikes, how do you know which organ is the cause?

**The Professional Solution:**
To build a truly resilient and observable interconnected system, you must evolve your organism with more advanced capabilities:

*   **Resilience Patterns (The Reflexes):** Implement patterns like **Circuit Breakers** (e.g., using `gobreaker`) in the `brain-service`. If the heart becomes unresponsive, the circuit breaker "trips" and stops calls, preventing the brain from failing and giving the heart time to recover.
*   **Distributed Tracing (The Circulatory System):** This is the correct tool for understanding inter-service calls. Technologies like **OpenTelemetry** and backends like **Jaeger** allow you to trace a single request as it flows between organs. This lets you see exactly how much time was spent in the brain, on the network, and in the heart, pinpointing the true source of latency.

Mastering these advanced concepts is the key to evolving your simple organism into a production-grade, resilient distributed system.
