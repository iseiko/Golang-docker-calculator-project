# Golang Calculator — Dockerized with IP-Hash Load Balancing

<div align="center">

A learning project exploring **HTTP protocols**, **Golang**, **Docker**, and **load balancing** through a containerized web-based calculator.

![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)
![Nginx](https://img.shields.io/badge/Nginx-009639?logo=nginx&logoColor=white)
![License](https://img.shields.io/badge/License-GPL%20v3-blue)

</div>

## 1 Overview

This is a **college project** designed to learn the fundamentals of:
- **HTTP Protocol** - Understanding how web requests and responses work
- **Golang** - Building efficient web applications with Go
- **Docker** - Containerizing applications for consistency across environments
- **Load Balancing** - Distributing requests across multiple instances using IP-hash strategy

The project implements a simple web-based calculator that performs basic arithmetic operations across multiple containerized instances, with **Nginx** as a load balancer routing requests based on client IP.

## 2 Features

- **Basic Mathematical Operations:**
  - Addition
  - Subtraction
  - Multiplication
  - Division (with zero-division handling)
  - Exponentiation (Power)

- **IP-Hash Load Balancing:** Nginx distributes incoming requests to calculator instances based on client IP address, ensuring session stickiness

- **Docker Containerization:** Multi-container setup with Docker Compose for easy deployment

- **HTTP Server:** Pure Golang HTTP implementation with no external web framework dependencies

## 3 Architecture

```
┌─────────────────────────────────────────┐
│           Client Requests               │
│              (Port 80)                  │
└────────────────┬────────────────────────┘
                 │
        ┌────────▼────────┐
        │  Nginx (LB)     │
        │  IP-Hash        │
        └────────┬────────┘
                 │
         ┌───────┴───────┐
         │               │
    ┌────▼───┐      ┌────▼───┐
    │ App 1  │      │ App 2  │
    │ :8080  │      │ :8080  │
    └────────┘      └────────┘
    (Go Process)    (Go Process)
```

### Components

| Component | Purpose | Technology |
|-----------|---------|-----------|
| **Nginx** | Reverse proxy & load balancer | Nginx Alpine |
| **App 1 & 2** | Calculator instances | Go 1.26 |
| **Dockerfile** | Container image definition | Multi-stage build |
| **docker-compose.yml** | Orchestration | Docker Compose |

## 4 Getting Started

### Prerequisites

- **Docker** and **Docker Compose** installed on your system
- (Optional) Go 1.26+ if building locally

### Installation & Running

1. **Clone the repository:**
```bash
git clone https://github.com/iseiko/Golang-docker-calculator-project.git
cd Golang-docker-calculator-project/calculator
```

2. **Start the services:**
```bash
docker-compose up --build
```

The calculator will be available at `http://localhost`

### Usage

Make HTTP requests to perform calculations. Example with `curl`:

```bash
# Addition: 5 + 3
curl "http://localhost/add?n1=5&n2=3"

# Subtraction: 10 - 4
curl "http://localhost/sub?n1=10&n2=4"

# Multiplication: 6 * 7
curl "http://localhost/mul?n1=6&n2=7"

# Division: 20 / 5
curl "http://localhost/div?n1=20&n2=5"

# Exponentiation: 2 ^ 8
curl "http://localhost/pow?n1=2&n2=8"
```

*Note: The exact endpoint format depends on the implementation. Check the source code for actual endpoint specifications.*

## 5 Project Structure

```
Golang-docker-calculator-project/
├── calculator/
│   ├── app/
│   │   ├── main.go           # Entry point
│   │   ├── funcCalc.go       # Mathematical operations
│   │   ├── Dockerfile        # Container definition
│   │   └── ...
│   ├── docker-compose.yml    # Multi-container setup
│   ├── nginx.conf            # Load balancer config (IP-Hash)
│   └── ...
├── LICENSE
└── README.md
```

## 6 Key Learning Points

### 1. **HTTP Protocol Understanding**
- Stateless request/response model
- HTTP methods and status codes
- Request parsing and routing

### 2. **Go Programming**
- Package structure and interfaces
- HTTP server implementation
- Error handling
- Goroutines and concurrency concepts

### 3. **Docker & Containerization**
- Multi-stage Dockerfile builds (smaller final images)
- Container communication via Docker network
- Volume mounting for configuration

### 4. **Load Balancing Strategy**
- **IP-Hash Algorithm:** Routes all requests from the same client IP to the same backend instance
- **Benefit:** Session affinity without sticky sessions
- **Implementation:** Configured in `nginx.conf` with `ip_hash` directive

## 🛠️ Implementation Details

### Load Balancer Configuration
The Nginx configuration uses IP-hash to consistently route requests:

```nginx
upstream calculator_backend {
    ip_hash;
    server app1:8080;
    server app2:8080;
}
```

### Multi-Stage Docker Build
Reduces final image size by building in a separate stage:

1. **Build stage** (golang:1.26) - Compiles the Go binary
2. **Runtime stage** (debian:bookworm-slim) - Runs only the binary

## 7 Notes

- **No Frontend:** This project is API-only (no web UI). Interact via HTTP requests.
- **Educational Purpose:** Optimized for learning, not production use.
- **Division by Zero:** Handled gracefully (returns 0 instead of panic).

## 8 License

This project is licensed under the **GNU General Public License v3.0** - see the [LICENSE](LICENSE) file for details.

## 9 Purpose

Created as a college project to understand how modern web applications work at the protocol and infrastructure level, bridging the gap between theoretical knowledge and practical implementation.

---

**Questions or improvements?** Feel free to open an issue or submit a pull request!
