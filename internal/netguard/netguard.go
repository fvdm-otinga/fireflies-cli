// Package netguard provides SSRF-protection helpers for outbound HTTP URLs.
//
// The CLI receives signed upload URLs from the Fireflies GraphQL API and PUTs
// user data to them. To defend against a compromised/spoofed response pointing
// the CLI at an internal address (RFC1918, loopback, link-local, CGNAT) or a
// cloud metadata endpoint, ValidateUploadURL enforces:
//
//   - scheme must be https
//   - hostname must not be localhost / metadata.google.internal / metadata.goog
//   - if the host is a literal IP, it must not fall in any private / loopback
//     / link-local / CGNAT range, and must not be the AWS/GCP metadata IP.
package netguard

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// cgnatNet is the RFC6598 carrier-grade NAT block 100.64.0.0/10.
var cgnatNet = &net.IPNet{
	IP:   net.IPv4(100, 64, 0, 0),
	Mask: net.CIDRMask(10, 32),
}

// awsMetadataIP is the well-known EC2/GCE metadata service IPv4.
var awsMetadataIP = net.ParseIP("169.254.169.254")

// blockedHostnames are literal DNS names we refuse regardless of resolution.
var blockedHostnames = map[string]struct{}{
	"localhost":                {},
	"metadata.google.internal": {},
	"metadata.goog":            {},
}

// ValidateUploadURL parses raw and rejects it if the scheme is not https or
// the host points at an internal / metadata endpoint. The returned *url.URL
// is safe to hand to http.NewRequest.
func ValidateUploadURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid upload URL scheme: %w", err)
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("invalid upload URL scheme: %q (expected https)", u.Scheme)
	}
	host := u.Hostname()
	if host == "" {
		return nil, fmt.Errorf("invalid upload URL scheme: missing host")
	}
	lower := strings.ToLower(host)
	if _, blocked := blockedHostnames[lower]; blocked {
		return nil, fmt.Errorf("invalid upload URL scheme: host %q is not allowed", host)
	}

	if ip := net.ParseIP(host); ip != nil {
		if err := checkIP(ip); err != nil {
			return nil, err
		}
	}

	return u, nil
}

// checkIP rejects IPs in private / loopback / link-local / CGNAT ranges and
// the AWS/GCP metadata address.
func checkIP(ip net.IP) error {
	if awsMetadataIP.Equal(ip) {
		return fmt.Errorf("invalid upload URL scheme: metadata IP %s is not allowed", ip)
	}
	if ip.IsLoopback() {
		return fmt.Errorf("invalid upload URL scheme: loopback address %s is not allowed", ip)
	}
	if ip.IsPrivate() {
		return fmt.Errorf("invalid upload URL scheme: private address %s is not allowed", ip)
	}
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return fmt.Errorf("invalid upload URL scheme: link-local address %s is not allowed", ip)
	}
	if v4 := ip.To4(); v4 != nil && cgnatNet.Contains(v4) {
		return fmt.Errorf("invalid upload URL scheme: CGNAT address %s is not allowed", ip)
	}
	return nil
}
