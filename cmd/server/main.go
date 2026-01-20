package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nees/omail/internal/server"
)

func main() {
	address := flag.String("address", ":51820", "Server listen address")
	password := flag.String("password", "", "Encryption password (required)")
	tunName := flag.String("tun", "omail0", "TUN interface name")
	tunIP := flag.String("tun-ip", "10.0.0.1", "TUN interface IP address")
	tunNetmask := flag.String("tun-netmask", "255.255.255.0", "TUN interface netmask")
	mtu := flag.Int("mtu", 1500, "MTU size")
	flag.Parse()

	if *password == "" {
		log.Fatal("Password is required. Use -password flag")
	}

	config := server.Config{
		Address:    *address,
		Password:   *password,
		TUNName:    *tunName,
		TUNIP:      *tunIP,
		TUNNetmask: *tunNetmask,
		MTU:        *mtu,
	}

	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	if err := srv.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
}
