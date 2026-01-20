package tun

import (
	"fmt"
	"io"
	"net"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/songgao/water"
	"golang.org/x/sys/unix"
)

const (
	// DefaultMTU is the default Maximum Transmission Unit
	DefaultMTU = 1500
	// TUNInterfaceName is the name of the TUN interface
	TUNInterfaceName = "omail0"
)

// Interface represents a TUN network interface
type Interface struct {
	ifce *water.Interface
	name string
	mtu  int
}

// New creates a new TUN interface
func New(name string, mtu int) (*Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	ifce, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN interface: %w", err)
	}

	tun := &Interface{
		ifce: ifce,
		name: name,
		mtu:  mtu,
	}

	// Set MTU
	if err := tun.SetMTU(mtu); err != nil {
		ifce.Close()
		return nil, fmt.Errorf("failed to set MTU: %w", err)
	}

	return tun, nil
}

// Name returns the interface name
func (t *Interface) Name() string {
	return t.name
}

// MTU returns the MTU
func (t *Interface) MTU() int {
	return t.mtu
}

// Read reads a packet from the TUN interface
func (t *Interface) Read(p []byte) (n int, err error) {
	return t.ifce.Read(p)
}

// Write writes a packet to the TUN interface
func (t *Interface) Write(p []byte) (n int, err error) {
	return t.ifce.Write(p)
}

// Close closes the TUN interface
func (t *Interface) Close() error {
	return t.ifce.Close()
}

// SetIP sets the IP address and netmask for the interface
func (t *Interface) SetIP(ip net.IP, mask net.IPMask) error {
	return setIP(t.name, ip, mask)
}

// Up brings the interface up
func (t *Interface) Up() error {
	return up(t.name)
}

// Down brings the interface down
func (t *Interface) Down() error {
	return down(t.name)
}

// AddRoute adds a route through this interface
func (t *Interface) AddRoute(dest *net.IPNet) error {
	return addRoute(t.name, dest)
}

// setIP sets the IP address on the interface (platform-specific)
func setIP(name string, ip net.IP, mask net.IPMask) error {
	switch runtime.GOOS {
	case "linux":
		return setIPLinux(name, ip, mask)
	case "darwin":
		return setIPDarwin(name, ip, mask)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func setIPLinux(name string, ip net.IP, mask net.IPMask) error {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	ifreq, err := unix.NewIfreq(name)
	if err != nil {
		return err
	}

	// Set IP address
	addr := unix.SockaddrInet4{}
	copy(addr.Addr[:], ip.To4())
	ifreq.SetInet4Addr(addr.Addr)

	if err := unix.IoctlIfreq(fd, unix.SIOCSIFADDR, ifreq); err != nil {
		return err
	}

	// Set netmask
	maskBytes := mask
	if len(maskBytes) == 16 {
		maskBytes = maskBytes[12:]
	}
	copy(addr.Addr[:], maskBytes)
	ifreq.SetInet4Addr(addr.Addr)

	if err := unix.IoctlIfreq(fd, unix.SIOCSIFNETMASK, ifreq); err != nil {
		return err
	}

	return nil
}

func setIPDarwin(name string, ip net.IP, mask net.IPMask) error {
	// Use ifconfig command on macOS
	ipStr := ip.String()
	maskStr := net.IP(mask).String()
	
	cmd := exec.Command("ifconfig", name, "inet", ipStr, "netmask", maskStr)
	return cmd.Run()
}

// up brings the interface up
func up(name string) error {
	switch runtime.GOOS {
	case "linux":
		return upLinux(name)
	case "darwin":
		return upDarwin(name)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func upLinux(name string) error {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	ifreq, err := unix.NewIfreq(name)
	if err != nil {
		return err
	}

	// Get current flags
	if err := unix.IoctlIfreq(fd, unix.SIOCGIFFLAGS, ifreq); err != nil {
		return err
	}

	flags := ifreq.Uint16()
	flags |= unix.IFF_UP | unix.IFF_RUNNING
	ifreq.SetUint16(flags)

	return unix.IoctlIfreq(fd, unix.SIOCSIFFLAGS, ifreq)
}

func upDarwin(name string) error {
	cmd := exec.Command("ifconfig", name, "up")
	return cmd.Run()
}

// down brings the interface down
func down(name string) error {
	switch runtime.GOOS {
	case "linux":
		return downLinux(name)
	case "darwin":
		return downDarwin(name)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func downLinux(name string) error {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	ifreq, err := unix.NewIfreq(name)
	if err != nil {
		return err
	}

	if err := unix.IoctlIfreq(fd, unix.SIOCGIFFLAGS, ifreq); err != nil {
		return err
	}

	flags := ifreq.Uint16()
	flags &^= unix.IFF_UP
	ifreq.SetUint16(flags)

	return unix.IoctlIfreq(fd, unix.SIOCSIFFLAGS, ifreq)
}

func downDarwin(name string) error {
	cmd := exec.Command("ifconfig", name, "down")
	return cmd.Run()
}

// SetMTU sets the MTU for the interface
func (t *Interface) SetMTU(mtu int) error {
	switch runtime.GOOS {
	case "linux":
		return setMTULinux(t.name, mtu)
	case "darwin":
		return setMTUDarwin(t.name, mtu)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func setMTULinux(name string, mtu int) error {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	ifreq, err := unix.NewIfreq(name)
	if err != nil {
		return err
	}

	ifreq.SetUint32(uint32(mtu))
	return unix.IoctlIfreq(fd, unix.SIOCSIFMTU, ifreq)
}

func setMTUDarwin(name string, mtu int) error {
	cmd := exec.Command("ifconfig", name, "mtu", fmt.Sprintf("%d", mtu))
	return cmd.Run()
}

// addRoute adds a route through the interface
func addRoute(name string, dest *net.IPNet) error {
	switch runtime.GOOS {
	case "linux":
		return addRouteLinux(name, dest)
	case "darwin":
		return addRouteDarwin(name, dest)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func addRouteLinux(name string, dest *net.IPNet) error {
	// Use ip command for simplicity and reliability
	cmd := exec.Command("ip", "route", "add", dest.String(), "dev", name)
	return cmd.Run()
}

func addRouteDarwin(name string, dest *net.IPNet) error {
	cmd := exec.Command("route", "add", "-net", dest.String(), "-interface", name)
	return cmd.Run()
}

// Copy copies data between TUN interface and a connection
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
