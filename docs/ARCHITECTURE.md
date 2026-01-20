# Architecture Overview

This document explains the internal architecture of the OMail VPN server.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        VPN Client                            │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐               │
│  │   App    │───►│   OS     │───►│   TUN    │               │
│  │          │    │ Routing  │    │Interface │               │
│  └──────────┘    └──────────┘    └────┬─────┘               │
│                                         │                     │
│                                         ▼                     │
│                                  ┌──────────┐                │
│                                  │  Client  │                │
│                                  │  Handler │                │
│                                  └────┬─────┘                │
│                                       │                       │
│                                       │ Encrypt               │
│                                       │ & Encode              │
│                                       ▼                       │
│                                  ┌──────────┐                │
│                                  │   UDP    │                │
│                                  │  Socket  │                │
│                                  └────┬─────┘                │
└───────────────────────────────────────┼───────────────────────┘
                                         │
                                         │ UDP (Encrypted)
                                         │
┌────────────────────────────────────────┼───────────────────────┐
│                                        ▼                       │
│                                  ┌──────────┐                │
│                                  │   UDP    │                │
│                                  │  Socket  │                │
│                                  └────┬─────┘                │
│                                       │                       │
│                                       │ Decrypt               │
│                                       │ & Decode              │
│                                       ▼                       │
│                                  ┌──────────┐                │
│                                  │  Server  │                │
│                                  │  Handler │                │
│                                  └────┬─────┘                │
│                                       │                       │
│                                       ▼                       │
│                                  ┌──────────┐                │
│                                  │   TUN    │                │
│                                  │Interface │                │
│                                  └────┬─────┘                │
│                                       │                       │
│                                       ▼                       │
│                                  ┌──────────┐                │
│                                  │   OS     │                │
│                                  │ Routing  │                │
│                                  └────┬─────┘                │
│                                       │                       │
│                                       ▼                       │
│                                  ┌──────────┐                │
│                                  │Internet  │                │
│                                  │          │                │
│                                  └──────────┘                │
│                        VPN Server                             │
└───────────────────────────────────────────────────────────────┘
```

## Component Breakdown

### 1. TUN Interface (`internal/tun/`)

**Purpose**: Creates and manages the virtual network interface.

**Key Functions**:
- `New()`: Creates TUN interface
- `SetIP()`: Assigns IP address and netmask
- `Up()/Down()`: Activates/deactivates interface
- `Read()/Write()`: Reads/writes IP packets
- `SetMTU()`: Sets Maximum Transmission Unit

**Platform Support**:
- Linux: Uses `ioctl` syscalls
- macOS: Uses `ifconfig` commands

### 2. Crypto Layer (`internal/crypto/`)

**Purpose**: Encrypts/decrypts VPN packets.

**Implementation**:
- **Algorithm**: AES-256-GCM
- **Key Derivation**: PBKDF2 with SHA-256 (4096 iterations)
- **Nonce**: 12-byte random nonce per packet
- **Overhead**: ~28 bytes per packet (nonce + GCM tag)

**Key Functions**:
- `NewCrypto()`: Creates crypto from password
- `Encrypt()`: Encrypts plaintext
- `Decrypt()`: Decrypts ciphertext

### 3. Protocol Layer (`internal/protocol/`)

**Purpose**: Defines packet format and encapsulation.

**Packet Structure**:
```
+--------+--------+--------+--------+
| Type   |Reserved| Length |SessionID|
+--------+--------+--------+--------+
|              Data (IP Packet)      |
+-----------------------------------+
```

**Packet Types**:
- `PacketTypeData`: Encapsulated IP packet
- `PacketTypeKeepAlive`: Heartbeat packet

**Key Functions**:
- `Encode()`: Converts packet to bytes
- `Decode()`: Parses bytes to packet
- `NewDataPacket()`: Creates data packet
- `NewKeepAlivePacket()`: Creates keep-alive

### 4. Routing Manager (`internal/routing/`)

**Purpose**: Manages routing table entries.

**Operations**:
- `AddRoute()`: Adds route through VPN
- `DeleteRoute()`: Removes route
- `ListRoutes()`: Lists current routes
- `SetupDefaultRoute()`: Full tunnel mode
- `SetupSplitTunnel()`: Split tunnel mode
- `Cleanup()`: Removes all VPN routes

**Platform Commands**:
- Linux: `ip route add/del`
- macOS: `route add/delete`

### 5. Server (`internal/server/`)

**Purpose**: VPN server that handles client connections.

**Components**:
- **UDP Listener**: Receives encrypted packets from clients
- **Client Manager**: Tracks connected clients
- **TUN Reader**: Reads packets from TUN, forwards to clients
- **UDP Reader**: Reads packets from UDP, forwards to TUN
- **Keep-Alive Handler**: Maintains client sessions

**Session Management**:
- Each client has unique SessionID
- Keep-alive packets maintain session
- Inactive clients timeout after 60 seconds

**Key Functions**:
- `NewServer()`: Creates server instance
- `Start()`: Starts server
- `Stop()`: Stops server
- `handleDataPacket()`: Processes data packets
- `sendToClient()`: Sends packet to client

### 6. Client (`internal/client/`)

**Purpose**: VPN client that connects to server.

**Components**:
- **UDP Connection**: Connects to server
- **TUN Reader**: Reads packets from TUN, sends to server
- **UDP Reader**: Reads packets from UDP, writes to TUN
- **Keep-Alive Sender**: Sends periodic keep-alives
- **Routing Setup**: Configures routing tables

**Connection Flow**:
1. Create TUN interface
2. Set IP address
3. Bring interface up
4. Connect to server via UDP
5. Send keep-alive to establish session
6. Setup routing tables
7. Start packet forwarding

**Key Functions**:
- `NewClient()`: Creates client instance
- `Connect()`: Connects to server
- `Disconnect()`: Disconnects and cleans up
- `sendToServer()`: Sends packet to server
- `setupRouting()`: Configures routes

## Data Flow

### Client → Server (Outbound)

1. **Application** sends IP packet
2. **OS Routing** routes packet to TUN interface (based on routing table)
3. **TUN Interface** captures packet
4. **Client Handler** reads packet from TUN
5. **Protocol Layer** encapsulates packet (adds header)
6. **Crypto Layer** encrypts packet
7. **UDP Socket** sends encrypted packet to server
8. **Server UDP Socket** receives packet
9. **Crypto Layer** decrypts packet
10. **Protocol Layer** decapsulates packet (removes header)
11. **Server Handler** writes packet to server's TUN
12. **Server OS Routing** routes packet to internet

### Server → Client (Inbound)

1. **Internet** sends IP packet to server
2. **Server OS Routing** routes packet to TUN interface
3. **TUN Interface** captures packet
4. **Server Handler** reads packet from TUN
5. **Protocol Layer** encapsulates packet
6. **Crypto Layer** encrypts packet
7. **UDP Socket** sends to client (based on SessionID)
8. **Client UDP Socket** receives packet
9. **Crypto Layer** decrypts packet
10. **Protocol Layer** decapsulates packet
11. **Client Handler** writes packet to client's TUN
12. **Client OS Routing** routes packet to application

## Concurrency Model

### Server

- **Main Goroutine**: Handles signals, cleanup
- **TUN Reader Goroutine**: Reads from TUN, forwards to clients
- **UDP Reader Goroutine**: Reads from UDP, forwards to TUN
- **Cleanup Goroutine**: Removes inactive clients

### Client

- **Main Goroutine**: Handles signals, cleanup
- **TUN Reader Goroutine**: Reads from TUN, sends to server
- **UDP Reader Goroutine**: Reads from UDP, writes to TUN
- **Keep-Alive Goroutine**: Sends periodic keep-alives

## Error Handling

### TUN Interface Errors
- **Creation Failure**: Returns error, doesn't start
- **Read/Write Errors**: Logged, continues operation
- **IP Assignment Failure**: Returns error during initialization

### Network Errors
- **UDP Errors**: Logged, continues operation
- **Connection Errors**: Logged, client retries
- **Decryption Errors**: Logged, packet dropped

### Routing Errors
- **Route Add Failure**: Logged, continues (may affect connectivity)
- **Route Delete Failure**: Logged, continues

## Security Considerations

### Current Implementation

✅ **Encryption**: AES-256-GCM
✅ **Key Derivation**: PBKDF2
✅ **Session Management**: Unique SessionIDs
✅ **Keep-Alive**: Prevents stale connections

⚠️ **Limitations** (for production):

- **Password-based**: No certificate authentication
- **No Perfect Forward Secrecy**: Same key for all packets
- **No Rate Limiting**: Vulnerable to DoS
- **No Authentication**: Any client with password can connect
- **No Audit Logging**: No security event logging

### Recommended Improvements

1. **TLS Handshake**: Use TLS for initial key exchange
2. **Certificate Auth**: Client certificates for authentication
3. **Key Rotation**: Rotate encryption keys periodically
4. **Rate Limiting**: Limit packets per client
5. **Audit Logging**: Log all security events
6. **IP Whitelisting**: Restrict client IPs

## Performance Considerations

### Packet Size
- **MTU**: Default 1500 bytes
- **Overhead**: ~28 bytes encryption + 8 bytes protocol header
- **Effective MTU**: ~1464 bytes

### Throughput
- **Single-threaded**: One goroutine per direction
- **Bottleneck**: Encryption/decryption (CPU-bound)
- **Optimization**: Could use AES-NI hardware acceleration

### Latency
- **UDP**: Low latency (no connection overhead)
- **Encryption**: Minimal overhead (~1-2ms per packet)
- **Routing**: OS routing table lookup (~microseconds)

## Scalability

### Current Limitations

- **Single Server**: One server instance
- **No Load Balancing**: No multi-server support
- **No Clustering**: No server coordination

### Potential Improvements

1. **Multiple Servers**: Load balance across servers
2. **Server Clustering**: Coordinate multiple servers
3. **Connection Pooling**: Reuse UDP connections
4. **Packet Batching**: Batch multiple packets

## Testing Strategy

### Unit Tests
- Crypto layer encryption/decryption
- Protocol packet encoding/decoding
- Routing table operations

### Integration Tests
- TUN interface creation/destruction
- End-to-end packet forwarding
- Routing table management

### Manual Testing
- Server startup/shutdown
- Client connection/disconnection
- Packet capture and verification
- Routing table inspection

## Future Enhancements

1. **IPv6 Support**: Full IPv6 packet handling
2. **Compression**: Compress packets before encryption
3. **Multipath**: Use multiple network paths
4. **QoS**: Quality of Service prioritization
5. **Statistics**: Connection and traffic statistics
6. **Web UI**: Web-based management interface
