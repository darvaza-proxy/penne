package horizon

import "net/netip"

// Config describes a [Horizon]
type Config struct {
	Name string `yaml:"name"`
	Next string `yaml:"next,omitempty" toml:",omitempty" json:",omitempty"`

	AllowForwarding bool `yaml:"allow_forwarding,omitempty" toml:",omitempty" json:",omitempty"`

	Networks []netip.Prefix `yaml:"networks,omitempty" toml:",omitempty" json:",omitempty"`
	Resolver string         `yaml:"resolver,omitempty" toml:",omitempty" json:",omitempty"`
}
