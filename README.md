# K8s-Anatomy
K8s-Anatomy is a project that models a microservices system as a living organism. Kubernetes forms the skeleton, while the Prometheus stack acts as a central nervous system that enables observation, pain responses (alerting), and autonomic reflexes (autoscaling) for a resilient, self-regulating system.

### The Philosophy

-   **Organs (`/organs`):** Our microservices (`heart-service`, `brain-service`) are the vital organs. They perform core functions and, crucially, are designed to emit their own health signals (metrics).
-   **Skeletal System (`/skeletal-system`):** Kubernetes provides the structure. Its manifests define the organism's anatomy, giving the organs a place to live, connect, and scale. This includes an **autonomic response system** (Horizontal Pod Autoscaler) to react to stress.
-   **Central Nervous System (`/nervous-system`):** Prometheus and Grafana form the brain and nerves. This system listens to the signals from the organs, processes them, and provides a window into the organism's overall consciousness and well-being. It can also trigger a **pain response** (Alerting) when vitals are abnormal.

This project will guide you through the process of bringing this advanced, resilient organism to life in your own Kubernetes cluster.

---

### Prerequisites: The Spark of Life

To begin, you will need the following tools installed and configured:
*   `kubectl` connected to a running Kubernetes cluster (e.g., Minikube, Kind, Docker Desktop, or a cloud provider).
*   `helm` for installing the nervous system.
*   `docker` for synthesizing the organs (if not relying on pre-built images).

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
This deploys the `heart-service` and `brain-service` into the `default` namespace.

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

**4. Generate a Pulse and Observe:**
The organism is idle. Let's give it something to do. Open another terminal and run a loop to send requests to the `heart-service`.

First, port-forward the `heart-service`:
```bash
kubectl port-forward svc/heart-service 8080:8080
```

Now, send traffic. We'll use the `hey` load testing tool for a stronger pulse:
```bash
# If you don't have hey: go install go.uber.org/hey
# This sends 5 requests per second for 5 minutes.
hey -z 5m -q 5 http://localhost:8080/beat
```

**Watch the dashboard come alive!**
*   The **"System Pulse"** panel will rise.
*   The **"CPU Respiration"** for the `heart-service` pods will increase.
*   **The Adrenaline Response:** After a minute, run `kubectl get hpa`. You'll see the HPA targeting ~80% CPU. Run `kubectl get pods -w` and watch as Kubernetes automatically creates new `heart-service` pods to handle the load!
*   **The Pain Response:** The "High Pulse Rate" alert we defined fires above 10 req/s. Port-forward Alertmanager (`kubectl port-forward svc/prometheus-stack-alertmanager 9093:9093 -n monitoring`) and visit `http://localhost:9093` to see the active alert.

You have successfully created and observed a living, resilient system on Kubernetes.

---

### Future Evolution: The Circulatory System (Distributed Tracing)

This project masters one of the three pillars of observability: **Metrics**. To achieve a truly holistic view, you would also implement:

*   **Logs (Memory):** Using a stack like Loki to aggregate logs from all organs, providing a searchable history of events.
*   **Traces (Circulatory System):** Using a technology like OpenTelemetry and a backend like Jaeger or Tempo, you can trace the path of a single request as it flows between organs. This is like tracking a blood cell through the body, showing exactly where time is spent and revealing bottlenecks that metrics alone cannot. This is the next step in evolving your organism's self-awareness.
