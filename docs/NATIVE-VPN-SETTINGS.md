# Using Native VPN Settings on Phone

## Why Native VPN Settings Don't Work (Yet)

Your phone's built-in VPN settings support **standard protocols** like:
- **OpenVPN** (most common)
- **WireGuard** (modern, fast)
- **IPSec/IKEv2**
- **L2TP/IPSec**
- **PPTP**

Our OMail VPN uses a **custom protocol** (UDP-based with our own packet format), so it's not directly compatible with native VPN settings.

## Solutions

### Option 1: Create Mobile App (Best Long-term Solution)

Create a native Android/iOS app that uses the VPN API.

**Android:**
- Uses `VpnService` API
- Can create TUN interface
- Works with our custom protocol

**iOS:**
- Uses `NetworkExtension` framework
- Requires app store or enterprise distribution
- More complex

### Option 2: Add OpenVPN Compatibility (Easier)

Modify our VPN to be OpenVPN-compatible, then use OpenVPN Connect app.

### Option 3: Add WireGuard Compatibility (Modern)

Implement WireGuard protocol compatibility.

### Option 4: Use Router/Gateway (Current Best Option)

Set up VPN on router, connect phone to router's WiFi.

## Quick Solution: OpenVPN-Compatible Wrapper

Let me create a wrapper that makes our VPN work with OpenVPN Connect app.
