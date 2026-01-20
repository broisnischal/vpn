# Using Phone's Native VPN Settings

## The Problem

Your phone's **Settings → VPN** menu only supports **standard VPN protocols**:
- OpenVPN
- WireGuard  
- IPSec/IKEv2
- L2TP/IPSec
- PPTP

Our OMail VPN uses a **custom protocol**, so it doesn't appear in native VPN settings.

## Why This Happens

Native VPN settings require apps that implement standard protocols. When you add a VPN configuration:
- **Android**: Uses `VpnService` API (requires an app)
- **iOS**: Uses `NetworkExtension` framework (requires an app)

Our VPN needs a **mobile app** to work with native settings.

## Current Solutions

### ✅ Solution 1: Router/Gateway (Works Now!)

**Best option for immediate use.**

Set up VPN on a Raspberry Pi or router, then connect your phone to that device's WiFi.

**Advantages:**
- ✅ Works immediately
- ✅ No app needed
- ✅ Works on all phones (Android & iOS)
- ✅ No root/jailbreak required
- ✅ Automatic connection

**How it works:**
```
Your Phone → Router WiFi → Router VPN Client → VPN Server → Internet
```

**Setup:** See `docs/MOBILE-CONNECTION.md`

### ⚠️ Solution 2: OpenVPN Connect App (Partial)

Use OpenVPN Connect app with a compatibility wrapper.

**Steps:**
1. Install **OpenVPN Connect** app (free)
2. Generate OpenVPN config from our server
3. Import config into app
4. Connect through app

**Status:** Requires OpenVPN compatibility layer (in development)

### ❌ Solution 3: Custom Mobile App (Future)

Develop native Android/iOS app.

**Android App:**
- Uses `VpnService` API
- Implements our custom protocol
- Appears in VPN settings
- No root required

**iOS App:**
- Uses `NetworkExtension` framework
- Requires App Store or enterprise distribution
- More complex

**Status:** Needs development

## What You Can Do Now

### Option A: Router Method (Recommended)

1. **Get a Raspberry Pi** (~$35)
2. **Set up VPN client** on Pi
3. **Configure Pi as WiFi gateway**
4. **Connect phone** to Pi's WiFi
5. **Done!** All phone traffic goes through VPN

**Time:** ~30 minutes  
**Cost:** ~$35 (one-time)  
**Works:** ✅ Immediately

### Option B: Wait for Mobile App

We can develop a mobile app, but it takes time:
- Android: Several days of development
- iOS: More complex, requires App Store approval

**Time:** Days/weeks  
**Cost:** Free (if we develop it)  
**Works:** ❌ Not yet

### Option C: Use Laptop as Gateway

Connect laptop to VPN, share WiFi from laptop, connect phone to laptop's WiFi.

**Steps:**
1. Connect laptop to VPN
2. Enable WiFi hotspot on laptop
3. Connect phone to laptop's WiFi
4. Configure laptop to route phone traffic through VPN

**Time:** ~15 minutes  
**Cost:** Free  
**Works:** ✅ Yes

## Detailed: Router/Gateway Setup

### Using Raspberry Pi

**Step 1: Install VPN Client on Pi**

```bash
# SSH into Raspberry Pi
ssh pi@raspberry-pi-ip

# Install Go
sudo apt update
sudo apt install golang-go

# Clone project
cd ~
git clone YOUR_REPO_URL omail
cd omail

# Build client
go build -o bin/omail-client ./cmd/client
```

**Step 2: Connect Pi to VPN**

```bash
sudo ./bin/omail-client \
  -server YOUR_SERVER_IP:51820 \
  -password YOUR_PASSWORD \
  -tun-ip 10.0.0.2
```

**Step 3: Configure Pi as Gateway**

```bash
# Enable IP forwarding
echo 'net.ipv4.ip_forward=1' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# Configure NAT
sudo iptables -t nat -A POSTROUTING -o omail0 -j MASQUERADE
sudo iptables -A FORWARD -i wlan0 -o omail0 -j ACCEPT
sudo iptables -A FORWARD -i omail0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT
```

**Step 4: Set up WiFi Hotspot**

```bash
# Install hostapd and dnsmasq
sudo apt install hostapd dnsmasq

# Configure WiFi hotspot (see Raspberry Pi WiFi hotspot guides)
```

**Step 5: Connect Phone**

- Open phone WiFi settings
- Connect to Raspberry Pi's WiFi network
- All traffic automatically goes through VPN!

## Future: Mobile App Development

To make it work with native VPN settings, we need:

### Android App

```kotlin
// Pseudo-code
class OMailVpnService : VpnService() {
    fun connect() {
        val tun = Builder()
            .setSession("OMail VPN")
            .addAddress("10.0.0.2", 24)
            .addRoute("0.0.0.0", 0)
            .establish()
        
        // Implement our protocol
        connectToServer(tun)
    }
}
```

### iOS App

```swift
// Pseudo-code
class OMailVPNManager: NEVPNManager {
    func connect() {
        // Configure VPN
        // Implement our protocol
    }
}
```

## Summary

**Current Situation:**
- ❌ Native VPN settings don't work (custom protocol)
- ✅ Router/Gateway method works perfectly
- ⚠️ Mobile app needed for native settings

**Best Option Now:**
Use **Router/Gateway** method - it's the easiest and works immediately!

**Future:**
Develop mobile app for native VPN settings support.

## Quick Answer

**Q: Can I use my phone's VPN settings?**

**A:** Not directly yet. Use one of these:
1. **Router/Gateway** (easiest, works now)
2. **Laptop as gateway** (quick workaround)
3. **Wait for mobile app** (future)

See `docs/MOBILE-CONNECTION.md` for router setup guide!
