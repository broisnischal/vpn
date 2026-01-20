# Mobile Connection Guide

How to connect your phone (Android/iOS) to OMail VPN.

## Overview

Connecting phones directly requires root access or a workaround. The easiest methods are:

1. **Router/Gateway Setup** (Recommended) - Set up VPN on a router
2. **Termux on Android** - Use Linux environment on Android
3. **SOCKS5 Proxy** - Use a proxy server

## Method 1: Router/Gateway Setup (Easiest)

Set up the VPN client on a Raspberry Pi or router, then connect your phone to that device's WiFi.

### Using Raspberry Pi

**Step 1: Set up Raspberry Pi**

```bash
# SSH into Raspberry Pi
ssh pi@raspberry-pi-ip

# Install Go
sudo apt update
sudo apt install golang-go

# Clone or upload project
cd ~
git clone YOUR_REPO_URL omail
cd omail

# Build client
go build -o bin/omail-client ./cmd/client
```

**Step 2: Connect Raspberry Pi to VPN**

```bash
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

**Step 3: Configure Raspberry Pi as Gateway**

```bash
# Enable IP forwarding
echo 'net.ipv4.ip_forward=1' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# Configure NAT
sudo iptables -t nat -A POSTROUTING -o omail0 -j MASQUERADE
sudo iptables -A FORWARD -i wlan0 -o omail0 -j ACCEPT
sudo iptables -A FORWARD -i omail0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT

# Save iptables
sudo apt install iptables-persistent
sudo netfilter-persistent save
```

**Step 4: Connect Phone to Raspberry Pi WiFi**

- Connect your phone to Raspberry Pi's WiFi hotspot
- All traffic automatically goes through VPN!

### Using OpenWrt Router

1. Install OpenWrt on compatible router
2. SSH into router
3. Install Go (if router supports it)
4. Build and run VPN client
5. Configure routing

## Method 2: Android with Termux

**Requirements:**
- Android device
- Root access (for TUN interface)
- Termux app

**Steps:**

1. **Install Termux** from F-Droid (not Google Play)

2. **Open Termux and install dependencies:**
```bash
pkg update
pkg install golang git

# For root access
pkg install root-repo
pkg install tsu
```

3. **Clone project:**
```bash
cd ~
git clone YOUR_REPO_URL omail
cd omail
```

4. **Build client:**
```bash
go build -o bin/omail-client ./cmd/client
```

5. **Connect (requires root):**
```bash
# Switch to root
su

# Navigate to project
cd /data/data/com.termux/files/home/omail

# Connect
./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

**Note:** This requires root access. Most modern Android devices don't allow root easily.

## Method 3: SOCKS5 Proxy (Workaround)

Create a SOCKS5 proxy on your server that phones can connect to.

### Server Setup

```bash
# Install SSH (for SOCKS proxy)
sudo apt install openssh-server

# Or use a dedicated SOCKS5 proxy
# Install Dante SOCKS server
sudo apt install dante-server

# Configure Dante
sudo nano /etc/danted.conf
```

### Client Setup

**Android:**
- Use apps like "ProxyDroid" or "Orbot"
- Configure SOCKS5 proxy: `SERVER_IP:1080`

**iOS:**
- Use apps like "Shadowrocket" or "Surge"
- Configure SOCKS5 proxy

## Method 4: VPN Gateway Script

Create an automated script for Raspberry Pi:

```bash
#!/bin/bash
# vpn-gateway.sh

# Connect to VPN
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2 &

# Wait for connection
sleep 5

# Enable forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward

# Configure NAT
sudo iptables -t nat -A POSTROUTING -o omail0 -j MASQUERADE
sudo iptables -A FORWARD -i wlan0 -o omail0 -j ACCEPT
sudo iptables -A FORWARD -i omail0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT

echo "VPN Gateway ready! Connect phones to this WiFi."
```

## Comparison

| Method | Difficulty | Requirements | Best For |
|--------|-----------|--------------|----------|
| Router/Gateway | Easy | Raspberry Pi/Router | Home use |
| Termux (Android) | Medium | Root access | Advanced users |
| SOCKS5 Proxy | Medium | Proxy server | Workaround |
| Native App | Hard | Mobile development | Future |

## Recommended Approach

**For most users:** Use **Router/Gateway** method with Raspberry Pi.

**Advantages:**
- ✅ No root required on phone
- ✅ Works with all phones
- ✅ Automatic connection
- ✅ Easy to set up
- ✅ Can connect multiple devices

**Setup Time:** ~30 minutes

## Troubleshooting

### Phone Can't Connect to Router

1. Check WiFi password
2. Check router is broadcasting WiFi
3. Check phone WiFi settings

### No Internet on Phone

1. Check VPN is connected on router
2. Check NAT rules on router
3. Check IP forwarding enabled

### Slow Connection

1. Check server location
2. Check server resources
3. Check network speed

## Next Steps

1. Choose a method
2. Set up router/gateway (recommended)
3. Connect phone
4. Test connectivity
5. Enjoy secure browsing!

Your phone is now connected through VPN! 📱🔒
