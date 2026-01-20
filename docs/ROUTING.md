# Understanding Routing Tables in VPNs

This document explains how routing tables work and how our VPN manages them.

## What is a Routing Table?

A routing table is a data structure stored in the operating system kernel that determines where network packets should be sent. It's like a map that tells your computer how to reach different network destinations.

## Routing Table Structure

Each entry in a routing table typically contains:

- **Destination**: The network address (CIDR notation)
- **Gateway**: The next hop router (or interface)
- **Interface**: The network interface to use
- **Metric**: Priority/preference (lower = higher priority)

## Viewing Routing Tables

### Linux

```bash
# Modern ip command
ip route show

# Traditional route command
route -n

# Detailed view
ip route show table all
```

Example output:
```
default via 192.168.1.1 dev eth0
10.0.0.0/24 dev omail0
192.168.1.0/24 dev eth0
```

### macOS

```bash
# View routing table
netstat -rn

# Get default route
route -n get default

# View specific route
route -n get 10.0.0.0
```

## How VPN Routing Works

### 1. Full Tunnel Mode

In full tunnel mode, **all** internet traffic goes through the VPN:

```
Routing Table:
  0.0.0.0/0 → dev omail0 (VPN interface)
```

This means:
- All packets are captured by the TUN interface
- Packets are encrypted and sent to VPN server
- Server decrypts and forwards to internet
- Return traffic comes back through VPN

**Advantages**:
- Complete privacy (all traffic encrypted)
- Bypass geo-restrictions
- Hide real IP address

**Disadvantages**:
- Slower (all traffic goes through VPN)
- Uses VPN server bandwidth
- May break local network access

### 2. Split Tunnel Mode

In split tunnel mode, only **specific networks** go through VPN:

```
Routing Table:
  10.0.0.0/8 → dev omail0 (VPN interface)
  192.168.1.0/24 → dev omail0 (VPN interface)
  0.0.0.0/0 → dev eth0 via 192.168.1.1 (normal internet)
```

This means:
- Only packets to specified networks use VPN
- Other traffic uses normal internet connection
- Faster for non-VPN traffic
- Local network still accessible

**Advantages**:
- Faster (only VPN traffic goes through VPN)
- Local network access maintained
- Reduced VPN server load

**Disadvantages**:
- Some traffic not encrypted
- Real IP visible for non-VPN traffic

## Routing Operations

### Adding a Route

**Linux**:
```bash
# Add route through VPN interface
ip route add 10.0.0.0/8 dev omail0

# Add default route through VPN
ip route add default dev omail0
```

**macOS**:
```bash
# Add route through VPN interface
route add -net 10.0.0.0/8 -interface omail0

# Add default route
route add -net 0.0.0.0/0 -interface omail0
```

### Deleting a Route

**Linux**:
```bash
ip route del 10.0.0.0/8 dev omail0
ip route del default dev omail0
```

**macOS**:
```bash
route delete -net 10.0.0.0/8 -interface omail0
route delete -net 0.0.0.0/0 -interface omail0
```

## How Our VPN Manages Routes

### Client Connection

When a client connects:

1. **TUN Interface Created**: Virtual network interface created
2. **IP Assigned**: TUN interface gets IP address (e.g., 10.0.0.2)
3. **Interface Brought Up**: Interface activated
4. **Routes Added**: Routes added based on configuration:
   - Full tunnel: Default route (0.0.0.0/0)
   - Split tunnel: Specific network routes

### Client Disconnection

When a client disconnects:

1. **Routes Removed**: All VPN routes deleted
2. **Interface Brought Down**: TUN interface deactivated
3. **Interface Closed**: TUN interface destroyed

## Route Priority and Metrics

Routes have priorities (metrics). Lower metric = higher priority:

```
10.0.0.0/8 dev omail0 metric 100
10.0.0.0/8 dev eth0 metric 200
```

In this case, traffic to 10.0.0.0/8 will use `omail0` (lower metric).

## Common Routing Scenarios

### Scenario 1: Accessing VPN Network

```
Client wants to reach: 10.0.0.5
Routing table: 10.0.0.0/24 → dev omail0
Result: Packet goes through VPN
```

### Scenario 2: Accessing Internet

**Full Tunnel**:
```
Client wants to reach: 8.8.8.8
Routing table: 0.0.0.0/0 → dev omail0
Result: Packet goes through VPN
```

**Split Tunnel**:
```
Client wants to reach: 8.8.8.8
Routing table: 0.0.0.0/0 → dev eth0 via 192.168.1.1
Result: Packet uses normal internet
```

### Scenario 3: Local Network Access

**Full Tunnel** (problematic):
```
Client wants to reach: 192.168.1.100 (local printer)
Routing table: 0.0.0.0/0 → dev omail0
Result: Packet goes through VPN (may fail!)
```

**Split Tunnel** (better):
```
Client wants to reach: 192.168.1.100
Routing table: 192.168.1.0/24 → dev eth0
Result: Packet uses local network
```

## Debugging Routing Issues

### Check Current Routes

```bash
ip route show
```

### Test Route

```bash
# See which route will be used
ip route get 8.8.8.8
ip route get 10.0.0.5
```

### Monitor Traffic

```bash
# See packets on TUN interface
tcpdump -i omail0

# See routing decisions
ip route get 8.8.8.8
```

### Common Issues

1. **Route Not Added**: Check permissions (need root/NET_ADMIN)
2. **Wrong Interface**: Verify route points to correct interface
3. **Route Conflicts**: Check for duplicate routes with different metrics
4. **Interface Down**: Ensure TUN interface is up (`ip link show omail0`)

## Advanced: Policy-Based Routing

For more complex scenarios, you can use policy-based routing:

```bash
# Create custom routing table
echo "200 vpn" >> /etc/iproute2/rt_tables

# Add route to custom table
ip route add 10.0.0.0/8 dev omail0 table vpn

# Add rule to use custom table
ip rule add from 10.0.0.2 table vpn
```

This allows routing based on source IP, not just destination.

## Summary

- **Routing tables** determine where packets go
- **Full tunnel** routes all traffic through VPN
- **Split tunnel** routes only specific networks through VPN
- **Routes are added** when VPN connects
- **Routes are removed** when VPN disconnects
- **Check routes** with `ip route show` or `netstat -rn`

Understanding routing is crucial for VPN operation and troubleshooting!
