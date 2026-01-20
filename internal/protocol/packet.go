package protocol

import (
	"encoding/binary"
	"errors"
	"net"
)

const (
	// PacketHeaderSize is the size of the packet header
	PacketHeaderSize = 8
	// MaxPacketSize is the maximum packet size (including header)
	MaxPacketSize = 65535
)

// PacketType represents the type of VPN packet
type PacketType uint8

const (
	// PacketTypeData is a data packet
	PacketTypeData PacketType = 0x01
	// PacketTypeKeepAlive is a keep-alive packet
	PacketTypeKeepAlive PacketType = 0x02
)

// PacketHeader is the header of a VPN packet
type PacketHeader struct {
	Type      PacketType
	Reserved  uint8
	Length    uint16
	SessionID uint32
}

// Packet represents a VPN packet
type Packet struct {
	Header PacketHeader
	Data   []byte
}

// Encode encodes a packet into bytes
func (p *Packet) Encode() []byte {
	buf := make([]byte, PacketHeaderSize+len(p.Data))
	
	buf[0] = byte(p.Header.Type)
	buf[1] = p.Header.Reserved
	binary.BigEndian.PutUint16(buf[2:4], p.Header.Length)
	binary.BigEndian.PutUint32(buf[4:8], p.Header.SessionID)
	
	copy(buf[PacketHeaderSize:], p.Data)
	
	return buf
}

// Decode decodes bytes into a packet
func Decode(data []byte) (*Packet, error) {
	if len(data) < PacketHeaderSize {
		return nil, errors.New("packet too short")
	}

	p := &Packet{}
	p.Header.Type = PacketType(data[0])
	p.Header.Reserved = data[1]
	p.Header.Length = binary.BigEndian.Uint16(data[2:4])
	p.Header.SessionID = binary.BigEndian.Uint32(data[4:8])
	
	if len(data) < PacketHeaderSize+int(p.Header.Length) {
		return nil, errors.New("packet length mismatch")
	}
	
	p.Data = make([]byte, p.Header.Length)
	copy(p.Data, data[PacketHeaderSize:PacketHeaderSize+int(p.Header.Length)])
	
	return p, nil
}

// NewDataPacket creates a new data packet
func NewDataPacket(sessionID uint32, data []byte) *Packet {
	return &Packet{
		Header: PacketHeader{
			Type:      PacketTypeData,
			Length:    uint16(len(data)),
			SessionID: sessionID,
		},
		Data: data,
	}
}

// NewKeepAlivePacket creates a new keep-alive packet
func NewKeepAlivePacket(sessionID uint32) *Packet {
	return &Packet{
		Header: PacketHeader{
			Type:      PacketTypeKeepAlive,
			Length:    0,
			SessionID: sessionID,
		},
		Data: nil,
	}
}

// IsIPv4 checks if the packet data is an IPv4 packet
func (p *Packet) IsIPv4() bool {
	if len(p.Data) < 1 {
		return false
	}
	return (p.Data[0] >> 4) == 4
}

// IsIPv6 checks if the packet data is an IPv6 packet
func (p *Packet) IsIPv6() bool {
	if len(p.Data) < 1 {
		return false
	}
	return (p.Data[0] >> 4) == 6
}

// GetDestinationIP extracts the destination IP from the packet (if IPv4/IPv6)
func (p *Packet) GetDestinationIP() (net.IP, error) {
	if p.IsIPv4() {
		if len(p.Data) < 20 {
			return nil, errors.New("packet too short for IPv4")
		}
		return net.IP(p.Data[16:20]), nil
	}
	if p.IsIPv6() {
		if len(p.Data) < 40 {
			return nil, errors.New("packet too short for IPv6")
		}
		return net.IP(p.Data[24:40]), nil
	}
	return nil, errors.New("not an IP packet")
}
