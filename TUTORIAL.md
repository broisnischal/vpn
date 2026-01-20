# How to Use OMail VPN - Complete Tutorial

This tutorial will teach you how to use your VPN server step by step.

## Prerequisites

1. **Go installed** (see SETUP.md)
2. **Root/sudo access** (for TUN interface)
3. **Two terminals** (one for server, one for client)

## Part 1: Understanding What We're Building

### What is a VPN?

A VPN (Virtual Private Network) creates a secure tunnel between your device and a server. All your internet traffic goes through this encrypted tunnel.

### How Our VPN Works

```
Your Computer → TUN Interface → Encrypted UDP → VPN Server → Internet
```

1. **TUN Interface**: Virtual network interface that captures packets
2. **Encryption**: All packets are encrypted with AES-256-GCM
3. **UDP Transport**: Encrypted packets sent via UDP
4. **Server Forwarding**: Server decrypts and forwards to internet

## Part 2: Installation

### Step 1: Install Go

**Arch Linux:**
```bash
sudo pacman -S go
```

**Verify:**
```bash
go version
# Should show: go version go1.21.x linux/amd64
```

### Step 2: Build the VPN

```bash
cd /home/nees/omail
go mod download
make build
```

This creates:
- `bin/omail-server` - VPN server
- `bin/omail-client` - VPN client

## Part 3: Running the VPN

### Scenario 1: Local Testing (Same Machine)

**Terminal 1 - Start Server:**
```bash
cd /home/nees/omail
sudo ./bin/omail-server \
  -address :51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.1
```

**What happens:**
- Creates TUN interface `omail0`
- Sets IP to `10.0.0.1`
- Listens on UDP port `51820`
- Waits for client connections

**Terminal 2 - Start Client:**
```bash
cd /home/nees/omail
sudo ./bin/omail-client \
  -server localhost:51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.2
```

**What happens:**
- Creates TUN interface `omail0` (client side)
- Sets IP to `10.0.0.2`
- Connects to server
- Sets up routing (all traffic through VPN)

### Scenario 2: Remote Server

**On Server Machine:**
```bash
sudo ./bin/omail-server \
  -address :51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.1
```

**On Client Machine:**
```bash
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.2
```

## Part 4: Understanding the Commands

### Server Options

```bash
./bin/omail-server \
  -address :51820          # Listen on all interfaces, port 51820
  -password secret123       # Encryption password (REQUIRED)
  -tun omail0              # TUN interface name
  -tun-ip 10.0.0.1         # Server's VPN IP address
  -tun-netmask 255.255.255.0  # Network mask
  -mtu 1500                 # Maximum packet size
```

### Client Options

```bash
./bin/omail-client \
  -server server.com:51820  # Server address (REQUIRED)
  -password secret123        # Same password as server
  -tun omail0               # TUN interface name
  -tun-ip 10.0.0.2          # Client's VPN IP address
  -split-tunnel "10.0.0.0/8"  # Only route these networks (optional)
```

## Part 5: Testing the Connection

### Check TUN Interface

```bash
# See TUN interface
ip addr show omail0

# Should show:
# omail0: <UP> ...
#    inet 10.0.0.1/24 (server) or 10.0.0.2/24 (client)
```

### Check Routing

```bash
# View routing table
ip route show

# Should see routes like:
# 10.0.0.0/24 dev omail0
# 0.0.0.0/0 dev omail0 (if full tunnel)
```

### Test Connectivity

```bash
# Ping server from client
ping 10.0.0.1

# Should see successful pings!
```

### Monitor Traffic

```bash
# See packets on TUN interface
sudo tcpdump -i omail0 -n

# You'll see encrypted packets going through
```

## Part 6: Understanding Routing Modes

### Full Tunnel (Default)

All internet traffic goes through VPN:

```bash
sudo ./bin/omail-client \
  -server server.com:51820 \
  -password secret123
```

**Routing table:**
```
0.0.0.0/0 → dev omail0  (all traffic through VPN)
```

**Use case:** Complete privacy, hide IP address

### Split Tunnel

Only specific networks go through VPN:

```bash
sudo ./bin/omail-client \
  -server server.com:51820 \
  -password secret123 \
  -split-tunnel "10.0.0.0/8,192.168.1.0/24"
```

**Routing table:**
```
10.0.0.0/8 → dev omail0
192.168.1.0/24 → dev omail0
0.0.0.0/0 → dev eth0 via 192.168.1.1  (normal internet)
```

**Use case:** Access private networks while keeping normal internet

## Part 7: Practical Examples

### Example 1: Secure Browsing

**Goal:** Route all web traffic through VPN

```bash
# Start server
sudo ./bin/omail-server -address :51820 -password secret123

# Start client (full tunnel)
sudo ./bin/omail-client \
  -server localhost:51820 \
  -password secret123

# Now all browsing goes through VPN!
curl ifconfig.me  # Shows server's IP
```

### Example 2: Access Private Network

**Goal:** Access server's private network (10.0.0.0/8)

```bash
# Start server on network with 10.0.0.0/8
sudo ./bin/omail-server -address :51820 -password secret123

# Start client with split tunnel
sudo ./bin/omail-client \
  -server server.com:51820 \
  -password secret123 \
  -split-tunnel "10.0.0.0/8"

# Now you can access 10.0.0.x addresses!
ping 10.0.0.100
```

### Example 3: Multiple Clients

**Goal:** Connect multiple clients to same server

```bash
# Server (one instance)
sudo ./bin/omail-server -address :51820 -password secret123

# Client 1
sudo ./bin/omail-client \
  -server server.com:51820 \
  -password secret123 \
  -tun-ip 10.0.0.2

# Client 2 (different IP)
sudo ./bin/omail-client \
  -server server.com:51820 \
  -password secret123 \
  -tun-ip 10.0.0.3
```

## Part 8: Troubleshooting

### Problem: "operation not permitted"

**Solution:** Run with `sudo`

### Problem: "TUN interface already exists"

**Solution:**
```bash
sudo ip link delete omail0
```

### Problem: "Connection refused"

**Solution:**
1. Check server is running
2. Check firewall: `sudo ufw allow 51820/udp`
3. Verify server address/port

### Problem: "Cannot ping server"

**Solution:**
1. Check both TUN interfaces are up: `ip link show omail0`
2. Check IPs are set: `ip addr show omail0`
3. Check routing: `ip route show`
4. Check server logs for errors

### Problem: "No internet after connecting"

**Solution:**
- Check routing: `ip route show`
- Verify default route: `ip route get 8.8.8.8`
- Check server can access internet
- Try split tunnel mode

## Part 9: Security Best Practices

### 1. Use Strong Passwords

```bash
# Generate random password
openssl rand -base64 32

# Use it:
./bin/omail-server -password $(openssl rand -base64 32)
```

### 2. Firewall Rules

```bash
# Allow only specific IPs
sudo ufw allow from CLIENT_IP to any port 51820 proto udp

# Or allow all (less secure)
sudo ufw allow 51820/udp
```

### 3. Monitor Connections

```bash
# See active connections
sudo netstat -unap | grep 51820

# Monitor traffic
sudo tcpdump -i omail0
```

## Part 10: Advanced Usage

### Using with Docker

```bash
# Start server
docker-compose up -d vpn-server

# View logs
docker-compose logs -f vpn-server
```

### Using Example Scripts

```bash
# Server
sudo VPN_PASSWORD=secret123 ./examples/run-server.sh

# Client
sudo VPN_SERVER=server.com:51820 \
  VPN_PASSWORD=secret123 \
  ./examples/run-client.sh
```

### Automated Testing

```bash
# Run test suite
sudo ./test-vpn.sh test
```

## Summary

You now know how to:

1. ✅ Install and build the VPN
2. ✅ Start server and client
3. ✅ Understand routing modes
4. ✅ Test connectivity
5. ✅ Troubleshoot issues
6. ✅ Use advanced features

**Next Steps:**
- Read `docs/ROUTING.md` to understand routing tables deeply
- Read `docs/ARCHITECTURE.md` to understand internals
- Experiment with different configurations
- Deploy on a remote server

Happy VPN-ing! 🔒
