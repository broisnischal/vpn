# Mobile App Solutions

Solutions for using OMail VPN with native phone VPN settings.

## Current Status

**Native VPN settings don't work directly** because our VPN uses a custom protocol, not a standard one (OpenVPN, WireGuard, etc.).

## Solutions

### Solution 1: Android App (Recommended)

Create a native Android app using `VpnService` API.

**Features:**
- ✅ Works with native Android VPN settings
- ✅ No root required
- ✅ Easy to use
- ✅ Background service

**Implementation:**
- Use Android `VpnService` API
- Implement our custom protocol
- Create TUN interface through Android API

**Status:** Needs development

### Solution 2: iOS App

Create a native iOS app using `NetworkExtension` framework.

**Features:**
- ✅ Works with native iOS VPN settings
- ✅ No jailbreak required
- ✅ Easy to use

**Challenges:**
- Requires App Store approval or enterprise distribution
- More complex than Android

**Status:** Needs development

### Solution 3: OpenVPN Compatibility

Make our VPN compatible with OpenVPN protocol.

**How it works:**
- Modify server to speak OpenVPN protocol
- Use OpenVPN Connect app on phone
- Configure through app

**Status:** Partial implementation (see openvpn-wrapper.go)

### Solution 4: WireGuard Compatibility

Implement WireGuard protocol compatibility.

**Advantages:**
- Modern, fast protocol
- Native support on many devices
- Easy to use

**Status:** Needs implementation

### Solution 5: Router/Gateway (Current Best Option)

Set up VPN on router, connect phone to router's WiFi.

**Advantages:**
- ✅ Works immediately
- ✅ No app needed
- ✅ Works on all phones
- ✅ No root required

**How:**
1. Set up Raspberry Pi with VPN client
2. Configure Pi as WiFi gateway
3. Connect phone to Pi's WiFi
4. All traffic automatically goes through VPN

See `docs/MOBILE-CONNECTION.md` for details.

## Quick Comparison

| Solution | Difficulty | Time | Works Now? |
|----------|-----------|------|------------|
| Router/Gateway | Easy | 30 min | ✅ Yes |
| Android App | Medium | Days | ❌ No |
| iOS App | Hard | Days | ❌ No |
| OpenVPN Wrapper | Medium | Hours | ⚠️ Partial |
| WireGuard Compat | Hard | Days | ❌ No |

## Recommended Approach

**For immediate use:** Router/Gateway method

**For long-term:** Develop Android/iOS app

## Next Steps

1. **Now:** Use router/gateway method (see `docs/MOBILE-CONNECTION.md`)
2. **Future:** Develop mobile app for native VPN settings

## Contributing

Want to help develop the mobile app? See the code structure and contribute!
