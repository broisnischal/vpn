#!/bin/bash

# Quick VPN Connection Script
# Usage: sudo ./connect.sh [server-ip] [password]

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Error: Please run as root:${NC} sudo $0"
    exit 1
fi

# Get server IP
SERVER_IP="${1:-${VPN_SERVER}}"
if [ -z "$SERVER_IP" ]; then
    read -p "Enter server IP or hostname: " SERVER_IP
fi

# Get password
PASSWORD="${2:-${VPN_PASSWORD}}"
if [ -z "$PASSWORD" ]; then
    read -sp "Enter VPN password: " PASSWORD
    echo ""
fi

# Get client IP (auto-increment)
CLIENT_IP="${VPN_CLIENT_IP:-10.0.0.2}"

echo -e "${GREEN}Connecting to VPN server...${NC}"
echo "Server: $SERVER_IP:51820"
echo "Client IP: $CLIENT_IP"
echo ""

# Check if binary exists
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLIENT_BINARY="$SCRIPT_DIR/bin/omail-client"

if [ ! -f "$CLIENT_BINARY" ]; then
    echo -e "${RED}Error: Client binary not found at $CLIENT_BINARY${NC}"
    echo "Please build it first: make build"
    exit 1
fi

# Clean up any existing TUN interface
if ip link show omail0 &>/dev/null; then
    echo -e "${YELLOW}Cleaning up existing omail0 interface...${NC}"
    ip link delete omail0 2>/dev/null || true
fi

# Connect
echo -e "${GREEN}Starting VPN client...${NC}"
echo "Press Ctrl+C to disconnect"
echo ""

exec "$CLIENT_BINARY" \
    -server "$SERVER_IP:51820" \
    -password "$PASSWORD" \
    -tun-ip "$CLIENT_IP"
