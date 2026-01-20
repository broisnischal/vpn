# 🚀 START HERE - Get Your VPN Running in 5 Minutes!

## Step-by-Step Instructions

### Step 1: Install Go (Required)

Open a terminal and run:

```bash
sudo pacman -S go
```

**Verify installation:**
```bash
go version
```

You should see: `go version go1.21.x linux/amd64`

### Step 2: Build the VPN

```bash
cd /home/nees/omail
./quick-start.sh
```

This will:
- ✅ Check Go installation
- ✅ Download dependencies
- ✅ Build server and client binaries

### Step 3: Test Locally

You need **3 terminals** for testing:

#### Terminal 1: Start VPN Server

```bash
cd /home/nees/omail
sudo ./bin/omail-server \
  -address :51820 \
  -password test-password-123 \
  -tun-ip 10.0.0.1
```

**Expected output:**
```
VPN Server listening on :51820
TUN interface: omail0
```

#### Terminal 2: Start VPN Client

```bash
cd /home/nees/omail
sudo ./bin/omail-client \
  -server localhost:51820 \
  -password test-password-123 \
  -tun-ip 10.0.0.2
```

**Expected output:**
```
Connecting to VPN server at localhost:51820
TUN interface: omail0
Connected to VPN server
```

#### Terminal 3: Test Connection

```bash
# Check TUN interface
ip addr show omail0

# Ping server from client
ping 10.0.0.1

# Check routing
ip route show | grep omail0
```

**Success indicators:**
- ✅ `ping 10.0.0.1` works
- ✅ TUN interface shows IP addresses
- ✅ Routes are configured

### Step 4: Automated Testing (Alternative)

Instead of manual testing, you can use:

```bash
sudo ./test-vpn.sh test
```

This runs automated tests and shows you what's working.

## Quick Reference

### Start Server
```bash
sudo ./bin/omail-server -address :51820 -password YOUR_PASSWORD -tun-ip 10.0.0.1
```

### Start Client
```bash
sudo ./bin/omail-client -server SERVER_IP:51820 -password YOUR_PASSWORD -tun-ip 10.0.0.2
```

### Check Status
```bash
ip addr show omail0      # See TUN interface
ip route show            # See routing table
ping 10.0.0.1           # Test connectivity
```

## Common Commands

```bash
# Build everything
make build

# Run server
sudo ./bin/omail-server -address :51820 -password secret123

# Run client
sudo ./bin/omail-client -server localhost:51820 -password secret123

# Clean up TUN interface
sudo ip link delete omail0

# View VPN traffic
sudo tcpdump -i omail0
```

## Troubleshooting

### "go: command not found"
→ Install Go: `sudo pacman -S go`

### "operation not permitted"
→ Run with `sudo`

### "TUN interface already exists"
→ Delete it: `sudo ip link delete omail0`

### "Connection refused"
→ Check server is running and firewall allows UDP port 51820

## Learn More

- **TUTORIAL.md** - Complete usage guide
- **SETUP.md** - Detailed setup instructions
- **README.md** - Full documentation
- **docs/ROUTING.md** - Learn about routing tables
- **docs/ARCHITECTURE.md** - Understand how it works

## What You've Built

✅ Self-hosted VPN server
✅ Encrypted tunnel (AES-256-GCM)
✅ TUN interface management
✅ Routing table automation
✅ Full tunnel and split tunnel support
✅ Docker deployment ready

**Congratulations!** You now have a working VPN! 🎉
