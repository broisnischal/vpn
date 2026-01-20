# Can I Use My Phone's VPN Settings? 🤔

## Short Answer

**Not directly yet**, but there are easy workarounds!

## Why Native VPN Settings Don't Work

Your phone's **Settings → VPN** menu only works with **standard VPN protocols**:
- ✅ OpenVPN
- ✅ WireGuard
- ✅ IPSec/IKEv2

Our OMail VPN uses a **custom protocol**, so it needs either:
1. A mobile app (we can build this)
2. A workaround (use router/gateway)

## ✅ Easy Solution: Router/Gateway Method

**This works RIGHT NOW and is the easiest option!**

### How It Works

```
Your Phone → Router WiFi → Router VPN → VPN Server → Internet
```

Instead of connecting phone directly to VPN, you:
1. Set up VPN on a Raspberry Pi (or router)
2. Connect your phone to the Pi's WiFi
3. **Done!** All phone traffic goes through VPN automatically

### Setup (30 minutes)

**Step 1:** Get a Raspberry Pi (~$35)

**Step 2:** Set up VPN client on Pi:
```bash
# SSH into Pi
ssh pi@raspberry-pi-ip

# Install and build
sudo apt install golang-go
cd omail
go build -o bin/omail-client ./cmd/client

# Connect to VPN
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD
```

**Step 3:** Configure Pi as WiFi gateway (see `docs/MOBILE-CONNECTION.md`)

**Step 4:** Connect phone to Pi's WiFi - **that's it!**

### Advantages

- ✅ Works immediately
- ✅ No app needed
- ✅ Works on Android AND iOS
- ✅ No root/jailbreak
- ✅ Automatic - just connect to WiFi

## 🔮 Future: Mobile App

To use native VPN settings, we need to build a mobile app:

### Android App
- Uses `VpnService` API
- Appears in VPN settings
- No root required
- **Status:** Can be developed

### iOS App  
- Uses `NetworkExtension` framework
- Appears in VPN settings
- Requires App Store approval
- **Status:** More complex

**Would you like us to develop a mobile app?** It takes a few days but then you can use native VPN settings!

## 🚀 Quick Alternatives

### Option 1: Laptop as Gateway

1. Connect laptop to VPN
2. Enable WiFi hotspot on laptop
3. Connect phone to laptop's WiFi
4. Phone traffic goes through VPN!

### Option 2: Use OpenVPN App

We can add OpenVPN compatibility, then use **OpenVPN Connect** app:
1. Install OpenVPN Connect (free app)
2. Import our VPN config
3. Connect through app

**Status:** Requires compatibility layer (in development)

## Comparison

| Method | Works Now? | Difficulty | Cost |
|--------|-----------|------------|------|
| Router/Gateway | ✅ Yes | Easy | $35 (Pi) |
| Laptop Gateway | ✅ Yes | Easy | Free |
| Mobile App | ❌ No | Medium | Free (if we build) |
| OpenVPN App | ⚠️ Partial | Medium | Free |

## Recommendation

**For immediate use:** Use **Router/Gateway** method
- Works right now
- Easy to set up
- Best user experience

**For future:** Develop mobile app for native VPN settings

## Next Steps

1. **Now:** Set up router/gateway (see `docs/MOBILE-CONNECTION.md`)
2. **Future:** We can develop mobile app if you want

**Want help setting up the router method?** I can guide you through it step by step!
