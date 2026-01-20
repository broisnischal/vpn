# Cloud Deployment Guide

Deploy your OMail VPN server on AWS, Azure, or any cloud provider.

## Overview

When you deploy the VPN server in the cloud, you can:
- ✅ **Access Private Networks**: Connect to resources on the server's network
- ✅ **Bypass Geo-restrictions**: Access content restricted by location
- ✅ **Secure Browsing**: Encrypt all your internet traffic
- ✅ **Hide Your IP**: Your traffic appears to come from the server's IP

## Prerequisites

- Cloud account (AWS, Azure, or similar)
- SSH access to your cloud instance
- Docker installed on the server (recommended)
- Or Go installed for direct binary deployment

## AWS Deployment

### Step 1: Launch EC2 Instance

1. **Go to EC2 Console** → Launch Instance
2. **Choose AMI**: Ubuntu 22.04 LTS or Amazon Linux 2023
3. **Instance Type**: t2.micro (free tier) or t3.small (recommended)
4. **Configure Security Group**:
   - **SSH (22)**: Your IP only
   - **UDP (51820)**: 0.0.0.0/0 (all IPs) or specific IPs
5. **Launch** and save your key pair

### Step 2: Connect to Instance

```bash
ssh -i your-key.pem ubuntu@YOUR_EC2_IP
```

### Step 3: Install Docker (Recommended)

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Log out and back in for group changes
```

### Step 4: Deploy VPN Server

**Option A: Using Docker (Recommended)**

```bash
# Clone or upload your project
git clone YOUR_REPO_URL
cd omail

# Or upload files via SCP
# scp -r omail/ ubuntu@YOUR_EC2_IP:~/

# Set password
export VPN_PASSWORD=$(openssl rand -base64 32)
echo "VPN Password: $VPN_PASSWORD"

# Start server
docker-compose up -d vpn-server

# View logs
docker-compose logs -f vpn-server
```

**Option B: Direct Binary**

```bash
# Install Go
sudo apt update
sudo apt install -y golang-go

# Build and run
cd omail
go mod download
go build -o bin/omail-server ./cmd/server

# Run as service (using systemd)
sudo nano /etc/systemd/system/omail-vpn.service
```

### Step 5: Configure Firewall

```bash
# Allow UDP port 51820
sudo ufw allow 51820/udp
sudo ufw enable

# Or using AWS Security Group (recommended)
# Add inbound rule: UDP 51820 from 0.0.0.0/0
```

### Step 6: Enable IP Forwarding

```bash
# Enable IP forwarding
echo 'net.ipv4.ip_forward=1' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# Verify
sysctl net.ipv4.ip_forward
# Should output: net.ipv4.ip_forward = 1
```

### Step 7: Configure NAT (For Internet Access)

```bash
# Get your network interface name
ip route | grep default
# Usually: eth0 or ens5

# Enable NAT masquerading
sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
sudo iptables -A FORWARD -i omail0 -o eth0 -j ACCEPT
sudo iptables -A FORWARD -i eth0 -o omail0 -m state --state RELATED,ESTABLISHED -j ACCEPT

# Save iptables rules (Ubuntu)
sudo apt install iptables-persistent
sudo netfilter-persistent save
```

## Azure Deployment

### Step 1: Create Virtual Machine

1. **Azure Portal** → Create Resource → Virtual Machine
2. **Image**: Ubuntu Server 22.04 LTS
3. **Size**: Standard_B1s (free tier) or Standard_B2s
4. **Networking**:
   - Create new Network Security Group
   - Add inbound rule: UDP 51820 (allow all or specific IPs)
   - Add inbound rule: SSH 22 (your IP only)
5. **Create** and note the public IP

### Step 2: Connect to VM

```bash
ssh azureuser@YOUR_AZURE_IP
```

### Step 3: Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt install docker-compose-plugin

# Log out and back in
```

### Step 4: Deploy VPN Server

Same as AWS Step 4 above.

### Step 5: Configure Azure Network Security Group

1. **Azure Portal** → Your VM → Networking
2. **Add inbound port rule**:
   - Port: 51820
   - Protocol: UDP
   - Source: Any or specific IPs
   - Action: Allow

### Step 6: Enable IP Forwarding

```bash
# Same as AWS Step 6
echo 'net.ipv4.ip_forward=1' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

### Step 7: Configure NAT

Same as AWS Step 7.

## Generic Cloud Deployment (Any Provider)

### Quick Setup Script

```bash
#!/bin/bash
# deploy-vpn.sh

set -e

echo "Deploying OMail VPN Server..."

# Install dependencies
sudo apt update
sudo apt install -y docker.io docker-compose iptables-persistent

# Enable IP forwarding
echo 'net.ipv4.ip_forward=1' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# Generate password
VPN_PASSWORD=$(openssl rand -base64 32)
echo "VPN Password: $VPN_PASSWORD"
echo "$VPN_PASSWORD" > /tmp/vpn-password.txt

# Configure firewall
sudo ufw allow 51820/udp
sudo ufw --force enable

# Deploy with Docker
cd /path/to/omail
VPN_PASSWORD=$VPN_PASSWORD docker-compose up -d vpn-server

echo "VPN Server deployed!"
echo "Password saved in /tmp/vpn-password.txt"
```

## Connecting Clients

### From Your Local Machine

```bash
# Get server IP
SERVER_IP="YOUR_CLOUD_SERVER_IP"

# Connect client
sudo ./bin/omail-client \
  -server $SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

### From Another Cloud Instance

Same as above, just use the server's IP address.

## What Can You Access?

### 1. Private Network Resources

If your VPN server is on a private network (VPC), you can access:
- Internal databases
- Private APIs
- Internal services
- Other VPC resources

**Example:**
```
Your Computer → VPN → AWS VPC → Private RDS Database
```

### 2. Geo-Restricted Content

Route traffic through the server's location:
- Access region-locked services
- Bypass content restrictions
- Use server's IP address

**Example:**
```
Your Computer → VPN → US Server → Access US-only content
```

### 3. Secure Browsing

All traffic encrypted:
- Hide browsing from ISP
- Protect on public WiFi
- Encrypt all connections

### 4. Server's Network

Access resources on the server's local network:
- Local services
- Network shares
- Internal tools

## Security Considerations

### 1. Strong Password

```bash
# Generate strong password
openssl rand -base64 32
```

### 2. Firewall Rules

**Restrict Access:**
```bash
# Only allow specific IPs
sudo ufw delete allow 51820/udp
sudo ufw allow from YOUR_IP to any port 51820 proto udp
```

### 3. Regular Updates

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Update Docker images
docker-compose pull
docker-compose up -d
```

### 4. Monitor Connections

```bash
# Check active connections
sudo netstat -unap | grep 51820

# Monitor traffic
sudo tcpdump -i omail0
```

## Troubleshooting

### Can't Connect from Client

1. **Check Firewall**: Ensure UDP 51820 is open
2. **Check Server**: `sudo docker-compose logs vpn-server`
3. **Check Routing**: Verify IP forwarding is enabled
4. **Test Connectivity**: `nc -u -v SERVER_IP 51820`

### No Internet Through VPN

1. **Check NAT**: Verify iptables rules
2. **Check Routing**: `ip route show`
3. **Check DNS**: Try `nslookup google.com`
4. **Check Server Internet**: `curl ifconfig.me` on server

### Connection Drops

1. **Check Keep-Alive**: Client sends every 10 seconds
2. **Check Timeout**: Server timeout is 60 seconds
3. **Check Logs**: `docker-compose logs -f vpn-server`

## Cost Estimation

### AWS
- **t2.micro**: Free tier (750 hours/month)
- **t3.small**: ~$15/month
- **Data Transfer**: First 1GB free, then ~$0.09/GB

### Azure
- **B1s**: Free tier (750 hours/month)
- **B2s**: ~$30/month
- **Data Transfer**: First 5GB free, then ~$0.087/GB

## Production Recommendations

1. **Use Systemd Service**: Run as a service, not manually
2. **Logging**: Set up log rotation
3. **Monitoring**: Monitor server health
4. **Backup**: Backup configuration
5. **SSL/TLS**: Consider adding TLS layer
6. **Rate Limiting**: Prevent abuse
7. **Multiple Servers**: Load balancing

## Next Steps

1. Deploy server on your chosen cloud provider
2. Configure firewall and networking
3. Connect clients
4. Test connectivity
5. Monitor usage

Your VPN server is now accessible from anywhere! 🌍
