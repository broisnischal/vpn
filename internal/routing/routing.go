package routing

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// Route represents a network route
type Route struct {
	Destination *net.IPNet
	Gateway    net.IP
	Interface  string
	Metric     int
}

// Manager handles routing table operations
type Manager struct {
	interfaceName string
}

// NewManager creates a new routing manager
func NewManager(interfaceName string) *Manager {
	return &Manager{
		interfaceName: interfaceName,
	}
}

// AddRoute adds a route through the VPN interface
func (m *Manager) AddRoute(dest *net.IPNet) error {
	switch runtime.GOOS {
	case "linux":
		return m.addRouteLinux(dest)
	case "darwin":
		return m.addRouteDarwin(dest)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func (m *Manager) addRouteLinux(dest *net.IPNet) error {
	cmd := exec.Command("ip", "route", "add", dest.String(), "dev", m.interfaceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Ignore "File exists" error (route already exists)
		if strings.Contains(string(output), "File exists") {
			return nil
		}
		return fmt.Errorf("failed to add route: %w: %s", err, string(output))
	}
	return nil
}

func (m *Manager) addRouteDarwin(dest *net.IPNet) error {
	cmd := exec.Command("route", "add", "-net", dest.String(), "-interface", m.interfaceName)
	return cmd.Run()
}

// DeleteRoute removes a route
func (m *Manager) DeleteRoute(dest *net.IPNet) error {
	switch runtime.GOOS {
	case "linux":
		return m.deleteRouteLinux(dest)
	case "darwin":
		return m.deleteRouteDarwin(dest)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func (m *Manager) deleteRouteLinux(dest *net.IPNet) error {
	cmd := exec.Command("ip", "route", "del", dest.String(), "dev", m.interfaceName)
	return cmd.Run()
}

func (m *Manager) deleteRouteDarwin(dest *net.IPNet) error {
	cmd := exec.Command("route", "delete", "-net", dest.String(), "-interface", m.interfaceName)
	return cmd.Run()
}

// ListRoutes lists all routes
func (m *Manager) ListRoutes() ([]Route, error) {
	switch runtime.GOOS {
	case "linux":
		return m.listRoutesLinux()
	case "darwin":
		return m.listRoutesDarwin()
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func (m *Manager) listRoutesLinux() ([]Route, error) {
	cmd := exec.Command("ip", "route", "show")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var routes []Route
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		
		// Parse route line: "10.0.0.0/24 dev omail0"
		var dest *net.IPNet
		var gateway net.IP
		var iface string
		
		for i, field := range fields {
			if field == "dev" && i+1 < len(fields) {
				iface = fields[i+1]
			}
			if field == "via" && i+1 < len(fields) {
				gateway = net.ParseIP(fields[i+1])
			}
			if i == 0 {
				_, dest, _ = net.ParseCIDR(field)
			}
		}
		
		if dest != nil {
			routes = append(routes, Route{
				Destination: dest,
				Gateway:     gateway,
				Interface:   iface,
			})
		}
	}
	
	return routes, nil
}

func (m *Manager) listRoutesDarwin() ([]Route, error) {
	cmd := exec.Command("netstat", "-rn")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var routes []Route
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		
		// Parse netstat output format
		dest := fields[0]
		gateway := fields[1]
		flags := fields[2]
		iface := fields[3]
		
		if strings.Contains(flags, "U") && iface == m.interfaceName {
			_, ipNet, err := net.ParseCIDR(dest)
			if err != nil {
				continue
			}
			
			routes = append(routes, Route{
				Destination: ipNet,
				Gateway:     net.ParseIP(gateway),
				Interface:   iface,
			})
		}
	}
	
	return routes, nil
}

// SetupDefaultRoute sets up default routing through VPN (for full tunnel)
func (m *Manager) SetupDefaultRoute() error {
	_, defaultRoute, _ := net.ParseCIDR("0.0.0.0/0")
	return m.AddRoute(defaultRoute)
}

// SetupSplitTunnel sets up split tunneling (only route specific networks)
func (m *Manager) SetupSplitTunnel(networks []*net.IPNet) error {
	for _, network := range networks {
		if err := m.AddRoute(network); err != nil {
			return fmt.Errorf("failed to add route for %s: %w", network.String(), err)
		}
	}
	return nil
}

// Cleanup removes all routes through the VPN interface
func (m *Manager) Cleanup() error {
	routes, err := m.ListRoutes()
	if err != nil {
		return err
	}
	
	for _, route := range routes {
		if route.Interface == m.interfaceName {
			if err := m.DeleteRoute(route.Destination); err != nil {
				// Log but continue
				fmt.Printf("Warning: failed to delete route %s: %v\n", route.Destination.String(), err)
			}
		}
	}
	
	return nil
}
