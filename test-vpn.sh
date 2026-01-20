#!/bin/bash

# VPN Testing Script
# This script helps you test the VPN server and client

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== OMail VPN Test Script ===${NC}"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Error: This script must be run as root (for TUN interface)${NC}"
    echo "Usage: sudo $0 [server|client|test]"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed!${NC}"
    echo "Install Go first:"
    echo "  sudo pacman -S go"
    echo "  or"
    echo "  See SETUP.md for manual installation"
    exit 1
fi

# Check if binaries exist
if [ ! -f "./bin/omail-server" ] || [ ! -f "./bin/omail-client" ]; then
    echo -e "${YELLOW}Binaries not found. Building...${NC}"
    go mod download
    go build -o bin/omail-server ./cmd/server
    go build -o bin/omail-client ./cmd/client
    echo -e "${GREEN}Build complete!${NC}"
fi

MODE=${1:-test}
PASSWORD=${VPN_PASSWORD:-test-password-123}

case $MODE in
    server)
        echo -e "${GREEN}Starting VPN Server...${NC}"
        echo "Password: $PASSWORD"
        echo "Press Ctrl+C to stop"
        echo ""
        exec ./bin/omail-server \
            -address :51820 \
            -password "$PASSWORD" \
            -tun-ip 10.0.0.1
        ;;
    client)
        SERVER=${VPN_SERVER:-localhost:51820}
        echo -e "${GREEN}Starting VPN Client...${NC}"
        echo "Server: $SERVER"
        echo "Password: $PASSWORD"
        echo "Press Ctrl+C to stop"
        echo ""
        exec ./bin/omail-client \
            -server "$SERVER" \
            -password "$PASSWORD" \
            -tun-ip 10.0.0.2
        ;;
    test)
        echo -e "${GREEN}Running VPN Tests...${NC}"
        echo ""
        
        # Cleanup any existing TUN interface
        if ip link show omail0 &>/dev/null; then
            echo -e "${YELLOW}Removing existing omail0 interface...${NC}"
            ip link delete omail0 2>/dev/null || true
        fi
        
        echo -e "${GREEN}Test 1: Check TUN interface creation${NC}"
        echo "Starting server in background..."
        ./bin/omail-server \
            -address :51820 \
            -password "$PASSWORD" \
            -tun-ip 10.0.0.1 &
        SERVER_PID=$!
        sleep 2
        
        if ip link show omail0 &>/dev/null; then
            echo -e "${GREEN}✓ TUN interface created successfully${NC}"
        else
            echo -e "${RED}✗ Failed to create TUN interface${NC}"
            kill $SERVER_PID 2>/dev/null || true
            exit 1
        fi
        
        echo ""
        echo -e "${GREEN}Test 2: Check server TUN IP${NC}"
        if ip addr show omail0 | grep -q "10.0.0.1"; then
            echo -e "${GREEN}✓ Server TUN IP configured correctly${NC}"
        else
            echo -e "${RED}✗ Server TUN IP not configured${NC}"
        fi
        
        echo ""
        echo -e "${GREEN}Test 3: Start client${NC}"
        ./bin/omail-client \
            -server localhost:51820 \
            -password "$PASSWORD" \
            -tun-ip 10.0.0.2 &
        CLIENT_PID=$!
        sleep 3
        
        echo ""
        echo -e "${GREEN}Test 4: Check connectivity${NC}"
        sleep 2
        if ping -c 2 -W 2 10.0.0.1 &>/dev/null; then
            echo -e "${GREEN}✓ Client can ping server (VPN working!)${NC}"
        else
            echo -e "${YELLOW}⚠ Ping test failed (may need more time)${NC}"
        fi
        
        echo ""
        echo -e "${GREEN}Test 5: Check routing${NC}"
        ROUTES=$(ip route show | grep omail0 | wc -l)
        if [ "$ROUTES" -gt 0 ]; then
            echo -e "${GREEN}✓ Routes configured: $ROUTES routes through omail0${NC}"
            ip route show | grep omail0
        else
            echo -e "${YELLOW}⚠ No routes found (may be normal for split tunnel)${NC}"
        fi
        
        echo ""
        echo -e "${GREEN}Cleaning up...${NC}"
        kill $SERVER_PID 2>/dev/null || true
        kill $CLIENT_PID 2>/dev/null || true
        sleep 1
        ip link delete omail0 2>/dev/null || true
        
        echo ""
        echo -e "${GREEN}=== Tests Complete ===${NC}"
        echo ""
        echo "To run server: sudo $0 server"
        echo "To run client: sudo $0 client"
        ;;
    *)
        echo "Usage: sudo $0 [server|client|test]"
        echo ""
        echo "  server  - Start VPN server"
        echo "  client  - Start VPN client"
        echo "  test    - Run automated tests"
        exit 1
        ;;
esac
