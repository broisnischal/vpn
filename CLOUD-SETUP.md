# Quick Cloud Setup Guide

## 🚀 Deploy to AWS in 5 Minutes

### Step 1: Launch EC2 Instance

1. Go to **AWS Console** → **EC2** → **Launch Instance**
2. Choose **Ubuntu 22.04 LTS**
3. Select **t2.micro** (free tier) or **t3.small**
4. **Security Group**:
   - SSH (22) from your IP
   - **UDP (51820) from 0.0.0.0/0** ← Important!
5. Launch and connect: `ssh -i key.pem ubuntu@YOUR_IP`

### Step 2: Run Deployment Script

```bash
# On your local machine, upload project
scp -r omail/ ubuntu@YOUR_EC2_IP:~/

# SSH into server
ssh ubuntu@YOUR_EC2_IP

# Run deployment
cd omail
sudo bash deploy/aws-deploy.sh
```

**That's it!** The script will:
- ✅ Install Docker
- ✅ Configure firewall
- ✅ Set up NAT
- ✅ Deploy VPN server
- ✅ Generate secure password

### Step 3: Connect from Your Computer

```bash
# Get the password (shown during deployment)
# Or: ssh ubuntu@YOUR_EC2_IP "sudo cat /root/vpn-password.txt"

# Connect client
sudo ./bin/omail-client \
  -server YOUR_EC2_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

## 🚀 Deploy to Azure in 5 Minutes

### Step 1: Create VM

1. **Azure Portal** → **Create Resource** → **Virtual Machine**
2. **Ubuntu Server 22.04 LTS**
3. **Size**: Standard_B1s (free tier)
4. **Networking** → **Add inbound port rule**:
   - UDP 51820 (allow all)
5. Create and connect: `ssh azureuser@YOUR_IP`

### Step 2: Run Deployment Script

```bash
# Upload project
scp -r omail/ azureuser@YOUR_IP:~/

# SSH and deploy
ssh azureuser@YOUR_IP
cd omail
sudo bash deploy/azure-deploy.sh
```

### Step 3: Connect Client

Same as AWS Step 3 above.

## What Can You Access?

### ✅ Private Network Resources

If your server is in a VPC/VNet, you can access:
- Private databases (RDS, Azure Database)
- Internal APIs
- Private services
- Other VPC resources

**Example:**
```
Your Computer → VPN → AWS VPC → Private RDS Database ✅
```

### ✅ Geo-Restricted Content

Access content based on server location:
- Streaming services
- Region-locked websites
- Country-specific content

**Example:**
```
Your Computer (EU) → VPN → US Server → Access US Netflix ✅
```

### ✅ Secure Browsing

All traffic encrypted:
- Hide from ISP
- Protect on public WiFi
- Encrypt all connections

### ✅ Server's Network

Access resources on server's network:
- Local services
- Network shares
- Internal tools

## Security Checklist

- [ ] Strong password generated
- [ ] Firewall configured (UDP 51820)
- [ ] Security Group/NSG rules set
- [ ] IP forwarding enabled
- [ ] NAT configured
- [ ] Password saved securely

## Troubleshooting

### Can't Connect

1. **Check Security Group**: UDP 51820 must be open
2. **Check Firewall**: `sudo ufw status`
3. **Check Server**: `docker-compose logs vpn-server`
4. **Test Port**: `nc -u -v SERVER_IP 51820`

### No Internet Through VPN

1. **Check NAT**: `sudo iptables -t nat -L`
2. **Check Forwarding**: `sysctl net.ipv4.ip_forward`
3. **Check Routes**: `ip route show`

## Cost Estimate

- **AWS t2.micro**: Free (750 hrs/month)
- **Azure B1s**: Free (750 hrs/month)
- **Data Transfer**: ~$0.09/GB after free tier

## Next Steps

1. ✅ Deploy server
2. ✅ Connect client
3. ✅ Test connectivity
4. ✅ Access protected resources!

Your VPN is now live in the cloud! 🌍
