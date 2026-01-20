# Quick Start Guide

Get your VPN server up and running in minutes!

## Prerequisites

- Linux or macOS
- Go 1.21+ installed
- Root/Administrator access
- Docker (optional, for containerized deployment)

## Method 1: Local Development (Recommended for Learning)

### Step 1: Clone and Setup

```bash
cd omail
go mod download
```

### Step 2: Build Binaries

```bash
make build
# or manually:
go build -o bin/omail-server ./cmd/server
go build -o bin/omail-client ./cmd/client
```

### Step 3: Start Server

**Terminal 1 - Server**:
```bash
sudo ./bin/omail-server \
  -address :51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.1
```

You should see:
```
VPN Server listening on :51820
TUN interface: omail0
```

### Step 4: Start Client

**Terminal 2 - Client** (on same or different machine):
```bash
sudo ./bin/omail-client \
  -server localhost:51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.2
```

You should see:
```
Connecting to VPN server at localhost:51820
TUN interface: omail0
Connected to VPN server
```

### Step 5: Verify Connection

**Check TUN interface**:
```bash
ip addr show omail0
# Should show: inet 10.0.0.1/24 (server) or 10.0.0.2/24 (client)
```

**Check routing**:
```bash
ip route show
# Should show routes through omail0
```

**Test connectivity**:
```bash
# From client, ping server's TUN IP
ping 10.0.0.1

# Check if traffic is going through VPN
tcpdump -i omail0
```

## Method 2: Docker Deployment

### Step 1: Configure Environment

Create `.env` file:
```bash
VPN_PASSWORD=my-secure-password
VPN_ADDRESS=:51820
VPN_TUN_IP=10.0.0.1
```

### Step 2: Start Server

```bash
docker-compose up -d vpn-server
```

### Step 3: View Logs

```bash
docker-compose logs -f vpn-server
```

### Step 4: Connect Client

On your local machine (or another Docker container):
```bash
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password my-secure-password
```

## Method 3: Using Example Scripts

### Server

```bash
sudo VPN_PASSWORD=my-password ./examples/run-server.sh
```

### Client

```bash
sudo VPN_SERVER=server.com:51820 \
  VPN_PASSWORD=my-password \
  ./examples/run-client.sh
```

## Testing the VPN

### 1. Check Interface

```bash
# Linux
ip link show omail0
ip addr show omail0

# macOS
ifconfig omail0
```

### 2. Check Routing

```bash
# Linux
ip route show

# macOS
netstat -rn | grep omail0
```

### 3. Test Connectivity

```bash
# Ping server from client
ping 10.0.0.1

# Check if traffic is encrypted
tcpdump -i omail0 -n
```

### 4. Test Internet (Full Tunnel)

If using full tunnel mode:
```bash
# Check your IP (should be server's IP)
curl ifconfig.me

# Test DNS
nslookup google.com
```

## Common Issues

### Issue: "operation not permitted"

**Solution**: Run with `sudo` or ensure you have NET_ADMIN capability.

### Issue: "TUN interface already exists"

**Solution**: 
```bash
sudo ip link delete omail0
# or use a different name: -tun omail1
```

### Issue: "Cannot assign requested address"

**Solution**: Ensure TUN IPs don't conflict:
- Server: 10.0.0.1
- Client: 10.0.0.2 (or different)

### Issue: "Connection refused"

**Solution**: 
- Check firewall: `sudo ufw allow 51820/udp`
- Verify server is running
- Check server address/port

## Next Steps

- Read [ROUTING.md](./ROUTING.md) to understand routing tables
- Read [README.md](../README.md) for detailed documentation
- Experiment with split tunneling
- Explore the code to learn VPN internals

## Security Note

⚠️ **Change the default password!** The examples use `changeme123` for demonstration only.

For production:
1. Use a strong, random password
2. Consider implementing certificate-based authentication
3. Use firewall rules to restrict access
4. Monitor logs for suspicious activity
