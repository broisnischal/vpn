#!/bin/bash

# Azure Deployment Script for OMail VPN Server
# This script sets up the VPN server on an Azure VM

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== OMail VPN - Azure Deployment ===${NC}"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Error: This script must be run as root or with sudo${NC}"
    exit 1
fi

# Generate secure password
VPN_PASSWORD=$(openssl rand -base64 32)
echo -e "${YELLOW}Generated VPN Password:${NC} $VPN_PASSWORD"
echo "$VPN_PASSWORD" > /root/vpn-password.txt
chmod 600 /root/vpn-password.txt
echo -e "${GREEN}Password saved to /root/vpn-password.txt${NC}"
echo ""

# Update system
echo -e "${YELLOW}Updating system packages...${NC}"
apt update -qq
apt upgrade -y -qq

# Install Docker
echo -e "${YELLOW}Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    systemctl start docker
    systemctl enable docker
    rm get-docker.sh
fi

# Install Docker Compose
if ! command -v docker-compose &> /dev/null; then
    apt install -y docker-compose-plugin
fi

# Install other dependencies
echo -e "${YELLOW}Installing dependencies...${NC}"
apt install -y iptables-persistent netfilter-persistent

# Enable IP forwarding
echo -e "${YELLOW}Configuring IP forwarding...${NC}"
if ! grep -q "net.ipv4.ip_forward=1" /etc/sysctl.conf; then
    echo 'net.ipv4.ip_forward=1' >> /etc/sysctl.conf
fi
sysctl -p > /dev/null

# Configure firewall
echo -e "${YELLOW}Configuring firewall...${NC}"
ufw --force enable
ufw allow 22/tcp
ufw allow 51820/udp

# Get network interface (Azure usually uses eth0)
INTERFACE=$(ip route | grep default | awk '{print $5}' | head -1)
echo -e "${YELLOW}Detected network interface:${NC} $INTERFACE"

# Configure NAT masquerading
echo -e "${YELLOW}Configuring NAT...${NC}"
iptables -t nat -C POSTROUTING -o $INTERFACE -j MASQUERADE 2>/dev/null || \
    iptables -t nat -A POSTROUTING -o $INTERFACE -j MASQUERADE

iptables -C FORWARD -i omail0 -o $INTERFACE -j ACCEPT 2>/dev/null || \
    iptables -A FORWARD -i omail0 -o $INTERFACE -j ACCEPT

iptables -C FORWARD -i $INTERFACE -o omail0 -m state --state RELATED,ESTABLISHED -j ACCEPT 2>/dev/null || \
    iptables -A FORWARD -i $INTERFACE -o omail0 -m state --state RELATED,ESTABLISHED -j ACCEPT

# Save iptables rules
netfilter-persistent save

# Check if project directory exists
if [ ! -d "/opt/omail" ]; then
    echo -e "${YELLOW}Project directory not found. Please clone or copy the project to /opt/omail${NC}"
    echo "Example:"
    echo "  git clone YOUR_REPO /opt/omail"
    echo "  or"
    echo "  scp -r omail/ azureuser@SERVER:/opt/"
    exit 1
fi

# Deploy with Docker
echo -e "${YELLOW}Deploying VPN server with Docker...${NC}"
cd /opt/omail

# Create .env file
cat > .env << EOF
VPN_PASSWORD=$VPN_PASSWORD
VPN_ADDRESS=:51820
VPN_TUN_NAME=omail0
VPN_TUN_IP=10.0.0.1
VPN_TUN_NETMASK=255.255.255.0
VPN_MTU=1500
EOF

# Start server
docker-compose up -d vpn-server

# Wait for server to start
sleep 5

# Check status
if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}✓ VPN Server is running!${NC}"
    echo ""
    echo -e "${GREEN}=== Deployment Complete ===${NC}"
    echo ""
    echo "Server Information:"
    echo "  Public IP: $(curl -s ifconfig.me)"
    echo "  Port: 51820/UDP"
    echo "  Password: $VPN_PASSWORD"
    echo "  Password file: /root/vpn-password.txt"
    echo ""
    echo "View logs: docker-compose -f /opt/omail/docker-compose.yml logs -f vpn-server"
    echo "Stop server: docker-compose -f /opt/omail/docker-compose.yml down"
    echo ""
    echo -e "${YELLOW}Don't forget to:${NC}"
    echo "  1. Configure Azure NSG to allow UDP 51820"
    echo "  2. Save the password securely"
    echo "  3. Test connection from client"
else
    echo -e "${RED}✗ VPN Server failed to start${NC}"
    echo "Check logs: docker-compose logs vpn-server"
    exit 1
fi
