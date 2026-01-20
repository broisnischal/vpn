package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nees/omail/internal/client"
)

func main() {
	serverAddr := flag.String("server", "", "Server address (e.g., server.com:51820)")
	password := flag.String("password", "", "Encryption password (required)")
	tunName := flag.String("tun", "omail0", "TUN interface name")
	tunIP := flag.String("tun-ip", "10.0.0.2", "TUN interface IP address")
	tunNetmask := flag.String("tun-netmask", "255.255.255.0", "TUN interface netmask")
	mtu := flag.Int("mtu", 1500, "MTU size")
	splitTunnelStr := flag.String("split-tunnel", "", "Comma-separated list of CIDR networks for split tunneling (empty for full tunnel)")
	flag.Parse()

	if *serverAddr == "" {
		log.Fatal("Server address is required. Use -server flag")
	}
	if *password == "" {
		log.Fatal("Password is required. Use -password flag")
	}

	// Parse split tunnel networks
	var splitTunnel []*net.IPNet
	if *splitTunnelStr != "" {
		networks := strings.Split(*splitTunnelStr, ",")
		for _, netStr := range networks {
			netStr = strings.TrimSpace(netStr)
			_, ipNet, err := net.ParseCIDR(netStr)
			if err != nil {
				log.Fatalf("Invalid CIDR network %s: %v", netStr, err)
			}
			splitTunnel = append(splitTunnel, ipNet)
		}
	}

	config := client.Config{
		ServerAddr:  *serverAddr,
		Password:    *password,
		TUNName:     *tunName,
		TUNIP:       *tunIP,
		TUNNetmask:  *tunNetmask,
		MTU:         *mtu,
		SplitTunnel: splitTunnel,
	}

	cli, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	if err := cli.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Disconnecting from VPN server...")
	if err := cli.Disconnect(); err != nil {
		log.Printf("Error disconnecting: %v", err)
	}
}
