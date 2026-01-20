#!/bin/bash

# Create systemd service for OMail VPN Server
# This allows the VPN to run as a system service

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Creating Systemd Service ===${NC}"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Error: This script must be run as root"
    exit 1
fi

# Get VPN password
read -sp "Enter VPN password: " VPN_PASSWORD
echo ""

# Get server IP (optional, for binding)
read -p "Server IP (leave empty for all interfaces): " SERVER_IP
SERVER_IP=${SERVER_IP:-""}
ADDRESS=${SERVER_IP:+$SERVER_IP:51820}
ADDRESS=${ADDRESS:-":51820"}

# Get binary path
BINARY_PATH=$(readlink -f $(which omail-server 2>/dev/null) || echo "/opt/omail/bin/omail-server")
read -p "Binary path [$BINARY_PATH]: " INPUT_BINARY
BINARY_PATH=${INPUT_BINARY:-$BINARY_PATH}

if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: Binary not found at $BINARY_PATH"
    exit 1
fi

# Create service file
SERVICE_FILE="/etc/systemd/system/omail-vpn.service"

cat > $SERVICE_FILE << EOF
[Unit]
Description=OMail VPN Server
After=network.target

[Service]
Type=simple
ExecStart=$BINARY_PATH \\
    -address $ADDRESS \\
    -password $VPN_PASSWORD \\
    -tun omail0 \\
    -tun-ip 10.0.0.1 \\
    -tun-netmask 255.255.255.0 \\
    -mtu 1500
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# Security
CapabilityBoundingSet=CAP_NET_ADMIN CAP_SYS_MODULE
AmbientCapabilities=CAP_NET_ADMIN CAP_SYS_MODULE
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

echo -e "${GREEN}Service file created: $SERVICE_FILE${NC}"

# Reload systemd
systemctl daemon-reload

# Enable service
systemctl enable omail-vpn.service

echo -e "${GREEN}Service enabled${NC}"
echo ""
echo "Commands:"
echo "  Start:   sudo systemctl start omail-vpn"
echo "  Stop:    sudo systemctl stop omail-vpn"
echo "  Status:  sudo systemctl status omail-vpn"
echo "  Logs:    sudo journalctl -u omail-vpn -f"
echo ""
read -p "Start service now? (y/n): " START_NOW

if [ "$START_NOW" = "y" ]; then
    systemctl start omail-vpn.service
    sleep 2
    systemctl status omail-vpn.service
fi
