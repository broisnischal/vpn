# Connecting Clients Guide

How to connect to your OMail VPN server from laptops, phones, and other devices.

## Table of Contents

1. [Laptop Connection (Linux/macOS/Windows)](#laptop-connection)
2. [Mobile Connection (Android/iOS)](#mobile-connection)
3. [Router/Gateway Setup](#router-setup)
4. [Troubleshooting](#troubleshooting)

## Laptop Connection

### Linux

**Prerequisites:**
- Go installed (to build client)
- Root/sudo access

**Steps:**

1. **Build or download client binary:**
```bash
# If you have the project
cd omail
go build -o bin/omail-client ./cmd/client

# Or download pre-built binary (if available)
```

2. **Connect to VPN:**
```bash
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

3. **Verify connection:**
```bash
# Check TUN interface
ip addr show omail0

# Ping server
ping 10.0.0.1

# Check routing
ip route show
```

4. **Disconnect:**
```bash
# Press Ctrl+C in the client terminal
# Or kill the process
sudo pkill omail-client

# Clean up TUN interface
sudo ip link delete omail0
```

### macOS

**Steps:**

1. **Build client:**
```bash
cd omail
go build -o bin/omail-client ./cmd/client
```

2. **Connect:**
```bash
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

3. **Verify:**
```bash
# Check interface
ifconfig omail0

# Ping server
ping 10.0.0.1

# Check routes
netstat -rn | grep omail0
```

### Windows

**Option 1: Using WSL (Windows Subsystem for Linux)**

1. **Install WSL:**
```powershell
wsl --install
```

2. **In WSL, build and run client:**
```bash
cd /mnt/c/path/to/omail
go build -o bin/omail-client.exe ./cmd/client
sudo ./bin/omail-client.exe \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

**Option 2: Native Windows Build**

1. **Build for Windows:**
```bash
# On Linux/macOS, cross-compile:
GOOS=windows GOARCH=amd64 go build -o bin/omail-client.exe ./cmd/client
```

2. **Run on Windows:**
   - Requires TAP adapter (install OpenVPN to get TAP driver)
   - Run as Administrator
   - May need additional Windows-specific TUN code

**Option 3: Use Linux VM**

Run a Linux VM (VirtualBox/VMware) and connect from there.

## Mobile Connection

### Android

**Method 1: Using Termux (Recommended)**

Termux provides a Linux environment on Android.

1. **Install Termux** from F-Droid or Google Play

2. **Install dependencies:**
```bash
# In Termux
pkg update
pkg install golang git

# Install TUN support (requires root or use VPN mode)
pkg install root-repo
pkg install tsu  # For root access
```

3. **Build client:**
```bash
cd ~
git clone YOUR_REPO_URL
cd omail
go build -o bin/omail-client ./cmd/client
```

4. **Connect (requires root):**
```bash
su
cd /data/data/com.termux/files/home/omail
./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

**Method 2: Using VPN Gateway Router**

Set up VPN on a router/gateway device (see Router Setup section below).

**Method 3: SOCKS5 Proxy (Workaround)**

Create a SOCKS5 proxy on your server that phones can connect to.

### iOS

**Method 1: Using VPN Gateway Router**

iOS doesn't allow custom VPN protocols easily. Best option is to use a router/gateway.

**Method 2: Jailbreak + Terminal**

If jailbroken, you can use terminal apps, but this is not recommended.

**Method 3: SOCKS5 Proxy**

Use a proxy app that supports SOCKS5.

## Router/Gateway Setup

The easiest way to connect phones is to set up the VPN on a router or gateway device.

### Using Raspberry Pi as VPN Gateway

1. **Install on Raspberry Pi:**
```bash
# SSH into Raspberry Pi
ssh pi@raspberry-pi-ip

# Install Go
sudo apt install golang-go

# Build and run VPN client
cd omail
go build -o bin/omail-client ./cmd/client
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

2. **Configure Raspberry Pi as Gateway:**
```bash
# Enable IP forwarding
echo 'net.ipv4.ip_forward=1' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# Configure NAT
sudo iptables -t nat -A POSTROUTING -o omail0 -j MASQUERADE
sudo iptables -A FORWARD -i wlan0 -o omail0 -j ACCEPT
sudo iptables -A FORWARD -i omail0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT
```

3. **Connect phones to Raspberry Pi WiFi:**
   - All traffic from phones goes through VPN automatically

### Using OpenWrt Router

1. **Install OpenWrt on compatible router**
2. **SSH into router**
3. **Install Go and build client** (if router has enough resources)
4. **Configure routing**

## Quick Connection Scripts

### Linux/macOS Connection Script

Create `connect-vpn.sh`:

```bash
#!/bin/bash

SERVER_IP="${VPN_SERVER:-YOUR_SERVER_IP}"
PASSWORD="${VPN_PASSWORD:-YOUR_PASSWORD}"

if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root: sudo $0"
    exit 1
fi

cd "$(dirname "$0")"

./bin/omail-client \
  -server "$SERVER_IP:51820" \
  -password "$PASSWORD" \
  -tun-ip 10.0.0.2
```

Usage:
```bash
chmod +x connect-vpn.sh
sudo VPN_SERVER=server.com VPN_PASSWORD=secret ./connect-vpn.sh
```

### Windows Connection Script

Create `connect-vpn.bat`:

```batch
@echo off
echo Connecting to VPN...
omail-client.exe -server YOUR_SERVER_IP:51820 -password YOUR_PASSWORD -tun-ip 10.0.0.2
pause
```

## Connection Examples

### Example 1: Connect from Home Laptop

```bash
# Your laptop (Linux)
sudo ./bin/omail-client \
  -server 54.123.45.67:51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.2

# Now all traffic goes through VPN
curl ifconfig.me  # Shows server's IP
```

### Example 2: Connect from Coffee Shop

```bash
# Same command, works from anywhere
sudo ./bin/omail-client \
  -server 54.123.45.67:51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.2

# Your traffic is encrypted even on public WiFi
```

### Example 3: Split Tunnel (Only Route Specific Networks)

```bash
# Only route private networks through VPN
sudo ./bin/omail-client \
  -server 54.123.45.67:51820 \
  -password my-secure-password \
  -tun-ip 10.0.0.2 \
  -split-tunnel "10.0.0.0/8,192.168.1.0/24"

# Normal internet uses your connection
# Private networks use VPN
```

## Troubleshooting

### Can't Connect from Laptop

1. **Check server is running:**
```bash
# On server
docker-compose ps
# or
systemctl status omail-vpn
```

2. **Check firewall:**
```bash
# On server
sudo ufw status
# Should show: 51820/udp ALLOW
```

3. **Test connectivity:**
```bash
# From laptop
nc -u -v SERVER_IP 51820
# Should connect
```

4. **Check password:**
```bash
# Verify password matches server
```

### Phone Can't Connect

1. **Use Termux on Android** (requires root)
2. **Use router/gateway method** (easiest)
3. **Use SOCKS5 proxy** (workaround)

### Connection Drops

1. **Check keep-alive:** Client sends every 10 seconds
2. **Check server timeout:** 60 seconds
3. **Check network stability:** Unstable connection causes drops
4. **Check logs:** `docker-compose logs vpn-server`

### No Internet Through VPN

1. **Check NAT on server:**
```bash
sudo iptables -t nat -L
```

2. **Check IP forwarding:**
```bash
sysctl net.ipv4.ip_forward
# Should be: 1
```

3. **Check routing:**
```bash
ip route show
```

## Security Tips

1. **Use strong passwords**
2. **Restrict server access** (firewall rules)
3. **Use split tunnel** when possible
4. **Monitor connections**
5. **Keep server updated**

## Next Steps

1. ✅ Connect laptop
2. ✅ Test connectivity
3. ✅ Set up router/gateway for phones
4. ✅ Configure split tunneling
5. ✅ Monitor usage

Your VPN is now accessible from all your devices! 📱💻
