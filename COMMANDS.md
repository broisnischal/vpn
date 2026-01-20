# Quick Command Reference

Copy-paste these commands to get started!

## Installation

```bash
# Install Go (Arch Linux)
sudo pacman -S go

# Verify Go
go version

# Build VPN
cd /home/nees/omail
./quick-start.sh
```

## Running VPN

### Server (Terminal 1)
```bash
cd /home/nees/omail
sudo ./bin/omail-server \
  -address :51820 \
  -password my-password-123 \
  -tun-ip 10.0.0.1
```

### Client (Terminal 2)
```bash
cd /home/nees/omail
sudo ./bin/omail-client \
  -server localhost:51820 \
  -password my-password-123 \
  -tun-ip 10.0.0.2
```

### Test (Terminal 3)
```bash
# Check interfaces
ip addr show omail0

# Ping test
ping 10.0.0.1

# Check routes
ip route show
```

## Useful Commands

```bash
# Automated test
sudo ./test-vpn.sh test

# Clean up TUN interface
sudo ip link delete omail0

# Monitor VPN traffic
sudo tcpdump -i omail0 -n

# Check what's using port 51820
sudo netstat -unap | grep 51820

# View server logs (if running in foreground)
# (logs appear in terminal where server is running)
```

## Docker Commands

```bash
# Start server
docker-compose up -d vpn-server

# View logs
docker-compose logs -f vpn-server

# Stop server
docker-compose down
```

## Split Tunnel Example

```bash
# Only route specific networks through VPN
sudo ./bin/omail-client \
  -server server.com:51820 \
  -password secret123 \
  -split-tunnel "10.0.0.0/8,192.168.1.0/24"
```

## Troubleshooting Commands

```bash
# Check if Go is installed
go version

# Check if binaries exist
ls -la bin/

# Rebuild if needed
make build

# Check TUN interface status
ip link show omail0

# Check IP configuration
ip addr show omail0

# Check routing table
ip route show

# Test DNS
nslookup google.com

# Test connectivity
ping -c 4 10.0.0.1
```
