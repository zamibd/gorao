package cmd

import (
	"encoding/json"
	"net/url"
)

// Options represents console arguments.
type Options struct {
	// Config is the path to the configuration file.
	Config string `long:"config" description:"Path to the YAML configuration file."`

	// DNSListenAddress is the IP address the DNS proxy server will be
	// listening to.
	DNSListenAddress string `long:"dns-address" description:"IP address that the DNS proxy server will be listening to." yaml:"dns_address"`

	// DNSPort is the port the DNS proxy server will be listening to.
	DNSPort int `long:"dns-port" description:"Port the DNS proxy server will be listening to." yaml:"dns_port"`

	// DNSUpstream is the address of the DNS server the proxy will forward
	// queries that are not rewritten to the SNI proxy.
	DNSUpstream string `long:"dns-upstream" description:"The address of the DNS server the proxy will forward queries that are not rewritten by gorao." yaml:"dns_upstream"`

	// DNSRedirectIPV4To is the IPv4 address of the SNI proxy domains will be
	// redirected to by rewriting responses to A queries.
	DNSRedirectIPV4To string `long:"dns-redirect-ipv4-to" description:"IPv4 address that will be used for redirecting type A DNS queries." yaml:"dns_redirect_ipv4_to"`

	// DNSRedirectIPV6To is the IPv6 address of the SNI proxy domains will be
	// redirected to by rewriting responses to AAAA queries.  If not set, the
	// program will try to automatically choose the public address of the SNI
	// proxy.
	DNSRedirectIPV6To string `long:"dns-redirect-ipv6-to" description:"IPv6 address that will be used for redirecting type AAAA DNS queries." yaml:"dns_redirect_ipv6_to"`

	// DNSRedirectRules is a list of wildcards that defines which domains
	// should be redirected to the SNI proxy.  Can be specified multiple times.
	DNSRedirectRules []string `long:"dns-redirect-rule" description:"Wildcard that defines which domains should be redirected to the SNI proxy. Can be specified multiple times." yaml:"dns_redirect_rules"`

	// DNSDropRules is a list of wildcards that define queries to which domains
	// should be dropped.  Can be specified multiple times.
	DNSDropRules []string `long:"dns-drop-rule" description:"Wildcard that defines DNS queries to which domains should be dropped. Can be specified multiple times." yaml:"dns_drop_rules"`

	// DNSCacheEnabled enables DNS response caching.
	DNSCacheEnabled bool `long:"dns-cache-enabled" description:"Enable DNS response caching." yaml:"dns_cache_enabled"`

	// DNSCacheSizeBytes is the size of the DNS cache in bytes.
	DNSCacheSizeBytes int `long:"dns-cache-size" description:"DNS cache size in bytes." yaml:"dns_cache_size"`

	// DNSCacheMinTTL is the minimum TTL for cached entries in seconds.
	DNSCacheMinTTL uint32 `long:"dns-cache-min-ttl" description:"Minimum TTL for cached DNS entries in seconds." yaml:"dns_cache_min_ttl"`

	// DNSCacheMaxTTL is the maximum TTL for cached entries in seconds.
	DNSCacheMaxTTL uint32 `long:"dns-cache-max-ttl" description:"Maximum TTL for cached DNS entries in seconds." yaml:"dns_cache_max_ttl"`

	// HTTPListenAddress is the IP address the HTTP proxy server will be
	// listening to.  Note, that the HTTP proxy will work pretty much the same
	// way the SNI proxy works, i.e. it will tunnel traffic to the hostname
	// that was specified in the "Host" header.
	HTTPListenAddress string `long:"http-address" description:"IP address the SNI proxy server will be listening for plain HTTP connections." yaml:"http_address"`

	// HTTPPort is the port the HTTP proxy server will be listening to.
	HTTPPort int `long:"http-port" description:"Port the SNI proxy server will be listening for plain HTTP connections." yaml:"http_port"`

	// TLSListenAddress is the IP address the SNI proxy server will be
	// listening to.
	TLSListenAddress string `long:"tls-address" description:"IP address the SNI proxy server will be listening for TLS connections." yaml:"tls_address"`

	// TLSPort is the port the SNI proxy server will be listening to.
	TLSPort int `long:"tls-port" description:"Port the SNI proxy server will be listening for TLS connections." yaml:"tls_port"`

	// DOTListenAddress is the IP address the DNS proxy server will be
	// listening for DNS-over-TLS connections.
	DOTListenAddress string `long:"dot-address" description:"IP address the DNS proxy server will be listening for DNS-over-TLS connections." yaml:"dot_address"`

	// DOTPort is the port the DNS proxy server will be listening for
	// DNS-over-TLS connections.
	DOTPort int `long:"dot-port" description:"Port the DNS proxy server will be listening for DNS-over-TLS connections." yaml:"dot_port"`

	// DOHListenAddress is the IP address the DNS proxy server will be
	// listening for DNS-over-HTTPS connections.
	DOHListenAddress string `long:"doh-address" description:"IP address the DNS proxy server will be listening for DNS-over-HTTPS connections." yaml:"doh_address"`

	// DOHPort is the port the DNS proxy server will be listening for
	// DNS-over-HTTPS connections.
	DOHPort int `long:"doh-port" description:"Port the DNS proxy server will be listening for DNS-over-HTTPS connections." yaml:"doh_port"`

	// DOQListenAddress is the IP address the DNS proxy server will be
	// listening for DNS-over-QUIC connections.
	DOQListenAddress string `long:"doq-address" description:"IP address the DNS proxy server will be listening for DNS-over-QUIC connections." yaml:"doq_address"`

	// DOQPort is the port the DNS proxy server will be listening for
	// DNS-over-QUIC connections.
	DOQPort int `long:"doq-port" description:"Port the DNS proxy server will be listening for DNS-over-QUIC connections." yaml:"doq_port"`

	// TLSCertFile is the path to the TLS certificate file.
	TLSCertFile string `long:"tls-cert-file" description:"Path to the TLS certificate file." yaml:"tls_cert_file"`

	// TLSKeyFile is the path to the TLS private key file.
	TLSKeyFile string `long:"tls-key-file" description:"Path to the TLS private key file." yaml:"tls_key_file"`

	// BandwidthRate is a number of bytes per second the connections speed will
	// be limited to.  Note, that the speed is shared between all connections.
	// If not set, there is no limit.
	BandwidthRate float64 `long:"bandwidth-rate" description:"Bytes per second the connections speed will be limited to. If not set, there is no limit." yaml:"bandwidth_rate"`

	// BandwidthRules is a map that allows to define connection speed for
	// domains that match the wildcards.  Has higher priority than
	// BandwidthRate.
	BandwidthRules map[string]float64 `long:"bandwidth-rule" description:"Allows to define connection speed in bytes/sec for domains that match the wildcard. Example: example.*:1024. Can be specified multiple times." yaml:"bandwidth_rules"`

	// ForwardProxy is the address of a SOCKS/HTTP/HTTPS proxy that the connections will
	// be forwarded to according to ForwardRules.
	ForwardProxy string `long:"forward-proxy" description:"Address of a SOCKS/HTTP/HTTPS proxy that the connections will be forwarded to according to forward-rule." yaml:"forward_proxy"`

	// ForwardRules is a list of wildcards that define what connections will be
	// forwarded to ForwardProxy.  If the list is empty and ForwardProxy is set,
	// all connections will be forwarded.
	ForwardRules []string `long:"forward-rule" description:"Wildcard that defines what connections will be forwarded to forward-proxy. Can be specified multiple times. If no rules are specified, all connections will be forwarded to the proxy." yaml:"forward_rules"`

	// ForwardRulesFile is the path to a CSV file containing forward rules (one pattern per line).
	ForwardRulesFile string `long:"forward-rules-file" description:"Path to CSV file with forward rules (one pattern per line)." yaml:"forward_rules_file"`

	// DNSRedirectRulesFile is the path to a CSV file containing DNS redirect rules (one pattern per line).
	DNSRedirectRulesFile string `long:"dns-redirect-rules-file" description:"Path to CSV file with DNS redirect rules (one pattern per line)." yaml:"dns_redirect_rules_file"`

	// BlockRules is a list of wildcards that define connections to which hosts
	// will be blocked.
	BlockRules []string `long:"block-rule" description:"Wildcard that defines connections to which domains should be blocked. Can be specified multiple times." yaml:"block_rules"`

	// BlockRulesFile is the path to a CSV file containing block rules (one pattern per line).
	BlockRulesFile string `long:"block-rules-file" description:"Path to CSV file with block rules (one pattern per line)." yaml:"block_rules_file"`

	// DropRules is a list of wildcards that define connections to which hosts
	// will be "dropped".  "Dropped" means that the connection will be delayed
	// for a hard-coded period of 3 minutes.
	DropRules []string `long:"drop-rule" description:"Wildcard that defines connections to which domains should be dropped (i.e. delayed for a hard-coded period of 3 minutes. Can be specified multiple times." yaml:"drop_rules"`

	// DropRulesFile is the path to a CSV file containing drop rules (one pattern per line).
	DropRulesFile string `long:"drop-rules-file" description:"Path to CSV file with drop rules (one pattern per line)." yaml:"drop_rules_file"`

	// Log settings
	// --

	// Verbose defines whether we should write the DEBUG-level log or not.
	Verbose bool `long:"verbose" description:"Verbose output (optional)" optional:"yes" optional-value:"true" yaml:"verbose"`

	// LogOutput is the optional path to the log file.
	LogOutput string `long:"output" description:"Path to the log file. If not set, write to stdout." yaml:"output"`
}

// String implements fmt.Stringer interface for Options.
// String implements fmt.Stringer interface for Options.
func (o *Options) String() (s string) {
	// Create a shallow copy to avoid modifying the original options
	oCopy := *o
	if oCopy.ForwardProxy != "" {
		if u, err := url.Parse(oCopy.ForwardProxy); err == nil {
			oCopy.ForwardProxy = u.Redacted()
		}
	}

	b, _ := json.MarshalIndent(oCopy, "", "    ")
	return string(b)
}

// DefaultOptions returns the default options.
func DefaultOptions() *Options {
	return &Options{
		DNSListenAddress:  "0.0.0.0",
		DNSPort:           53,
		DNSUpstream:       "8.8.8.8",
		DNSRedirectRules:  []string{"*"},
		DNSCacheEnabled:   true,
		DNSCacheSizeBytes: 64 * 1024 * 1024, // 64MB
		DNSCacheMinTTL:    60,               // 1 minute
		DNSCacheMaxTTL:    3600,             // 1 hour
		HTTPListenAddress: "0.0.0.0",
		HTTPPort:          80,
		TLSListenAddress:  "0.0.0.0",
		TLSPort:           443,
		DOTListenAddress:  "0.0.0.0",
		DOTPort:           853,
		DOHListenAddress:  "0.0.0.0",
		DOHPort:           8443,
		DOQListenAddress:  "0.0.0.0",
		DOQPort:           8853,
	}
}
