#!/bin/bash

# Quick Start Script for OMail VPN
# This script guides you through installation and testing

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     OMail VPN - Quick Start Guide     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Step 1: Check Go installation
echo -e "${YELLOW}Step 1: Checking Go installation...${NC}"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo -e "${GREEN}✓ Go is installed: $GO_VERSION${NC}"
else
    echo -e "${RED}✗ Go is not installed${NC}"
    echo ""
    echo "Please install Go first:"
    echo ""
    echo -e "${BLUE}For Arch Linux:${NC}"
    echo "  sudo pacman -S go"
    echo ""
    echo -e "${BLUE}For Ubuntu/Debian:${NC}"
    echo "  sudo apt update && sudo apt install golang-go"
    echo ""
    echo -e "${BLUE}For macOS:${NC}"
    echo "  brew install go"
    echo ""
    echo "After installing, run this script again."
    exit 1
fi

# Step 2: Install dependencies
echo ""
echo -e "${YELLOW}Step 2: Installing dependencies...${NC}"
cd "$(dirname "$0")"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Step 3: Build binaries
echo ""
echo -e "${YELLOW}Step 3: Building VPN binaries...${NC}"
mkdir -p bin
go build -o bin/omail-server ./cmd/server
go build -o bin/omail-client ./cmd/client
echo -e "${GREEN}✓ Build complete!${NC}"
echo ""
echo "Binaries created:"
echo "  - bin/omail-server"
echo "  - bin/omail-client"

# Step 4: Instructions
echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║        Ready to Test!                 ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}To test the VPN, open TWO terminals:${NC}"
echo ""
echo -e "${YELLOW}Terminal 1 - Start Server:${NC}"
echo "  cd $(pwd)"
echo "  sudo ./bin/omail-server \\"
echo "    -address :51820 \\"
echo "    -password test-password-123 \\"
echo "    -tun-ip 10.0.0.1"
echo ""
echo -e "${YELLOW}Terminal 2 - Start Client:${NC}"
echo "  cd $(pwd)"
echo "  sudo ./bin/omail-client \\"
echo "    -server localhost:51820 \\"
echo "    -password test-password-123 \\"
echo "    -tun-ip 10.0.0.2"
echo ""
echo -e "${YELLOW}Terminal 3 - Test Connection:${NC}"
echo "  ping 10.0.0.1"
echo "  ip addr show omail0"
echo "  ip route show"
echo ""
echo -e "${BLUE}Or use the automated test script:${NC}"
echo "  sudo ./test-vpn.sh test"
echo ""
echo -e "${BLUE}For detailed instructions, see:${NC}"
echo "  - TUTORIAL.md (complete guide)"
echo "  - SETUP.md (setup instructions)"
echo "  - README.md (full documentation)"
echo ""
