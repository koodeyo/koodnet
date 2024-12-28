package models

import (
	"time"

	"github.com/google/uuid"
)

type Configuration struct {
	// PKI defines the location of credentials for this node. Each of these can also be inlined by using the yaml ": |" syntax.
	PKI configPKI `yaml:"pki" json:"pki" gorm:"serializer:json"`

	// The static host map defines a set of hosts with fixed IP addresses on the internet (or any network).
	// A host can have multiple fixed IP addresses defined here, and nebula will try each when establishing a tunnel.
	// The syntax is:
	// "{nebula ip}": ["{routable ip/dns name}:{routable port}"]
	// Example, if your lighthouse has the nebula IP of 192.168.100.1 and has the real ip address of 100.64.22.11 and runs on port 4242:
	StaticHostmap map[string][]string `yaml:"static_host_map,omitempty" json:"staticHostmap" gorm:"serializer:json"`

	// The static_map config stanza can be used to configure how the static_host_map behaves.
	StaticMap configStaticMap `yaml:"static_map,omitempty" json:"staticMap" gorm:"serializer:json"`

	Lighthouse configLighthouse `yaml:"lighthouse,omitempty" json:"lighthouse" gorm:"serializer:json"`

	// Port Nebula will be listening on. The default here is 4242. For a lighthouse node, the port should be defined,
	// however using port 0 will dynamically assign a port and is recommended for roaming nodes.
	Listen configListen `yaml:"listen,omitempty" json:"listen" gorm:"serializer:json"`

	// Routines is the number of thread pairs to run that consume from the tun and UDP queues.
	// Currently, this defaults to 1 which means we have 1 tun queue reader and 1
	// UDP queue reader. Setting this above one will set IFF_MULTI_QUEUE on the tun
	// device and SO_REUSEPORT on the UDP socket to allow multiple queues.
	// This option is only supported on Linux.
	Routines uint `yaml:"routines,omitempty" json:"routines" example:"1" gorm:"default:1"`

	Punchy configPunchy `yaml:"punchy,omitempty" json:"punchy" gorm:"serializer:json"`

	// Cipher allows you to choose between the available ciphers for your network. Options are chachapoly or aes
	// IMPORTANT: this value must be identical on ALL NODES/LIGHTHOUSES. We do not/will not support use of different ciphers simultaneously!
	Cipher string `yaml:"cipher,omitempty" json:"cipher" example:"aes" enum:"chachapoly,aes"`

	// Preferred ranges is used to define a hint about the local network ranges, which speeds up discovering the fastest
	// path to a network adjacent nebula node.
	// This setting is reloadable.
	PreferredRanges string `yaml:"preferred_ranges,omitempty" json:"preferredRanges" example:"172.16.0.0/24"`

	LocalRange string `yaml:"local_range,omitempty" json:"localRange"`

	// sshd can expose informational and administrative functions via ssh. This can expose informational and administrative
	// functions, and allows manual tweaking of various network settings when debugging or testing.
	SSHD configSSHD `yaml:"sshd,omitempty" json:"sshd" gorm:"serializer:json"`

	// Configure the private interface. Note: addr is baked into the nebula certificate
	Tun configTun `yaml:"tun,omitempty" json:"tun" gorm:"serializer:json"`

	Logging configLogging `yaml:"logging,omitempty" json:"logging" gorm:"serializer:json"`

	Stats configStats `yaml:"stats,omitempty" json:"stats" gorm:"serializer:json"`

	Handshakes configHandshakes `yaml:"handshakes,omitempty" json:"handshakes" gorm:"serializer:json"`

	// Nebula security group configuration
	Firewall configFirewall `yaml:"firewall,omitempty" json:"firewall" gorm:"serializer:json"`

	// EXPERIMENTAL: relay support for networks that can't establish direct connections.
	Relay configRelay `yaml:"relay,omitempty" json:"relay" gorm:"serializer:json"`

	// model
	ID        uuid.UUID `yaml:"id,omitempty" json:"id,omitempty" gorm:"type:uuid;primary_key;" swaggerignore:"true"`
	OwnerID   uuid.UUID `yaml:"owner_id,omitempty" json:"ownerId,omitempty" gorm:"type:uuid;not null;" swaggerignore:"true"`
	OwnerType string    `yaml:"owner_type,omitempty" json:"ownerType,omitempty" gorm:"not null" swaggerignore:"true"`
	CreatedAt time.Time `yaml:"created_at,omitempty" json:"createdAt,omitempty" gorm:"autoCreateTime" swaggerignore:"true"`
	UpdatedAt time.Time `yaml:"updated_at,omitempty" json:"updatedAt,omitempty" gorm:"autoUpdateTime" swaggerignore:"true"`
}

// The configStaticMap config stanza can be used to configure how the static_host_map behaves.
type configStaticMap struct {
	// Cadence determines how frequently DNS is re-queried for updated IP addresses when a static_host_map entry contains a DNS name.
	Cadence string `yaml:"cadence" json:"cadence" example:"30s" gorm:"default:'30s'"`

	// Network determines the type of IP addresses to ask the DNS server for.
	// Valid options are "ip4" (default), "ip6", or "ip" (returns both).
	Network string `yaml:"network" json:"network" example:"ip4" gorm:"default:'ip4'"`

	// LookupTimeout is the DNS query timeout.
	LookupTimeout string `yaml:"lookup_timeout" json:"lookupTimeout" example:"250ms" gorm:"default:'250ms'"`
}

// configPKI defines the structure for storing Public Key Infrastructure (PKI) credentials and associated data.
type configPKI struct {
	// The CA certificate path, used to validate other certificates. Typically located in '/etc/nebula/ca.crt'.
	CA string `yaml:"ca,omitempty" json:"ca" example:"/etc/nebula/ca.crt"`

	// The certificate path for this node. Typically located in '/etc/nebula/host.crt'.
	Cert string `yaml:"cert,omitempty" json:"cert" example:"/etc/nebula/host.crt"`

	// The private key path for this node. Typically located in '/etc/nebula/host.key'.
	Key string `yaml:"key,omitempty" json:"key" example:"/etc/nebula/host.key"`

	// A list of certificate fingerprints that should be blocked. These are certificates the node will not communicate with.
	Blacklist []string `yaml:"blacklist,omitempty" json:"blacklist" example:"c99d4e650533b92061b09918e838a5a0a6aaee21eed1d12fd937682865936c72" gorm:"default:'[]'"`

	// Flag to toggle whether to disconnect clients with expired or invalid certificates.
	DisconnectInvalid bool `yaml:"disconnect_invalid,omitempty" json:"disconnectInvalid" example:"false" gorm:"default:false"`
}

// configLighthouse defines the configuration for lighthouse functionality in the network.
type configLighthouse struct {
	// AmLighthouse is used to enable lighthouse functionality for a node. This should ONLY be true on nodes
	// you have configured to be lighthouses in your network.
	AmLighthouse bool `yaml:"am_lighthouse,omitempty" json:"amLighthouse" example:"false" gorm:"default:false"`

	// ServeDNS optionally starts a DNS listener that responds to various queries and can even be
	// delegated to for resolution.
	ServeDNS bool `yaml:"serve_dns,omitempty" json:"serveDns" example:"false" gorm:"default:false"`

	// DNS holds the DNS configuration for this node.
	DNS configDNS `yaml:"dns,omitempty" json:"dns"`

	// Interval is the number of seconds between updates from this node to a lighthouse.
	// During updates, a node sends information about its current IP addresses to each node.
	Interval uint `yaml:"interval,omitempty" json:"interval" example:"60" gorm:"default:60"`

	// Hosts is a list of lighthouse hosts this node should report to and query from.
	// IMPORTANT: THIS SHOULD BE EMPTY ON LIGHTHOUSE NODES
	// IMPORTANT2: THIS SHOULD BE LIGHTHOUSES' NEBULA IPs, NOT LIGHTHOUSES' REAL ROUTABLE IPs
	Hosts []string `yaml:"hosts,omitempty" json:"hosts" example:"192.168.100.1" gorm:"default:'[]'"`

	// RemoteAllowList is a map of remote hosts that are allowed to communicate with this node.
	// The key is the host's IP address, and the value is a boolean indicating if it's allowed.
	RemoteAllowList map[string]bool `yaml:"remote_allow_list,omitempty" json:"remoteAllowList" example:"172.16.0.0/12:true,0.0.0.0/0:false"`

	// LocalAllowList is a map of local hosts that are allowed to communicate with this node.
	// The key is the host's IP address, and the value is an arbitrary interface{} for future extensibility.
	LocalAllowList map[string]interface{} `yaml:"local_allow_list,omitempty" json:"localAllowList"`

	// RemoteAllowRanges defines more specific remote IP rules for VPN CIDR ranges.
	// This feature is experimental and may change in the future.
	RemoteAllowRanges map[string]map[string]bool `yaml:"remote_allow_ranges,omitempty" json:"remoteAllowRanges"`

	// AdvertiseAddrs are routable addresses that will be included with discovered addresses to report to the lighthouse.
	// This is mainly used for static IPs or port forwarding scenarios where Nebula might not automatically discover them.
	AdvertiseAddrs []string `yaml:"advertise_addrs,omitempty" json:"advertiseAddrs" example:"1.1.1.1:4242,1.2.3.4:0"`

	// CalculatedRemotes is an experimental feature that allows for "guessing" the remote IPs based on Nebula IPs,
	// while waiting for the lighthouse response.
	CalculatedRemotes map[string][]CalculatedRemote `yaml:"calculated_remotes,omitempty" json:"calculatedRemotes"`
}

// configDNS defines the DNS settings for a lighthouse node.
type configDNS struct {
	// Host defines the DNS IP address to bind the DNS listener to. This can also bind to the Nebula node IP.
	Host string `yaml:"host,omitempty" json:"host" example:"0.0.0.0"`

	// Port defines the port for the DNS listener, typically 53.
	Port uint `yaml:"port,omitempty" json:"port" example:"53"`
}

// CalculatedRemote defines a remote IP address based on a mask, used in the "calculated_remotes" field.
type CalculatedRemote struct {
	// Mask defines the network range used to calculate the remote IP.
	Mask string `yaml:"mask,omitempty" json:"mask" example:"192.168.1.0/24"`

	// Port defines the port to be used for the calculated remote IP.
	Port uint `yaml:"port,omitempty" json:"port" example:"4242"`
}

// configListen defines the configuration for listening on a specific host and port for network traffic.
type configListen struct {
	// Host defines the IP address to listen on. To listen on all interfaces, use "0.0.0.0" (IPv4) or "::" (IPv6).
	Host string `yaml:"host,omitempty" json:"host" example:"[::]" gorm:"default:'[::]'"`

	// Port defines the port to listen on. This should be an open port on the system.
	Port uint `yaml:"port,omitempty" json:"port" example:"4242" gorm:"default:4242"`

	// Batch sets the maximum number of packets to pull from the kernel for each syscall.
	// This is used in systems that support recvmmsg (a system call for receiving multiple messages).
	// The default value is 64, and it cannot be reloaded dynamically.
	Batch uint `yaml:"batch,omitempty" json:"batch" example:"64" gorm:"default:64"`

	// ReadBuffer defines the size of the read buffer for the UDP socket. This can be adjusted for performance,
	// especially if the system is receiving a high volume of traffic. The default value is set to 10 MB.
	// This value can be configured in the system's network settings.
	ReadBuffer int64 `yaml:"read_buffer,omitempty" json:"readBuffer" example:"10485760" gorm:"default:10485760"`

	// WriteBuffer defines the size of the write buffer for the UDP socket. Similar to ReadBuffer, it controls
	// the buffer size for outgoing packets. The default value is set to 10 MB.
	WriteBuffer int64 `yaml:"write_buffer,omitempty" json:"writeBuffer" example:"10485760" gorm:"default:10485760"`

	// SendRecvError controls whether Nebula sends "recv_error" packets when it receives data on an unknown tunnel.
	// These packets can help with reconnecting after an abnormal shutdown but could potentially leak information about
	// the system's state. Valid values: "always", "never", or "private".
	// "always" sends the packet in all cases, "never" disables it, and "private" sends it only to private network remotes.
	SendRecvError string `yaml:"send_recv_error,omitempty" json:"sendRecvError" example:"always" gorm:"default:'always'"`
}

// configPunchy defines the configuration for NAT hole punching and response behavior in the network.
type configPunchy struct {
	// Punch defines whether the node should continuously attempt to punch inbound and outbound NAT mappings.
	// This helps avoid the expiration of firewall NAT mappings, ensuring the connection remains active.
	Punch bool `yaml:"punch,omitempty" json:"punch" example:"true" gorm:"default:true"`

	// Respond defines whether the node should connect back if a hole-punching attempt fails. This is useful for
	// situations where one node is behind a difficult NAT (such as symmetric NAT), allowing it to establish a connection.
	// The default value is false.
	Respond bool `yaml:"respond,omitempty" json:"respond" example:"false" gorm:"default:false"`

	// Delay specifies the delay before attempting a punch response for misbehaving NATs. This is particularly useful
	// when dealing with NATs that behave incorrectly. The default value is "1s".
	Delay string `yaml:"delay,omitempty" json:"delay" example:"1s" gorm:"default:'1s'"`

	// RespondDelay sets the delay before attempting punchy.respond, which controls how long the node waits
	// before trying to connect back after a failed hole punch. This only applies if `respond` is set to true.
	// The default value is "5s".
	RespondDelay string `yaml:"respond_delay,omitempty" json:"respondDelay" example:"5s" gorm:"default:'5s'"`
}

// configSSHD defines the configuration for exposing informational and administrative functions via SSH.
type configSSHD struct {
	// Enabled toggles the SSHD feature, allowing SSH access to the node for administrative and debugging tasks.
	Enabled bool `yaml:"enabled,omitempty" json:"enabled" example:"true" gorm:"default:true"`

	// Listen specifies the IP address and port that the SSH server should bind to. The default port 22 is not allowed for safety reasons.
	Listen string `yaml:"listen,omitempty" json:"listen" example:"127.0.0.1:2222"`

	// HostKey specifies the file path to the private key used for SSH host identification.
	HostKey string `yaml:"host_key,omitempty" json:"hostKey" example:"./ssh_host_ed25519_key"`

	// AuthorizedUsers lists the users allowed to authenticate via SSH, along with their corresponding public keys.
	AuthorizedUsers []configAuthorizedUser `yaml:"authorized_users,omitempty" json:"authorizedUsers"`

	// TrustedCAs is a list of trusted SSH Certificate Authorities (CAs) public keys that can sign SSH user keys.
	TrustedCAs []string `yaml:"trusted_cas,omitempty" json:"trustedCas" example:"ssh public key string"`
}

// configAuthorizedUser defines a user authorized to connect via SSH, along with their associated public keys.
type configAuthorizedUser struct {
	// Name specifies the username for SSH access.
	Name string `yaml:"name,omitempty" json:"name" example:"steeeeve"`

	// Keys is a list of authorized SSH public keys for the user. It can contain multiple keys.
	Keys []string `yaml:"keys,omitempty" json:"keys" example:"ssh public key string"`
}

// configTun defines the configuration for the private network interface (TUN).
type configTun struct {
	// When tun is disabled, a lighthouse can be started without a local tun interface (and therefore without root)
	Disabled bool `yaml:"disabled,omitempty" json:"disabled" example:"false" gorm:"default:false"`

	// Dev specifies the name of the TUN device to use. If not set, the OS will choose a default.
	// For macOS: Must be in the form `utun[0-9]+`.
	// For NetBSD: Must be in the form `tun[0-9]+`.
	Dev string `yaml:"dev,omitempty" json:"dev" example:"nebula1" gorm:"default:'nebula1'"`

	// DropLocalBroadcast toggles forwarding of local broadcast packets.
	// The address depends on the IP/mask encoded in the PKI certificate.
	DropLocalBroadcast bool `yaml:"drop_local_broadcast,omitempty" json:"dropLocalBroadcast" example:"false" gorm:"default:false"`

	// DropMulticast toggles forwarding of multicast packets.
	DropMulticast bool `yaml:"drop_multicast,omitempty" json:"dropMulticast" example:"false" gorm:"default:false"`

	// TxQueue sets the transmit queue length. It can help prevent packet drops if increased.
	// The default value is 500.
	TxQueue uint `yaml:"tx_queue,omitempty" json:"txQueue" example:"500" gorm:"default:500"`

	// MTU defines the maximum transmission unit for each packet. The safe default for internet-based traffic is 1300.
	MTU uint `yaml:"mtu,omitempty" json:"mtu" example:"1300" gorm:"default:1300"`

	// Routes defines the network routes that should be added for this TUN interface with MTU overrides.
	Routes []configRoute `yaml:"routes,omitempty" json:"routes"`

	// UnsafeRoutes defines potentially unsafe routes to non-Nebula nodes.
	UnsafeRoutes []configUnsafeRoute `yaml:"unsafe_routes,omitempty" json:"unsafeRoutes"`

	// UseSystemRouteTable allows controlling unsafe routes directly in the system's route table (Linux only).
	UseSystemRouteTable bool `yaml:"use_system_route_table,omitempty" json:"useSystemRouteTable" example:"false" gorm:"default:false"`
}

// configRoute defines a route to be added to the TUN interface with optional MTU overrides.
type configRoute struct {
	// MTU for this specific route. If not set, the default TUN MTU is used.
	MTU uint `yaml:"mtu,omitempty" json:"mtu" example:"8800"`

	// Route is the destination network in CIDR format for this route.
	Route string `yaml:"route,omitempty" json:"route" example:"10.0.0.0/16"`
}

// configUnsafeRoute defines an unsafe route configuration, typically used for routing over Nebula to non-Nebula nodes.
type configUnsafeRoute struct {
	// MTU for the unsafe route. If not set, the default TUN MTU will be used.
	MTU uint `yaml:"mtu,omitempty" json:"mtu" example:"1300"`

	// Route is the destination network in CIDR format.
	Route string `yaml:"route,omitempty" json:"route" example:"172.16.1.0/24"`

	// Via is the gateway IP address for the unsafe route.
	Via string `yaml:"via,omitempty" json:"via" example:"192.168.100.99"`

	// Metric for the unsafe route.
	Metric uint `yaml:"metric,omitempty" json:"metric" example:"100"`

	// Install flag controls whether the route should be installed in the system's routing table.
	Install bool `yaml:"install,omitempty" json:"install" example:"true"`
}

type configLogging struct {
	// Level specifies the logging level.
	// Available options are: panic, fatal, error, warning, info, or debug.
	Level string `yaml:"level,omitempty" json:"level" example:"info" gorm:"default:'info'"`

	// Format specifies the format of the log output.
	// Available options are: json or text.
	Format string `yaml:"format,omitempty" json:"format" example:"text" gorm:"default:'text'"`

	// DisableTimestamp controls whether timestamps are logged. Defaults to false.
	DisableTimestamp bool `yaml:"disable_timestamp,omitempty" json:"disableTimestamp" example:"false" gorm:"default:false"`

	// TimestampFormat specifies the format for timestamps in the log output.
	// Uses Go's time format constants. Leave empty for default behavior.
	TimestampFormat string `yaml:"timestamp_format,omitempty,omitempty" json:"timestampFormat" example:"2006-01-02T15:04:05.000Z07:00" gorm:"default:'2006-01-02T15:04:05Z07:00'"`
}

type configStats struct {
	Type     string `yaml:"type" json:"type" example:"graphite" enum:"graphite,prometheus"`        // Type of stats, e.g., graphite or prometheus
	Interval string `yaml:"interval,omitempty" json:"interval" example:"10s" gorm:"default:'10s'"` // Stats reporting interval

	// Fields for Graphite
	Prefix   string `yaml:"prefix,omitempty" json:"prefix" example:"nebula"`     // Prefix for stats, used in Graphite
	Protocol string `yaml:"protocol,omitempty" json:"protocol" example:"tcp"`    // Protocol for Graphite, e.g., "tcp"
	Host     string `yaml:"host,omitempty" json:"host" example:"127.0.0.1:9999"` // Host for Graphite,

	// Fields for Prometheus
	Listen    string `yaml:"listen,omitempty" json:"listen" example:"127.0.0.1:8080"`     // Address for Prometheus to listen on
	Path      string `yaml:"path,omitempty" json:"path" example:"/metrics"`               // Path for Prometheus metrics
	Namespace string `yaml:"namespace,omitempty" json:"namespace" example:"prometheusns"` // Namespace for Prometheus metrics
	Subsystem string `yaml:"subsystem,omitempty" json:"subsystem" example:"nebula"`       // Subsystem for Prometheus metrics

	// Additional fields
	MessageMetrics    bool `yaml:"message_metrics,omitempty" json:"messageMetrics" example:"false" gorm:"default:false"`       // Enable message metrics
	LighthouseMetrics bool `yaml:"lighthouse_metrics,omitempty" json:"lighthouseMetrics" example:"false" gorm:"default:false"` // Enable lighthouse metrics
}

// Handshake Manager Settings
type configHandshakes struct {
	// Time interval between handshake retries
	TryInterval string `yaml:"try_interval,omitempty" json:"tryInterval" example:"100ms" gorm:"default:'100ms'"`

	// Number of retries for handshakes
	Retries uint `yaml:"retries,omitempty" json:"retries" example:"20" gorm:"default:20"`

	// Size of the buffer channel for querying lighthouses
	QueryBuffer uint `yaml:"query_buffer,omitempty" json:"queryBuffer" example:"64" gorm:"default:64"`

	// Size of the buffer channel for quickly sending handshakes
	TriggerBuffer uint `yaml:"trigger_buffer,omitempty" json:"triggerBuffer" example:"64" gorm:"default:64"`
}

type configFirewall struct {
	// Action for unmatched outbound packets
	OutboundAction string `yaml:"outbound_action,omitempty" json:"outboundAction" example:"drop" gorm:"default:'drop'"`

	// Action for unmatched inbound packets
	InboundAction string `yaml:"inbound_action,omitempty" json:"inboundAction" example:"drop" gorm:"default:'drop'"`

	// Controls the default value for local_cidr. Default is true, will be deprecated after v1.9 and defaulted to false.
	// This setting only affects nebula hosts with subnets encoded in their certificate. A nebula host acting as an
	// unsafe router with `default_local_cidr_any: true` will expose their unsafe routes to every inbound rule regardless
	// of the actual destination for the packet. Setting this to false requires each inbound rule to contain a `local_cidr`
	// if the intention is to allow traffic to flow to an unsafe route.
	DefaultLocalCIDRAny bool `yaml:"default_local_cidr_any,omitempty" json:"defaultLocalCIDRAny" example:"true" gorm:"default:true"`

	// Connection tracking configuration
	ConnTrack configConnTrack `yaml:"conntrack,omitempty" json:"connTrack"`

	// Outbound firewall rules
	Outbound []configFirewallRule `yaml:"outbound,omitempty" json:"outbound"`

	// Inbound firewall rules
	Inbound []configFirewallRule `yaml:"inbound,omitempty" json:"inbound"`
}

type configConnTrack struct {
	// TCP connection timeout
	TcpTimeout string `yaml:"tcp_timeout,omitempty" json:"tcpTimeout" example:"12m" gorm:"default:'12m'"`

	// UDP connection timeout
	UdpTimeout string `yaml:"udp_timeout,omitempty" json:"udpTimeout" example:"3m" gorm:"default:'3m'"`

	// Default connection timeout
	DefaultTimeout string `yaml:"default_timeout,omitempty" json:"defaultTimeout" example:"10m" gorm:"default:'10m'"`
}

type configFirewallRule struct {
	// Port or range of ports
	Port string `yaml:"port,omitempty" json:"port" example:"80" gorm:"default:'any'"`

	// ICMP code (for ICMP-specific rules)
	Code string `yaml:"code,omitempty" json:"code" example:"any" gorm:"default:'any'"`

	// Protocol
	Proto string `yaml:"proto,omitempty" json:"proto" example:"tcp" enum:"tcp,udp,icmp" gorm:"default:'any'"`

	// Specific host to match
	Host string `yaml:"host,omitempty" json:"host"`

	// Specific group to match
	Group string `yaml:"group,omitempty" json:"group"`

	// Multiple groups to match (AND logic)
	Groups []string `yaml:"groups,omitempty" json:"groups"`

	// Remote CIDR to match
	CIDR string `yaml:"cidr,omitempty" json:"cidr"`

	// Local CIDR for unsafe routes
	LocalCIDR string `yaml:"local_cidr,omitempty" json:"localCIDR"`

	// Certificate authority SHA
	CASha string `yaml:"ca_sha,omitempty" json:"caSha" example:"An issuing CA shasum"`

	// Certificate authority name
	CAName string `yaml:"ca_name,omitempty" json:"caName" example:"An issuing CA name"`
}

type configRelay struct {
	// Set to true to permit other hosts to list my IP in their relays config. Default is false.
	AmRelay bool `yaml:"am_relay,omitempty" json:"amRelay" example:"false" gorm:"default:false"`

	// Set to false to prevent this instance from attempting to establish connections through relays. Default is true.
	UseRelays bool `yaml:"use_relays" json:"useRelays" example:"true" gorm:"default:true"`

	// List of Nebula IPs that peers can use to relay packets to this instance.
	// IPs in this list must have am_relay set to true in their configs, or they will reject relay requests.
	Relays []string `yaml:"relays,omitempty" json:"relays" example:"192.168.100.1"`
}

func NewConfig() *Configuration {
	return &Configuration{
		PKI: configPKI{
			Blacklist: []string{},
		},
		StaticHostmap: map[string][]string{},
		Lighthouse: configLighthouse{
			DNS:      configDNS{},
			Interval: 60,
			Hosts:    []string{},
		},
		Listen: configListen{
			Host:  "[::]",
			Port:  4242,
			Batch: 64,
		},
		Punchy: configPunchy{
			Punch: true,
			Delay: "1s",
		},
		Relay: configRelay{
			UseRelays: true,
		},
		Cipher: "aes",
		SSHD: configSSHD{
			AuthorizedUsers: []configAuthorizedUser{},
		},
		Tun: configTun{
			Dev:                "tun1",
			DropLocalBroadcast: true,
			DropMulticast:      true,
			TxQueue:            500,
			MTU:                1300,
			Routes:             []configRoute{},
			UnsafeRoutes:       []configUnsafeRoute{},
		},
		Logging: configLogging{
			Level:  "info",
			Format: "text",
		},
		Stats: configStats{},
		Handshakes: configHandshakes{
			TryInterval: "100ms",
			Retries:     20,
		},
		Firewall: configFirewall{
			ConnTrack: configConnTrack{
				TcpTimeout:     "120h",
				UdpTimeout:     "3m",
				DefaultTimeout: "10m",
			},
			Outbound: []configFirewallRule{
				{
					Port:  "any",
					Proto: "any",
					Host:  "any",
				},
			},
			Inbound: []configFirewallRule{},
		},
	}
}
