#!/bin/bash

# Example script to run VPN client

set -e

# Configuration
SERVER="${VPN_SERVER:-localhost:51820}"
PASSWORD="${VPN_PASSWORD:-changeme123}"
TUN_NAME="${VPN_TUN_NAME:-omail0}"
TUN_IP="${VPN_CLIENT_TUN_IP:-10.0.0.2}"
TUN_NETMASK="${VPN_TUN_NETMASK:-255.255.255.0}"
MTU="${VPN_MTU:-1500}"
SPLIT_TUNNEL="${VPN_SPLIT_TUNNEL:-}"

echo "Starting OMail VPN Client"
echo "========================="
echo "Server: $SERVER"
echo "TUN Interface: $TUN_NAME"
echo "TUN IP: $TUN_IP"
echo "TUN Netmask: $TUN_NETMASK"
echo "MTU: $MTU"
if [ -n "$SPLIT_TUNNEL" ]; then
    echo "Split Tunnel: $SPLIT_TUNNEL"
else
    echo "Mode: Full Tunnel (all traffic)"
fi
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Error: This script must be run as root (for TUN interface creation)"
    echo "Usage: sudo $0"
    exit 1
fi

# Build if binary doesn't exist
if [ ! -f "./bin/omail-client" ]; then
    echo "Building client binary..."
    go build -o bin/omail-client ./cmd/client
fi

# Build command
CMD="./bin/omail-client \
    -server \"$SERVER\" \
    -password \"$PASSWORD\" \
    -tun \"$TUN_NAME\" \
    -tun-ip \"$TUN_IP\" \
    -tun-netmask \"$TUN_NETMASK\" \
    -mtu \"$MTU\""

# Add split tunnel if specified
if [ -n "$SPLIT_TUNNEL" ]; then
    CMD="$CMD -split-tunnel \"$SPLIT_TUNNEL\""
fi

# Run client
eval exec $CMD
