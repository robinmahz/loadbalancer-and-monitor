# Go Load Balancer: Comparative Analysis of Load Balancing Algorithms

This project implements and compares three load balancing algorithms—Weighted Round-Robin (WRR), Round-Robin (RR), and Least Connections—using Go. The load balancer distributes HTTP requests across three backend servers with varying capacities (small, medium, large). Prometheus and Grafana are used for monitoring and performance analysis. All services are containerized with Docker and orchestrated using Docker Compose. A detailed report analyzes the performance of the implemented algorithms.

## Project Structure

go-load-balancer/  
├── FinalProjectReportToPrint.pdf  
├── backend/  
│ ├── Dockerfile  
│ └── main.go  
├── docker-compose.yml  
├── grafana/  
│ └── datasource.yml  
├── loadbalancer/  
│ ├── Dockerfile  
│ └── main.go  
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
```

### 2. Configure Container Resources

Defined in `docker-compose.yml`:

- `backend1`: Small capacity, 1 CPU core, 512MB memory, 256MB reserved, port 8081
- `backend2`: Medium capacity, 1 CPU core, 512MB memory, 256MB reserved, port 8082
- `backend3`: Large capacity, 1 CPU core, 512MB memory, 256MB reserved, port 8083
- `loadbalancer`: Runs on port 8000
- `prometheus`: Port 9090
- `grafana`: Port 3000

To modify resource limits or environment variables (e.g., `CAPACITY`, `SERVER_NAME`, `PORT`), edit `docker-compose.yml`.

```yaml
services:
  backend1:
    build: ./backend
    environment:
      - SERVER_NAME=backend1
      - CAPACITY=small
      - PORT=8081
    deploy:
      resources:
        limits:
          cpus: "1"
          memory: 512M
        reservations:
          memory: 256M
```

### 3. Configure Load Balancing Algorithms

The file `loadbalancer/main.go` supports the following algorithms:

- Round-Robin (RR)
- Weighted Round-Robin (WRR)
- Least Connections

To switch between algorithms, you can:

- Modify the code in `main.go` directly
- Or, use environment variables (if this functionality is implemented)

> The `CAPACITY` environment variable (`small`, `medium`, `large`) influences the WRR algorithm by assigning different weights to backends.

For more details on each algorithm’s behavior and performance, refer to `FinalProjectReportToPrint.pdf`.

---

### 4. Build and Run

Use Docker Compose to build and launch the services:

```bash
docker-compose up --build
```

This command will:

- Build and start:
  - The load balancer
  - Three backend services
  - Prometheus (for metrics collection)
  - Grafana (for visualization)
- Expose the following services:
  - **Load Balancer**: `http://localhost:8000`
  - **Grafana**: `http://localhost:3000` (default login: `admin/admin`)
  - **Prometheus**: `http://localhost:9090`

---

### 5. Access the Services

- **Load Balancer**  
  Access via:  
  `http://localhost:8000`

- **Grafana Dashboard**  
  Open in a browser:  
  `http://localhost:3000`  
  _(Default credentials: admin / admin)_

- **Prometheus**  
  View raw metrics at:  
  `http://localhost:9090`

> Note: Backend services are only accessible internally through the load balancer.

---

### 6. Stop the Services

To stop and remove all running containers:

```bash
docker-compose down
```

### Development

To make changes to the load balancer or backend services:

1. Edit `loadbalancer/main.go` or `backend/main.go` as needed.
2. Rebuild and restart the stack:

```
docker-compose up --build
```

---

### Monitoring and Analysis

- **Prometheus**  
  Configured in `prometheus/prometheus.yml` to scrape:

  - Request counts
  - Response times
  - Active connections

- **Grafana**  
  Uses `grafana/datasource.yml` to connect to Prometheus.  
  You can create dashboards to visualize:
  - Request latency
  - Load distribution
  - Throughput trends

---

### Comparative Analysis

The file `FinalProjectReportToPrint.pdf` includes:

- Implementation details for:

  - Round-Robin
  - Weighted Round-Robin
  - Least Connections

- Performance comparison based on:
  - Response latency
  - Request throughput
  - Load distribution under varying capacities

---

### Network

All services communicate using a bridge network (`app-network`) defined in `docker-compose.yml`.

This provides:

- Secure inter-container communication
- Logical isolation of application services

---

### Documentation

See `FinalProjectReportToPrint.pdf` for:

- Setup and deployment instructions
- Load balancing algorithm internals
- Architecture diagrams
- Collected metrics and visual analysis

---

### Contributing

Follow these steps to contribute:

```
# Fork the repository
# Create a new branch
git checkout -b feature-branch

# Make your changes
git commit -m "Add new feature"

# Push the branch
git push origin feature-branch
```

Then open a pull request for review.

---

### License

This project is licensed under the MIT License.  
Refer to the `LICENSE` file for full terms.

---

### Notes

- **WRR Weights**:  
  Assumes `CAPACITY` (small, medium, large) maps to different WRR weights. Update if weights are configured differently.

- **CPU Limits**:  
  Current config uses `cpus: "1"` for each service. Use `"0.2"` or `"0.5"` for fractional limits if needed.

- **Metrics**:  
  Collected metrics include latency, throughput, and connection counts. Add more if your project tracks them.

- **Algorithm Switching**:  
  Currently assumed to be handled via code changes or environment variables. Let me know if this is dynamic or uses a config file.

- **Performance Highlights**:  
  Want to showcase WRR’s superior performance under heavy load? I can help include key findings from the report here.

- **To change load balancing algorithm change the branches**
