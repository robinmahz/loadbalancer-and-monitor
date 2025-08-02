# Go Load Balancer: Comparative Analysis of Load Balancing Algorithms

This project implements and compares three load balancing algorithms—Weighted Round-Robin (WRR), Round-Robin (RR), and Least Connections—using Go. The load balancer distributes HTTP requests across three backend servers with varying capacities (small, medium, large). Prometheus and Grafana are used for monitoring and performance analysis. All services are containerized with Docker and orchestrated using Docker Compose. A detailed report analyzes the performance of the implemented algorithms.

## Project Structure

go-load-balancer/
│
├── FinalProjectReportToPrint.pdf
│
├── backend/
│ ├── Dockerfile
│ └── main.go
│
├── docker-compose.yml
│
├── grafana/
│ └── datasource.yml
│
├── loadbalancer/
│ ├── Dockerfile
│ └── main.go
│
└── prometheus/
└── prometheus.yml

## Features

- **Load Balancer**: Implements three algorithms:
  - **Round-Robin (RR)**: Distributes requests equally across backends in a cyclic order.
  - **Weighted Round-Robin (WRR)**: Distributes requests based on backend capacities (small, medium, large).
  - **Least Connections**: Routes requests to the backend with the fewest active connections.
- **Backend Services**: Three Go-based HTTP servers:
  - `backend1`: Small capacity, port 8081.
  - `backend2`: Medium capacity, port 8082.
  - `backend3`: Large capacity, port 8083.
- **Monitoring**: Prometheus collects metrics, and Grafana visualizes performance for algorithm comparison.
- **Containerized**: Services run in Docker containers, connected via a bridge network (`app-network`).
- **Configurable Resources**: CPU and memory limits are defined in `docker-compose.yml`.

## Prerequisites

- Docker (with Docker Compose)
- Go (optional, for development)
- Basic understanding of HTTP, load balancing, and containerization

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-load-balancer

2. Configure Container ResourcesThe docker-compose.yml defines three backend services with distinct configurations:backend1: Small capacity, 1 CPU core, 512MB memory limit, 256MB reserved, port 8081.
   backend2: Medium capacity, 1 CPU core, 512MB memory limit, 256MB reserved, port 8082.
   backend3: Large capacity, 1 CPU core, 512MB memory limit, 256MB reserved, port 8083.
   loadbalancer: Runs on port 8000, depends on all backends.
   prometheus: Runs on port 9090, scrapes metrics from all services.
   grafana: Runs on port 3000, visualizes Prometheus metrics.

To modify resource limits or environment variables (e.g., CAPACITY, SERVER_NAME, PORT), edit docker-compose.yml. Example:yaml

services:
backend1:
build: ./backend
environment: - SERVER_NAME=backend1 - CAPACITY=small - PORT=8081
deploy:
resources:
limits:
cpus: "1" # Adjust CPU limit
memory: 512M # Adjust memory limit
reservations:
memory: 256M

3. Configure Load Balancing AlgorithmsThe load balancer (loadbalancer/main.go) supports RR, WRR, and Least Connections. To switch algorithms, modify the configuration in main.go or set environment variables (if implemented). The CAPACITY environment variable (small, medium, large) influences WRR weights. Refer to FinalProjectReportToPrint.pdf for algorithm details and performance analysis.4. Build and RunStart all services using Docker Compose:bash

docker-compose up --build

This will:Build and start the load balancer, three backend services, Prometheus, and Grafana.
Expose the load balancer on http://localhost:8000.
Expose Grafana on http://localhost:3000 (default credentials: admin/admin).
Expose Prometheus on http://localhost:9090.

5. Access the ServicesLoad Balancer: Send HTTP requests to http://localhost:8000 to test RR, WRR, or Least Connections.
   Grafana Dashboard: Open http://localhost:3000 to visualize metrics (e.g., request latency, connection counts).
   Prometheus: Access http://localhost:9090 for raw metrics.
   Backend Services: Not directly exposed; accessed via the load balancer.

6. Stop the ServicesStop and remove the containers:bash

docker-compose down

DevelopmentTo modify the load balancer or backend services:Edit loadbalancer/main.go to adjust algorithms or backend/main.go for backend logic.
Rebuild and restart the services:

bash

docker-compose up --build

Monitoring and AnalysisPrometheus: Configured in prometheus/prometheus.yml to scrape metrics (e.g., request counts, response times, active connections) from the load balancer and backends.
Grafana: Configured in grafana/datasource.yml to use Prometheus as a data source. Create dashboards to compare algorithm performance (e.g., latency, throughput, connection distribution).
Comparative Analysis: The FinalProjectReportToPrint.pdf compares RR, WRR, and Least Connections based on metrics like latency, throughput, and load distribution across backends with different capacities.

NetworkAll services communicate over a bridge network (app-network) defined in docker-compose.yml, ensuring isolated and secure communication.DocumentationThe FinalProjectReportToPrint.pdf provides:Implementation details for Round-Robin, Weighted Round-Robin, and Least Connections.
Comparative analysis of algorithm performance under varying loads and backend capacities.
Setup instructions and system architecture.

ContributingFork the repository.
Create a new branch (git checkout -b feature-branch).
Make your changes and commit (git commit -m "Add feature").
Push to the branch (git push origin feature-branch).
Open a pull request.

LicenseThis project is licensed under the MIT License. See the LICENSE file for details.NotesAlgorithms: The README explicitly lists Round-Robin, Weighted Round-Robin, and Least Connections as the implemented algorithms, with brief descriptions of each. I assumed WRR uses the CAPACITY environment variable (small, medium, large) to determine weights; let me know if weights are configured differently.
CPU Limits: The docker-compose.yml specifies cpus: "1" for all backends, indicating one full CPU core, despite comments suggesting fractional values (e.g., "20% of one CPU core"). If you intended fractional CPU limits (e.g., 0.2, 0.5), please confirm, and I’ll update the README.
Metrics: I included example metrics (latency, throughput, connection counts) relevant to comparing load balancing algorithms. If your project tracks specific metrics, provide them, and I can refine the monitoring section.
Algorithm Switching: I assumed algorithms are switched via code changes in main.go or environment variables. If your load balancer supports a specific mechanism (e.g., API, config file), let me know, and I’ll update the instructions.
Report: The README highlights the comparative analysis in FinalProjectReportToPrint.pdf. If you want to emphasize specific findings (e.g., WRR outperforms RR under high load), I can add them.
Repository URL and License: These are placeholders; replace them with actual values.
```
