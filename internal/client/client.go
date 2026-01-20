package client

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/nees/omail/internal/crypto"
	"github.com/nees/omail/internal/protocol"
	"github.com/nees/omail/internal/routing"
	"github.com/nees/omail/internal/tun"
)

// Client represents a VPN client
type Client struct {
	serverAddr  string
	crypto      *crypto.Crypto
	tun         *tun.Interface
	udpConn     *net.UDPConn
	sessionID   uint32
	serverUDP   *net.UDPAddr
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	routing     *routing.Manager
	splitTunnel []*net.IPNet
}

// Config holds client configuration
type Config struct {
	ServerAddr  string
	Password    string
	TUNName     string
	TUNIP       string
	TUNNetmask  string
	MTU         int
	SplitTunnel []*net.IPNet // If empty, full tunnel
}

// NewClient creates a new VPN client
func NewClient(config Config) (*Client, error) {
	crypto, err := crypto.NewCrypto(config.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto: %w", err)
	}

	tunInterface, err := tun.New(config.TUNName, config.MTU)
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN interface: %w", err)
	}

	// Parse TUN IP and netmask
	tunIP := net.ParseIP(config.TUNIP)
	if tunIP == nil {
		return nil, fmt.Errorf("invalid TUN IP: %s", config.TUNIP)
	}

	var mask net.IPMask
	if config.TUNNetmask != "" {
		mask = net.IPMask(net.ParseIP(config.TUNNetmask).To4())
	} else {
		mask = net.CIDRMask(24, 32) // Default /24
	}

	// Set IP and bring interface up
	if err := tunInterface.SetIP(tunIP, mask); err != nil {
		tunInterface.Close()
		return nil, fmt.Errorf("failed to set TUN IP: %w", err)
	}

	if err := tunInterface.Up(); err != nil {
		tunInterface.Close()
		return nil, fmt.Errorf("failed to bring TUN up: %w", err)
	}

	// Resolve server address
	serverUDP, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	if err != nil {
		tunInterface.Close()
		return nil, fmt.Errorf("failed to resolve server address: %w", err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, serverUDP)
	if err != nil {
		tunInterface.Close()
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	// Generate session ID
	sessionID := generateSessionID()

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		serverAddr:  config.ServerAddr,
		crypto:      crypto,
		tun:         tunInterface,
		udpConn:     conn,
		sessionID:   sessionID,
		serverUDP:   serverUDP,
		ctx:         ctx,
		cancel:      cancel,
		routing:     routing.NewManager(config.TUNName),
		splitTunnel: config.SplitTunnel,
	}

	return client, nil
}

// Connect connects to the VPN server
func (c *Client) Connect() error {
	log.Printf("Connecting to VPN server at %s", c.serverAddr)
	log.Printf("TUN interface: %s", c.tun.Name())

	// Send initial keep-alive to establish session
	if err := c.sendKeepAlive(); err != nil {
		return fmt.Errorf("failed to establish session: %w", err)
	}

	// Setup routing
	if err := c.setupRouting(); err != nil {
		log.Printf("Warning: failed to setup routing: %v", err)
		// Continue anyway
	}

	// Start reading from TUN
	c.wg.Add(1)
	go c.readFromTUN()

	// Start reading from UDP
	c.wg.Add(1)
	go c.readFromUDP()

	// Start keep-alive goroutine
	c.wg.Add(1)
	go c.keepAlive()

	log.Println("Connected to VPN server")

	return nil
}

// Disconnect disconnects from the VPN server
func (c *Client) Disconnect() error {
	c.cancel()

	// Cleanup routing
	if err := c.routing.Cleanup(); err != nil {
		log.Printf("Warning: failed to cleanup routing: %v", err)
	}

	if c.udpConn != nil {
		c.udpConn.Close()
	}

	if c.tun != nil {
		c.tun.Down()
		c.tun.Close()
	}

	c.wg.Wait()
	return nil
}

// setupRouting sets up routing tables
func (c *Client) setupRouting() error {
	if len(c.splitTunnel) == 0 {
		// Full tunnel - route all traffic through VPN
		log.Println("Setting up full tunnel (all traffic through VPN)")
		return c.routing.SetupDefaultRoute()
	} else {
		// Split tunnel - only route specific networks
		log.Printf("Setting up split tunnel for %d networks", len(c.splitTunnel))
		return c.routing.SetupSplitTunnel(c.splitTunnel)
	}
}

// readFromTUN reads packets from TUN and sends them to server
func (c *Client) readFromTUN() {
	defer c.wg.Done()

	buf := make([]byte, 65535)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			n, err := c.tun.Read(buf)
			if err != nil {
				log.Printf("Error reading from TUN: %v", err)
				continue
			}

			packet := buf[:n]
			c.sendToServer(packet)
		}
	}
}

// readFromUDP reads packets from UDP and forwards them to TUN
func (c *Client) readFromUDP() {
	defer c.wg.Done()

	buf := make([]byte, 65535)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.udpConn.SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := c.udpConn.Read(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("Error reading from UDP: %v", err)
				continue
			}

			// Decrypt packet
			encrypted := buf[:n]
			decrypted, err := c.crypto.Decrypt(encrypted)
			if err != nil {
				log.Printf("Failed to decrypt packet: %v", err)
				continue
			}

			// Decode protocol packet
			pkt, err := protocol.Decode(decrypted)
			if err != nil {
				log.Printf("Failed to decode packet: %v", err)
				continue
			}

			// Handle data packet
			if pkt.Header.Type == protocol.PacketTypeData {
				// Write packet data to TUN
				if _, err := c.tun.Write(pkt.Data); err != nil {
					log.Printf("Error writing to TUN: %v", err)
				}
			}
		}
	}
}

// sendToServer sends a packet to the server
func (c *Client) sendToServer(data []byte) {
	// Create protocol packet
	pkt := protocol.NewDataPacket(c.sessionID, data)

	// Encode packet
	encoded := pkt.Encode()

	// Encrypt packet
	encrypted, err := c.crypto.Encrypt(encoded)
	if err != nil {
		log.Printf("Failed to encrypt packet: %v", err)
		return
	}

	// Send to server
	if _, err := c.udpConn.Write(encrypted); err != nil {
		log.Printf("Error sending to server: %v", err)
	}
}

// sendKeepAlive sends a keep-alive packet
func (c *Client) sendKeepAlive() error {
	pkt := protocol.NewKeepAlivePacket(c.sessionID)
	encoded := pkt.Encode()
	encrypted, err := c.crypto.Encrypt(encoded)
	if err != nil {
		return err
	}

	_, err = c.udpConn.Write(encrypted)
	return err
}

// keepAlive periodically sends keep-alive packets
func (c *Client) keepAlive() {
	defer c.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.sendKeepAlive(); err != nil {
				log.Printf("Failed to send keep-alive: %v", err)
			}
		}
	}
}

// generateSessionID generates a random session ID
func generateSessionID() uint32 {
	buf := make([]byte, 4)
	rand.Read(buf)
	return binary.BigEndian.Uint32(buf)
}
