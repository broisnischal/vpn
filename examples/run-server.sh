#!/bin/bash

# Example script to run VPN server

set -e

# Configuration
PASSWORD="${VPN_PASSWORD:-changeme123}"
ADDRESS="${VPN_ADDRESS:-:51820}"
TUN_NAME="${VPN_TUN_NAME:-omail0}"
TUN_IP="${VPN_TUN_IP:-10.0.0.1}"
TUN_NETMASK="${VPN_TUN_NETMASK:-255.255.255.0}"
MTU="${VPN_MTU:-1500}"

echo "Starting OMail VPN Server"
echo "========================="
echo "Address: $ADDRESS"
echo "TUN Interface: $TUN_NAME"
echo "TUN IP: $TUN_IP"
echo "TUN Netmask: $TUN_NETMASK"
echo "MTU: $MTU"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Error: This script must be run as root (for TUN interface creation)"
    echo "Usage: sudo $0"
    exit 1
fi

# Build if binary doesn't exist
if [ ! -f "./bin/omail-server" ]; then
    echo "Building server binary..."
    go build -o bin/omail-server ./cmd/server
fi

# Run server
exec ./bin/omail-server \
    -address "$ADDRESS" \
    -password "$PASSWORD" \
    -tun "$TUN_NAME" \
    -tun-ip "$TUN_IP" \
    -tun-netmask "$TUN_NETMASK" \
    -mtu "$MTU"
