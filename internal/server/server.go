package server

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
	"github.com/nees/omail/internal/tun"
)

// Server represents a VPN server
type Server struct {
	address    string
	crypto     *crypto.Crypto
	tun        *tun.Interface
	clients    map[uint32]*Client
	clientsMu  sync.RWMutex
	udpConn    *net.UDPConn
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// Client represents a connected VPN client
type Client struct {
	SessionID   uint32
	RemoteAddr  *net.UDPAddr
	LastSeen    time.Time
	mu          sync.Mutex
}

// Config holds server configuration
type Config struct {
	Address     string
	Password    string
	TUNName     string
	TUNIP       string
	TUNNetmask  string
	MTU         int
}

// NewServer creates a new VPN server
func NewServer(config Config) (*Server, error) {
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

	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		address: config.Address,
		crypto:  crypto,
		tun:     tunInterface,
		clients: make(map[uint32]*Client),
		ctx:     ctx,
		cancel:  cancel,
	}

	return s, nil
}

// Start starts the VPN server
func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.udpConn = conn

	log.Printf("VPN Server listening on %s", s.address)
	log.Printf("TUN interface: %s", s.tun.Name())

	// Start reading from TUN
	s.wg.Add(1)
	go s.readFromTUN()

	// Start reading from UDP
	s.wg.Add(1)
	go s.readFromUDP()

	// Start client cleanup goroutine
	s.wg.Add(1)
	go s.cleanupClients()

	return nil
}

// Stop stops the VPN server
func (s *Server) Stop() error {
	s.cancel()
	
	if s.udpConn != nil {
		s.udpConn.Close()
	}
	
	if s.tun != nil {
		s.tun.Down()
		s.tun.Close()
	}
	
	s.wg.Wait()
	return nil
}

// readFromTUN reads packets from TUN and forwards them to clients
func (s *Server) readFromTUN() {
	defer s.wg.Done()
	
	buf := make([]byte, 65535)
	
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			n, err := s.tun.Read(buf)
			if err != nil {
				log.Printf("Error reading from TUN: %v", err)
				continue
			}

			packet := buf[:n]
			
			// Determine which client to send to based on destination IP
			// For simplicity, we'll broadcast to all clients
			// In production, you'd want to maintain a routing table
			s.clientsMu.RLock()
			for _, client := range s.clients {
				s.sendToClient(client, packet)
			}
			s.clientsMu.RUnlock()
		}
	}
}

// readFromUDP reads packets from UDP and forwards them to TUN
func (s *Server) readFromUDP() {
	defer s.wg.Done()
	
	buf := make([]byte, 65535)
	
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			n, clientAddr, err := s.udpConn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Error reading from UDP: %v", err)
				continue
			}

			// Decrypt packet
			encrypted := buf[:n]
			decrypted, err := s.crypto.Decrypt(encrypted)
			if err != nil {
				log.Printf("Failed to decrypt packet from %s: %v", clientAddr, err)
				continue
			}

			// Decode protocol packet
			pkt, err := protocol.Decode(decrypted)
			if err != nil {
				log.Printf("Failed to decode packet: %v", err)
				continue
			}

			// Handle keep-alive
			if pkt.Header.Type == protocol.PacketTypeKeepAlive {
				s.handleKeepAlive(pkt.Header.SessionID, clientAddr)
				continue
			}

			// Handle data packet
			if pkt.Header.Type == protocol.PacketTypeData {
				s.handleDataPacket(pkt, clientAddr)
			}
		}
	}
}

// handleKeepAlive handles keep-alive packets
func (s *Server) handleKeepAlive(sessionID uint32, addr *net.UDPAddr) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	client, exists := s.clients[sessionID]
	if !exists {
		// Create new client
		client = &Client{
			SessionID:  sessionID,
			RemoteAddr: addr,
			LastSeen:   time.Now(),
		}
		s.clients[sessionID] = client
		log.Printf("New client connected: %s (session: %d)", addr, sessionID)
	} else {
		client.mu.Lock()
		client.LastSeen = time.Now()
		client.mu.Unlock()
	}
}

// handleDataPacket handles data packets from clients
func (s *Server) handleDataPacket(pkt *protocol.Packet, addr *net.UDPAddr) {
	// Update client last seen
	s.clientsMu.Lock()
	client, exists := s.clients[pkt.Header.SessionID]
	if !exists {
		client = &Client{
			SessionID:  pkt.Header.SessionID,
			RemoteAddr: addr,
			LastSeen:   time.Now(),
		}
		s.clients[pkt.Header.SessionID] = client
		log.Printf("New client connected: %s (session: %d)", addr, pkt.Header.SessionID)
	} else {
		client.mu.Lock()
		client.LastSeen = time.Now()
		client.mu.Unlock()
	}
	s.clientsMu.Unlock()

	// Write packet data to TUN
	if _, err := s.tun.Write(pkt.Data); err != nil {
		log.Printf("Error writing to TUN: %v", err)
	}
}

// sendToClient sends a packet to a client
func (s *Server) sendToClient(client *Client, data []byte) {
	client.mu.Lock()
	addr := client.RemoteAddr
	client.mu.Unlock()

	// Create protocol packet
	pkt := protocol.NewDataPacket(client.SessionID, data)
	
	// Encode packet
	encoded := pkt.Encode()
	
	// Encrypt packet
	encrypted, err := s.crypto.Encrypt(encoded)
	if err != nil {
		log.Printf("Failed to encrypt packet: %v", err)
		return
	}

	// Send to client
	if _, err := s.udpConn.WriteToUDP(encrypted, addr); err != nil {
		log.Printf("Error sending to client %s: %v", addr, err)
	}
}

// cleanupClients removes inactive clients
func (s *Server) cleanupClients() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.clientsMu.Lock()
			now := time.Now()
			for sessionID, client := range s.clients {
				client.mu.Lock()
				if now.Sub(client.LastSeen) > 60*time.Second {
					delete(s.clients, sessionID)
					log.Printf("Client disconnected: %s (session: %d)", client.RemoteAddr, sessionID)
				}
				client.mu.Unlock()
			}
			s.clientsMu.Unlock()
		}
	}
}

// generateSessionID generates a random session ID
func generateSessionID() uint32 {
	buf := make([]byte, 4)
	rand.Read(buf)
	return binary.BigEndian.Uint32(buf)
}
