# Setup and Testing Guide

## Step 1: Install Go

### Option A: Using Package Manager (Recommended)

**Arch Linux:**
```bash
sudo pacman -S go
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go
```

**macOS:**
```bash
brew install go
```

### Option B: Manual Installation

```bash
# Download Go (latest version)
cd /tmp
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# Remove old installation (if any)
sudo rm -rf /usr/local/go

# Extract to /usr/local
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin

# Reload shell
source ~/.bashrc  # or source ~/.zshrc
```

**Verify installation:**
```bash
go version
# Should show: go version go1.21.6 linux/amd64
```

## Step 2: Install Project Dependencies

```bash
cd /home/nees/omail
go mod download
go mod tidy
```

## Step 3: Build the VPN

```bash
# Build both server and client
make build

# Or manually:
go build -o bin/omail-server ./cmd/server
go build -o bin/omail-client ./cmd/client
```

## Step 4: Test Locally

### Terminal 1: Start VPN Server

```bash
cd /home/nees/omail

# Run server (requires sudo for TUN interface)
sudo ./bin/omail-server \
  -address :51820 \
  -password my-secure-password-123 \
  -tun-ip 10.0.0.1
```

**Expected output:**
```
VPN Server listening on :51820
TUN interface: omail0
```

### Terminal 2: Start VPN Client

Open a new terminal:

```bash
cd /home/nees/omail

# Run client (requires sudo for TUN interface)
sudo ./bin/omail-client \
  -server localhost:51820 \
  -password my-secure-password-123 \
  -tun-ip 10.0.0.2
```

**Expected output:**
```
Connecting to VPN server at localhost:51820
TUN interface: omail0
Connected to VPN server
```

### Terminal 3: Test the Connection

```bash
# Check TUN interfaces
ip addr show omail0

# Check routing table
ip route show | grep omail0

# Ping server from client
ping 10.0.0.1

# Monitor VPN traffic
sudo tcpdump -i omail0 -n
```

## Step 5: Verify Everything Works

### Check Server TUN Interface

```bash
# Should show: inet 10.0.0.1/24
ip addr show omail0
```

### Check Client TUN Interface

```bash
# Should show: inet 10.0.0.2/24
ip addr show omail0
```

### Test Connectivity

From client terminal:
```bash
# Ping server
ping -c 4 10.0.0.1

# Should see successful pings!
```

### Check Routing

```bash
# See routes through VPN
ip route show

# Should see routes like:
# 10.0.0.0/24 dev omail0
```

## Common Issues and Solutions

### Issue: "go: command not found"

**Solution:** Go is not in PATH. Add to your shell config:
```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
source ~/.zshrc
```

### Issue: "operation not permitted" when creating TUN

**Solution:** Run with sudo:
```bash
sudo ./bin/omail-server ...
```

### Issue: "TUN interface already exists"

**Solution:** Delete existing interface:
```bash
sudo ip link delete omail0
```

### Issue: "bind: address already in use"

**Solution:** Port 51820 is in use. Use different port:
```bash
./bin/omail-server -address :51821 ...
```

## Next Steps

1. **Read the documentation:**
   - `README.md` - Full documentation
   - `docs/QUICKSTART.md` - Quick start guide
   - `docs/ROUTING.md` - Learn about routing

2. **Experiment:**
   - Try split tunneling
   - Test with different IP ranges
   - Monitor traffic with tcpdump

3. **Deploy:**
   - Use Docker for production
   - Configure firewall rules
   - Set up on remote server
