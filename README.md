# Leviathan-Wake-Cockpit

**Leviathan-Wake-Cockpit** is the primary **signal-ingestion microservice** for a **High-Frequency Trading (HFT)** crypto ecosystem.  
It is engineered for **ultra-low latency** and serves as the gateway for real-time trade signal processing.

---

## ✨ Features
- Connects to real-time **WebSocket** streams from both **CEXs** and **on-chain sources**  
- Filters massive streams of trading data to detect **significant whale transactions** from a curated **whitelist**  
- Emits standardized **WhaleSignal** events for downstream services to analyze and execute  
- Fully implemented in **Go (Golang)** to leverage its **powerful concurrency model** and **performance**  

---

## 📂 Project Structure

Leviathan-Wake-Cockpit/
├── cmd/            # Main entrypoint for the service
├── internal/       # Core business logic
├── pkg/            # Shared libraries and reusable modules
├── configs/        # Configuration files
├── scripts/        # DevOps and helper scripts
└── README.md       # Project documentation

---

## 🚀 Quick Start

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)  
- Access to CEX and on-chain WebSocket endpoints  

### Build & Run
```bash
git clone https://github.com/your-org/Leviathan-Wake-Cockpit.git
cd Leviathan-Wake-Cockpit
go build ./cmd/watcher
./watcher 
```

## ⚡ Architecture
Airplane Architecture
- Watcher Service connects directly to live trading streams
- Filters trades and normalizes them into WhaleSignal events
- Downstream services such as Analyzer and Executor consume these signals for strategy execution