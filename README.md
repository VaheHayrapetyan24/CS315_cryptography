# CS315 Cryptography

A distributed cryptography system with server and node components.

## Requirements

- Go 1.16 or later

## Project Structure

```
src/
├── distribute/          # Server for distributing cryptographic resources
│   ├── initialize/      # Initialization logic (runs before server starts)
│   ├── server/          # HTTP server
│   └── parameters/      # Configuration files
└── node/                # Node server for key exchange between nodes
    └── (not yet implemented)
```

## Running

### Distribute Server

Initialize the distribute server:
```bash
go run ./src/distribute/initialize
```

Start the distribute HTTP server:
```bash
go run ./src/distribute/server
```

### Node Server

Node initialization and server (not yet implemented):
```bash
go run ./src/node
```