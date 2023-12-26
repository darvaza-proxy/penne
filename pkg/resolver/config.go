package resolver

// Config describes a [Resolver]
type Config struct {
	Name string `yaml:"name"`
	Next string `yaml:"next,omitempty" toml:",omitempty" json:",omitempty"`

	DisableAAAA bool     `yaml:"disable_aaaa,omitempty" toml:",omitempty" json:",omitempty"`
	Iterative   bool     `yaml:"iterative,omitempty"    toml:",omitempty" json:",omitempty"`
	Recursive   bool     `yaml:"recursive,omitempty"    toml:",omitempty" json:",omitempty"`
	Servers     []string `yaml:"servers,omitempty"      toml:",omitempty" json:",omitempty"`
	Suffixes    []string `yaml:"suffixes,omitempty"     toml:",omitempty" json:",omitempty"`

	Rewrites []RewriteConfig `yaml:"rewrite,omitempty" toml:",omitempty" json:",omitempty"`
}

// RewriteConfig describes an expression used to alter a request
type RewriteConfig struct {
	From string `yaml:"from,omitempty" toml:",omitempty" json:",omitempty"`
	To   string `yaml:"to,omitempty" toml:",omitempty" json:",omitempty"`
}
