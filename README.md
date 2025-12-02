# CS315 Cryptography

A distributed cryptography system with server and node components.

## Requirements

- Go 1.16 or later

## Project Structure

```
src/
├── distribute/          # Server for distributing cryptographic resources
│   ├── initialize/      # Initialization logic
│   ├── models/          # Data structures for distribute module
│   ├── server/          # HTTP server for key distribution
│   └── parameters/      # Configuration files
├── node/                # Node server for inter-node communication
│   ├── initialize/      # Node initialization (fetches parameters from distribute)
│   ├── server/          # HTTP server handling /message and /gcol endpoints
│   └── sendMessage/     # CLI tool for sending messages between nodes
└── shared/              # Shared utilities and types
    ├── types.go         # Common data structures
    └── utils.go         # Helper functions for user input
```

## Running

This example demonstrates running a system with 2 nodes (Node A and Node B).

### Prerequisites

Navigate to the source directory:
```bash
cd src
```

### Step 1: Initialize Distribute Server (if not already done)

```bash
go run ./distribute/initialize
```

This will create the distribute server configuration file with the necessary parameters.

### Step 2: Start Distribute Server

```bash
go run ./distribute/server -port 8080
```

The distribute server now runs on `http://localhost:8080/distribute`

### Step 3: Initialize Nodes

Initialize Node A:
```bash
go run ./node/initialize
```

When prompted:
- Enter distribute server URL: `http://localhost:8080/distribute`
- Enter relative path to store parameters file: `./node_a/parameters.json`

Initialize Node B:
```bash
go run ./node/initialize
```

When prompted:
- Enter distribute server URL: `http://localhost:8080/distribute`
- Enter relative path to store parameters file: `./node_b/parameters.json`

**Note:** Each node must have a unique configuration file path. If running multiple nodes on the same machine, ensure the paths are different (e.g., `./node_a/parameters.json`, `./node_b/parameters.json`).

At this point, the distribute server is no longer needed and can be stopped.

### Step 4: Start Node A Server

```bash
go run ./node/server -port 8080 -config ./node_a/parameters.json
```

Node A now listens on `http://localhost:8080`

### Step 5: Send Message from Node B to Node A

In a new terminal (still in the `src` directory):
```bash
go run ./node/sendMessage
```

When prompted:
- Enter config file path: `./node_b/parameters.json`
- Enter target node URL: `http://localhost:8080`
- Enter message: (your message here)

The program will:
1. Load Node B's configuration
2. Fetch Node A's Gcol from `http://localhost:8080/gcol`
3. Calculate the shared key
4. Send Node B's Gcol and the message to `http://localhost:8080/message`