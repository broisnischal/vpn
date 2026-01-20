# Quick Connect Guide

Fastest way to connect from your laptop or phone.

## From Laptop (Linux/macOS)

### Quick Method

```bash
# Make script executable (first time only)
chmod +x connect.sh

# Connect
sudo ./connect.sh YOUR_SERVER_IP YOUR_PASSWORD
```

### Manual Method

```bash
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

### With Environment Variables

```bash
export VPN_SERVER=54.123.45.67
export VPN_PASSWORD=my-password

sudo ./connect.sh
```

## From Phone

### Android (Termux - Requires Root)

```bash
# In Termux
pkg install golang git
git clone YOUR_REPO_URL
cd omail
go build -o bin/omail-client ./cmd/client

# Connect (requires root)
su
./bin/omail-client -server SERVER_IP:51820 -password PASSWORD
```

### iOS / Android (Easiest - Router Method)

**Set up VPN on Raspberry Pi, then connect phone to Pi's WiFi.**

See `docs/MOBILE-CONNECTION.md` for details.

## Verify Connection

```bash
# Check TUN interface
ip addr show omail0

# Ping server
ping 10.0.0.1

# Check your IP (should be server's IP)
curl ifconfig.me
```

## Disconnect

Press `Ctrl+C` in the terminal where client is running.

## Examples

### Connect to AWS Server

```bash
sudo ./connect.sh 54.123.45.67 my-secure-password
```

### Connect to Azure Server

```bash
sudo ./connect.sh 20.123.45.67 my-secure-password
```

### Split Tunnel (Only Route Specific Networks)

```bash
sudo ./bin/omail-client \
  -server 54.123.45.67:51820 \
  -password my-password \
  -split-tunnel "10.0.0.0/8,192.168.1.0/24"
```

That's it! You're connected! 🎉
