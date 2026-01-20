# Deployment Scripts

Automated deployment scripts for cloud providers.

## AWS Deployment

```bash
# 1. Upload project to EC2
scp -r omail/ ubuntu@YOUR_EC2_IP:~/

# 2. SSH into server
ssh ubuntu@YOUR_EC2_IP

# 3. Run deployment script
cd omail
sudo bash deploy/aws-deploy.sh
```

**What it does:**
- Installs Docker
- Configures firewall (UFW)
- Enables IP forwarding
- Sets up NAT masquerading
- Deploys VPN server
- Generates secure password

## Azure Deployment

```bash
# 1. Upload project to VM
scp -r omail/ azureuser@YOUR_AZURE_IP:~/

# 2. SSH into VM
ssh azureuser@YOUR_AZURE_IP

# 3. Run deployment script
cd omail
sudo bash deploy/azure-deploy.sh
```

## Systemd Service

Create a systemd service for production:

```bash
sudo bash deploy/systemd-service.sh
```

This creates `/etc/systemd/system/omail-vpn.service` and enables it.

## Manual Deployment

See `docs/DEPLOYMENT.md` for manual step-by-step instructions.

## Requirements

- Ubuntu 20.04+ or Debian 11+
- Root/sudo access
- Internet connection
- Project files in `/opt/omail` or current directory

## After Deployment

1. **Save the password** (shown during deployment)
2. **Configure Security Group/NSG** to allow UDP 51820
3. **Connect client** using the server IP and password
4. **Test connectivity** with `ping 10.0.0.1`

## Troubleshooting

- **Check logs**: `docker-compose logs -f vpn-server`
- **Check status**: `docker-compose ps`
- **Restart**: `docker-compose restart vpn-server`
- **Stop**: `docker-compose down`
