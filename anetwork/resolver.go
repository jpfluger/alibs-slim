package anetwork

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Resolver contains the data from resolv.conf such as domains, nameservers, search domains, and sortlist.
type Resolver struct {
	Domains     []string // Domains is a list of domain names.
	Nameservers []string // Nameservers is a list of nameserver addresses.
	Search      []string // Search is a list of search domains.
	Sortlist    []string // Sortlist is a list of IP address sorting directives.
}

// GetResolv reads /etc/resolv.conf and returns it as a Resolver.
func GetResolv() (Resolver, error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return Resolver{}, err
	}
	defer file.Close()
	return parseResolvConf(file)
}

// parseResolvConf parses the contents of an io.Reader assumed to be in resolv.conf format.
func parseResolvConf(reader io.Reader) (Resolver, error) {
	var (
		domains     []string
		nameservers []string
		search      []string
		sortlist    []string
	)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue // Skip comments and empty lines
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue // Skip lines that don't have enough fields
		}

		switch fields[0] {
		case "domain":
			domains = append(domains, fields[1:]...)
		case "nameserver":
			nameservers = append(nameservers, fields[1:]...)
		case "search":
			search = append(search, fields[1:]...)
		case "sortlist":
			sortlist = append(sortlist, fields[1:]...)
		}
	}

	if err := scanner.Err(); err != nil {
		return Resolver{}, fmt.Errorf("error reading resolv.conf: %v", err)
	}

	return Resolver{
		Domains:     domains,
		Nameservers: nameservers,
		Search:      search,
		Sortlist:    sortlist,
	}, nil
}

// GetNetResolver returns a custom net.Resolver based on the provided nameserver address.
// If no nameserver is provided, the default system resolver is returned.
func GetNetResolver(nameserver string) *net.Resolver {
	if nameserver == "" {
		return net.DefaultResolver // Use system resolver if no nameserver is specified
	}

	if !strings.Contains(nameserver, ":") {
		nameserver += ":53" // Append the default DNS port if not specified
	}

	dialer := func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{}
		return d.DialContext(ctx, "udp", nameserver) // Use UDP for DNS queries
	}

	return &net.Resolver{PreferGo: true, Dial: dialer}
}
