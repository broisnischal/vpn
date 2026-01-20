# Installation Guide

## Quick Install (Arch Linux)

```bash
# 1. Install Go
sudo pacman -S go

# 2. Run quick start script
cd /home/nees/omail
./quick-start.sh
```

## Manual Install

### Step 1: Install Go

**Arch Linux:**
```bash
sudo pacman -S go
```

**Verify:**
```bash
go version
```

### Step 2: Build VPN

```bash
cd /home/nees/omail
go mod download
make build
```

### Step 3: Test

```bash
# Automated test
sudo ./test-vpn.sh test

# Or manual test (see TUTORIAL.md)
```

## What Gets Installed

- `bin/omail-server` - VPN server binary
- `bin/omail-client` - VPN client binary

No system files are modified - everything runs from the project directory!

## Next Steps

After installation:
1. Read `TUTORIAL.md` to learn how to use it
2. Read `SETUP.md` for detailed setup instructions
3. Run `sudo ./test-vpn.sh test` to verify everything works
