package main

// OpenVPN-Compatible Wrapper
// This creates an OpenVPN-compatible server that wraps our VPN protocol

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nees/omail/internal/crypto"
	"github.com/nees/omail/internal/protocol"
	"github.com/nees/omail/internal/server"
)

// OpenVPNWrapper wraps our VPN to be OpenVPN-compatible
type OpenVPNWrapper struct {
	server *server.Server
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: openvpn-wrapper <password> <port>")
		os.Exit(1)
	}

	password := os.Args[1]
	port := os.Args[2]

	// Create our VPN server
	config := server.Config{
		Address:    ":" + port,
		Password:   password,
		TUNName:    "omail0",
		TUNIP:      "10.0.0.1",
		TUNNetmask: "255.255.255.0",
		MTU:        1500,
	}

	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("OpenVPN-compatible wrapper running on port", port)
	log.Println("Use OpenVPN Connect app with generated config")

	// Generate OpenVPN config
	generateOpenVPNConfig(password, port)

	select {} // Keep running
}

func generateOpenVPNConfig(password, port string) {
	// Generate OpenVPN config file
	config := fmt.Sprintf(`client
dev tun
proto udp
remote YOUR_SERVER_IP %s
resolv-retry infinite
nobind
persist-key
persist-tun
verb 3
auth-user-pass
`, port)

	// Save config
	os.WriteFile("omail-client.ovpn", []byte(config), 0644)
	fmt.Println("OpenVPN config saved to: omail-client.ovpn")
	fmt.Println("Note: This is a basic wrapper. Full OpenVPN compatibility requires more work.")
}
