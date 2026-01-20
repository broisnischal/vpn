# OMail VPN Server

A self-hosted VPN server implementation with TUN interface support, built in Go. This project demonstrates VPN concepts including packet encapsulation, encryption, routing tables, and network interface management.

## Features

- **TUN Interface**: Full TUN/TAP support for creating virtual network interfaces
- **Encryption**: AES-256-GCM encryption for all VPN traffic
- **Routing Management**: Automatic routing table configuration (full tunnel or split tunnel)
- **Docker Support**: Ready-to-use Docker containers for easy deployment
- **Cloud Ready**: Automated deployment scripts for AWS and Azure
- **Cross-Platform**: Works on Linux and macOS
- **Educational**: Learn about VPN protocols, routing, and network programming

## Architecture

```
┌─────────────┐         UDP (Encrypted)         ┌─────────────┐
│   Client    │◄───────────────────────────────►│   Server    │
│             │                                   │             │
│  ┌────────┐│                                   │  ┌────────┐│
│  │  TUN   ││                                   │  │  TUN   ││
│  │Interface│                                   │  │Interface│
│  └────────┘│                                   │  └────────┘│
│      │      │                                   │      │      │
│      ▼      │                                   │      ▼      │
│  Routing    │                                   │  Routing    │
│   Tables    │                                   │   Tables    │
└─────────────┘                                   └─────────────┘
```

### Components

1. **TUN Interface**: Virtual network interface that captures IP packets
2. **Protocol Layer**: Packet encapsulation/decapsulation with session management
3. **Crypto Layer**: AES-256-GCM encryption for secure communication
4. **Routing Manager**: Handles routing table operations
5. **Server**: Listens on UDP, manages client sessions, forwards packets
6. **Client**: Connects to server, routes traffic through VPN

## Prerequisites

- Go 1.21 or later
- Linux or macOS (for TUN support)
- Root/Administrator privileges (for TUN interface creation)
- Docker and Docker Compose (for containerized deployment)

## Quick Start

### Cloud Deployment (AWS/Azure)

Deploy your VPN server in the cloud in minutes:

**AWS:**
```bash
# Upload project to EC2
scp -r omail/ ubuntu@YOUR_EC2_IP:~/

# SSH and deploy
ssh ubuntu@YOUR_EC2_IP
cd omail && sudo bash deploy/aws-deploy.sh
```

**Azure:**
```bash
# Upload project to VM
scp -r omail/ azureuser@YOUR_AZURE_IP:~/

# SSH and deploy
ssh azureuser@YOUR_AZURE_IP
cd omail && sudo bash deploy/azure-deploy.sh
```

See [CLOUD-SETUP.md](./CLOUD-SETUP.md) for detailed instructions.

### Local Development

1. **Install dependencies**:
```bash
go mod download
```

2. **Build the binaries**:
```bash
make build
# or
go build -o bin/omail-server ./cmd/server
go build -o bin/omail-client ./cmd/client
```

3. **Run the server** (requires root):
```bash
sudo ./bin/omail-server -address :51820 -password your-secure-password
```

4. **Run the client** (requires root):
```bash
sudo ./bin/omail-client -server localhost:51820 -password your-secure-password
```

### Docker Deployment

1. **Set environment variables**:
```bash
export VPN_PASSWORD=your-secure-password
```

2. **Start the server**:
```bash
docker-compose up -d vpn-server
```

3. **View logs**:
```bash
docker-compose logs -f vpn-server
```

## Configuration

### Server Options

```
-address string
    Server listen address (default ":51820")
-password string
    Encryption password (required)
-tun string
    TUN interface name (default "omail0")
-tun-ip string
    TUN interface IP address (default "10.0.0.1")
-tun-netmask string
    TUN interface netmask (default "255.255.255.0")
-mtu int
    MTU size (default 1500)
```

### Client Options

```
-server string
    Server address (e.g., server.com:51820) (required)
-password string
    Encryption password (required)
-tun string
    TUN interface name (default "omail0")
-tun-ip string
    TUN interface IP address (default "10.0.0.2")
-tun-netmask string
    TUN interface netmask (default "255.255.255.0")
-mtu int
    MTU size (default 1500)
-split-tunnel string
    Comma-separated list of CIDR networks for split tunneling
    (empty for full tunnel)
```

### Example: Split Tunneling

Only route specific networks through VPN:

```bash
sudo ./bin/omail-client \
  -server vpn.example.com:51820 \
  -password mypassword \
  -split-tunnel "10.0.0.0/8,192.168.1.0/24"
```

## Understanding Routing Tables

### What are Routing Tables?

Routing tables tell the operating system how to route network packets. When you send a packet, the OS checks the routing table to determine:
- Which network interface to use
- What gateway (next hop) to send the packet to
- What route has the highest priority (metric)

### Viewing Routing Tables

**Linux**:
```bash
ip route show
# or
route -n
```

**macOS**:
```bash
netstat -rn
# or
route -n get default
```

### How VPN Routing Works

1. **Full Tunnel**: All traffic goes through VPN
   ```
   0.0.0.0/0 → dev omail0
   ```
   This routes ALL internet traffic through the VPN.

2. **Split Tunnel**: Only specific networks go through VPN
   ```
   10.0.0.0/8 → dev omail0
   192.168.1.0/24 → dev omail0
   ```
   Only traffic to these networks goes through VPN.

3. **Default Route**: Normal internet traffic
   ```
   0.0.0.0/0 → dev eth0 via 192.168.1.1
   ```
   This is your normal internet connection.

### Routing Table Operations

Our VPN client automatically manages routing:

- **Add Route**: Adds routes through the TUN interface
- **Delete Route**: Removes routes when disconnecting
- **List Routes**: Shows current routing configuration

## How It Works

### 1. TUN Interface Creation

The TUN (Tunnel) interface is a virtual network interface that operates at Layer 3 (IP layer). When created:

```go
tun, err := tun.New("omail0", 1500)
tun.SetIP(net.ParseIP("10.0.0.1"), net.CIDRMask(24, 32))
tun.Up()
```

This creates a virtual network interface that can capture and inject IP packets.

### 2. Packet Flow

**Client → Server**:
1. Application sends IP packet
2. OS routes packet to TUN interface (based on routing table)
3. VPN client reads packet from TUN
4. Packet is encrypted and encapsulated
5. Encrypted packet sent via UDP to server
6. Server decrypts and writes to server's TUN interface
7. Server's OS routes packet to destination

**Server → Client**:
1. Server receives IP packet on its TUN interface
2. Packet is encrypted and sent to client
3. Client decrypts and writes to client's TUN interface
4. Client's OS routes packet to application

### 3. Encryption

All packets are encrypted using AES-256-GCM:
- **Key Derivation**: PBKDF2 with SHA-256 (4096 iterations)
- **Encryption**: AES-256-GCM with random nonce
- **Authentication**: GCM provides authentication

### 4. Protocol

Each packet has a header:
```
+--------+--------+--------+--------+
| Type   | Reserved| Length | SessionID |
+--------+--------+--------+--------+
|              Data                  |
+-----------------------------------+
```

- **Type**: Data packet or keep-alive
- **Length**: Payload length
- **SessionID**: Client session identifier

## Security Considerations

⚠️ **This is an educational project**. For production use, consider:

1. **Stronger Key Exchange**: Implement ECDH or similar for key exchange
2. **Certificate-based Auth**: Use TLS certificates instead of passwords
3. **Perfect Forward Secrecy**: Rotate keys periodically
4. **Rate Limiting**: Prevent DoS attacks
5. **Connection Authentication**: Verify client identity
6. **Audit Logging**: Log security events

## Troubleshooting

### TUN Interface Not Created

**Error**: `operation not permitted`

**Solution**: Run with root privileges:
```bash
sudo ./bin/omail-server ...
```

### Cannot Set IP Address

**Error**: `SIOCSIFADDR: operation not permitted`

**Solution**: Ensure you have NET_ADMIN capability (Docker) or root access.

### Routing Not Working

**Check routing table**:
```bash
ip route show
```

**Manually add route** (if needed):
```bash
sudo ip route add 10.0.0.0/8 dev omail0
```

### Docker: TUN Device Not Available

**Error**: `open /dev/net/tun: no such file or directory`

**Solution**: Ensure Docker has TUN device access:
```yaml
devices:
  - /dev/net/tun:/dev/net/tun
cap_add:
  - NET_ADMIN
```

## Development

### Project Structure

```
omail/
├── cmd/
│   ├── server/          # Server entry point
│   └── client/          # Client entry point
├── internal/
│   ├── crypto/          # Encryption layer
│   ├── protocol/        # Packet protocol
│   ├── routing/         # Routing management
│   ├── server/          # Server implementation
│   ├── client/          # Client implementation
│   └── tun/             # TUN interface wrapper
├── docker/
│   ├── Dockerfile.server
│   └── Dockerfile.client
├── docker-compose.yml
├── Makefile
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
make fmt
# or
go fmt ./...
```

## Learning Resources

### VPN Concepts

- **TUN/TAP**: Virtual network interfaces
- **Packet Encapsulation**: Wrapping IP packets in another protocol
- **Routing Tables**: How OS routes network traffic
- **NAT Traversal**: Getting through firewalls

### Network Programming

- **UDP Socket Programming**: Connectionless communication
- **Raw Sockets**: Low-level packet access
- **Network Namespaces**: Isolated network environments

### Cryptography

- **Symmetric Encryption**: AES-GCM
- **Key Derivation**: PBKDF2
- **Nonce Management**: Random nonces for GCM

## License

MIT License - feel free to use this for learning and experimentation!

## Contributing

Contributions welcome! This is an educational project, so focus on:
- Code clarity and documentation
- Educational examples
- Security improvements
- Cross-platform support

## Acknowledgments

Inspired by WireGuard, OpenVPN, and other VPN implementations. Built for educational purposes to understand VPN internals.
